package components

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
)

type mysqlDataCol struct {
	Name string
	Type string
}

type mysqlDB struct {
	db *sql.DB
	l  *sync.Mutex
}

func mysqlInit(currentNode Node) error {
	mysqlDataType := map[string]string{"bigint": "int64", "bigserial": "int64",
		"boolean": "bool", "bytea": "[]uint8", "date": "time.Time",
		"integer": "int32", "smallint": "int16", "smallserial": "int16",
		"serial": "int32", "text": "string", "time without time zone": "time.Time",
		"time with time zone": "time.Time", "timestamp without time zone": "time.Time",
		"timestamp with time zone": "time.Time", "double precision": "float64", "numeric": "float64"}
	currentNode.Config["mysqlDataType"] = mysqlDataType
	currentNode.Config["mysqlDB"] = &mysqlDB{l: new(sync.Mutex)}
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	go func() {
		mysqluri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["dbname"].(string))
		db, err := sql.Open("mysql", mysqluri)
		if err != nil {
			log.Errorf("Mysql组件(%s)初始化数据库连接失败，请检查配置: %s", currentNode.Id, err.Error())
			currentNode.Config["mysqlConfigFail"] = true
		} else {
			currentNode.Config["mysqlConfigFail"] = false
		}
		if err = db.Ping(); err != nil {
			log.Errorf("Mysql组件(%s)数据库测试连接失败，请检查配置, 具体原因为: %s", currentNode.Id, err.Error())
			currentNode.Config["mysqlConfigFail"] = true
		}
		currentNode.Config["mysqlDB"].(*mysqlDB).db = db
		defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	}()
	return nil
}

func mysqlRlease(currentNode Node) error {
	if currentNode.Config["mysqlDB"] == nil {
		return nil
	}
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	go func() {
		if currentNode.Config["mysqlDB"].(*mysqlDB).db != nil {
			currentNode.Config["mysqlDB"].(*mysqlDB).db.Close()
		}
		currentNode.Config["mysqlDB"].(*mysqlDB).db = nil
		defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	}()
	return nil
}

func rebuildMysqlConnection(currentNode Node) error {
	log.Infof("Mysql组件(%s)尝试重新建立链接", currentNode.Id)
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	mysqluri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["dbname"].(string))
	db, err := sql.Open("mysql", mysqluri)
	if err != nil {
		log.Errorf("Mysql组件(%s)初始化数据库连接失败，请检查配置: %s", currentNode.Id, err.Error())
		currentNode.Config["mysqlConfigFail"] = true
		return err
	} else {
		currentNode.Config["mysqlConfigFail"] = false
	}
	if err = db.Ping(); err != nil {
		log.Errorf("Mysql组件(%s)数据库测试连接失败，请检查配置, 具体原因为: %s", currentNode.Id, err.Error())
		currentNode.Config["mysqlConfigFail"] = true
		return err
	}
	currentNode.Config["mysqlDB"].(*mysqlDB).db = db
	defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	return nil
}

func mysqlReaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	if currentNode.Config["mysqlConfigFail"].(bool) {
		err := rebuildMysqlConnection(currentNode)
		if err != nil {
			return map[string]interface{}{}, err
		}
	}
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	db := currentNode.Config["mysqlDB"].(*mysqlDB).db
	tableCols := make([]mysqlDataCol, 0)
	tableQueryStr := ""
	if len(currentNode.Config["sql"].(string)) == 0 {
		tablename := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
		tableQueryStr = fmt.Sprintf("SELECT * FROM %s", tablename)
	} else {
		tableQueryStr = loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	}
	rows, err := db.Query(tableQueryStr)
	if err != nil {
		log.Errorf("Mysql读取组件(%s)数据表检索失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	columnNames, err := rows.Columns()
	if err != nil {
		log.Errorf("Mysql读取组件(%s)查询数据表结构失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Errorf("Mysql读取组件(%s)查询数据表类型失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	for i, col := range columnNames {
		tableCol := mysqlDataCol{Name: col, Type: columnTypes[i].DatabaseTypeName()}
		tableCols = append(tableCols, tableCol)
	}
	records := make([][]string, 0)
	headers := make([]string, 0)
	headers = append(headers, "indexCol")
	for _, col := range tableCols {
		headers = append(headers, col.Name)
	}
	records = append(records, headers)
	recordNum := 0
	defer rows.Close()
	for rows.Next() {
		record := make([]interface{}, 0, len(tableCols))
		for _, col := range tableCols {
			switch strings.ToLower(col.Type) {
			case "date", "time without time zone", "time with time zone", "timestamp without time zone", "timestamp with time zone":
				record = append(record, sql.NullTime{})
			default:
				record = append(record, sql.NullString{})
			}
		}
		recordP := make([]interface{}, len(tableCols))
		for i := range record {
			recordP[i] = &record[i]
		}
		err = rows.Scan(recordP...)
		if err != nil {
			log.Errorf("Mysql读取组件(%s)数据表数据检索失败: %s", currentNode.Id, err.Error())
			log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
			currentNode.Config["mysqlConfigFail"] = true
			return map[string]interface{}{}, err
		}
		data := make([]string, 0)
		data = append(data, strconv.FormatInt(int64(recordNum), 10))
		for i := range record {
			switch v := record[i].(type) {
			case int64, int16, int32, int8, int, uint, uint16, uint32, uint64:
				data = append(data, strconv.FormatInt(v.(int64), 10))
			case bool:
				data = append(data, strconv.FormatBool(v))
			case float32, float64:
				data = append(data, strconv.FormatFloat(v.(float64), 'E', -1, 32))
			case time.Time:
				if strings.ToLower(tableCols[i].Type) == "date" {
					data = append(data, v.Format("2006-01-02"))
				} else {
					data = append(data, v.Format("2006-01-02 15:04:05"))
				}
			case nil:
				data = append(data, "")
			case []uint8:
				data = append(data, string([]byte(v)))
			default:
				data = append(data, v.(string))
			}
		}
		recordNum += 1
		records = append(records, data)
	}
	os.Mkdir(currentNode.Id, os.ModePerm)
	tmpPath := fmt.Sprintf("%s/data.csv", currentNode.Id)
	os.Remove(tmpPath)
	file, err := os.Create(tmpPath)
	if err != nil {
		log.Errorf("Mysql读取组件(%s)无法创建临时文件: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		return map[string]interface{}{}, err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	err = w.WriteAll(records)
	if err != nil {
		log.Errorf("Mysql读取组件(%s)无法写入csv数据: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		return map[string]interface{}{}, err
	}

	return map[string]interface{}{"out1": tmpPath}, nil
}

func mysqlJsonReaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	if currentNode.Config["mysqlConfigFail"].(bool) {
		err := rebuildMysqlConnection(currentNode)
		if err != nil {
			return map[string]interface{}{}, err
		}
	}
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	db := currentNode.Config["mysqlDB"].(*mysqlDB).db
	tableCols := make([]mysqlDataCol, 0)
	tableQueryStr := ""
	if len(currentNode.Config["sql"].(string)) == 0 {
		tablename := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
		tableQueryStr = fmt.Sprintf("SELECT * FROM %s", tablename)
	} else {
		tableQueryStr = loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	}
	rows, err := db.Query(tableQueryStr)
	if err != nil {
		log.Errorf("Mysql读取组件(%s)数据表检索失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	columnNames, err := rows.Columns()
	if err != nil {
		log.Errorf("Mysql读取组件(%s)查询数据表结构失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Errorf("Mysql读取组件(%s)查询数据表类型失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	for i, col := range columnNames {
		tableCol := mysqlDataCol{Name: col, Type: columnTypes[i].DatabaseTypeName()}
		tableCols = append(tableCols, tableCol)
	}
	records := make([][]string, 0)
	headers := make([]string, 0)
	headers = append(headers, "indexCol")
	for _, col := range tableCols {
		headers = append(headers, col.Name)
	}
	records = append(records, headers)
	recordNum := 0
	defer rows.Close()
	for rows.Next() {
		record := make([]interface{}, 0, len(tableCols))
		for _, col := range tableCols {
			switch strings.ToLower(col.Type) {
			case "date", "time without time zone", "time with time zone", "timestamp without time zone", "timestamp with time zone":
				record = append(record, sql.NullTime{})
			default:
				record = append(record, sql.NullString{})
			}
		}
		recordP := make([]interface{}, len(tableCols))
		for i := range record {
			recordP[i] = &record[i]
		}
		err = rows.Scan(recordP...)
		if err != nil {
			log.Errorf("Mysql读取组件(%s)数据表数据检索失败: %s", currentNode.Id, err.Error())
			log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
			currentNode.Config["mysqlConfigFail"] = true
			return map[string]interface{}{}, err
		}
		data := make([]string, 0)
		data = append(data, strconv.FormatInt(int64(recordNum), 10))
		for i := range record {
			switch v := record[i].(type) {
			case int64, int16, int32, int8, int, uint, uint16, uint32, uint64:
				data = append(data, strconv.FormatInt(v.(int64), 10))
			case bool:
				data = append(data, strconv.FormatBool(v))
			case float32, float64:
				data = append(data, strconv.FormatFloat(v.(float64), 'E', -1, 32))
			case time.Time:
				if strings.ToLower(tableCols[i].Type) == "date" {
					data = append(data, v.Format("2006-01-02"))
				} else {
					data = append(data, v.Format("2006-01-02 15:04:05"))
				}
			case nil:
				data = append(data, "")
			case []uint8:
				data = append(data, string([]byte(v)))
			default:
				data = append(data, v.(string))
			}
		}
		recordNum += 1
		records = append(records, data)
	}

	return map[string]interface{}{"out1": records}, nil
}

func mysqlExecutorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	if currentNode.Config["mysqlConfigFail"].(bool) {
		err := rebuildMysqlConnection(currentNode)
		if err != nil {
			return map[string]interface{}{}, err
		}
	}
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	db := currentNode.Config["mysqlDB"].(*mysqlDB).db
	tableQueryStr := loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	_, err := db.Exec(tableQueryStr)
	if err != nil {
		log.Errorf("Mysql执行组件(%s)数据表执行sql语句失败: %s", currentNode.Id, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		currentNode.Config["mysqlConfigFail"] = true
		return map[string]interface{}{}, err
	}
	return map[string]interface{}{"out1": "success"}, nil
}

func mysqlWriterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	args := config.GetArgs()
	tmpPath := currentNode.InputData["in1"].(string)
	if _, err := os.Stat(tmpPath); errors.Is(err, os.ErrNotExist) {
		tmpPath = path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], currentNode.InputData["in1"].(string), currentNode.Id, "data.csv")
		tmpKey := path.Join(currentNode.InputData["in1"].(string), "data.csv")
		os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
		storageErr := storage.FGetObject(tmpKey, tmpPath)
		if storageErr != nil {
			log.Errorf("Mysql写入组件(%s)无法下载文件: %s, 报错信息为: %s", currentNode.Id, tmpKey, storageErr.Error())
			log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
			return map[string]interface{}{}, err
		}
	}
	csvFile, err := os.Open(tmpPath)
	if err != nil {
		log.Errorf("Mysql写入组件(%s)无法打开文件: %s, 报错信息为: %s", currentNode.Id, tmpPath, err.Error())
		log.Errorf("消息ID为: %s, 消息内容为: %v 的消息运行失败", inputData.ID, currentNode.InputData)
		return map[string]interface{}{}, err
	}
	defer func() {
		csvFile.Close()
		err = os.Remove(tmpPath)
		if err != nil {
			log.Errorf("Mysql写入组件(%s)无法删除临时文件: %s, 报错信息为: %s", currentNode.Id, tmpPath, err.Error())
		}
	}()
	csvToSqlErr := ReadCsvToMySql(csvFile, currentNode)
	if csvToSqlErr != nil {
		log.Errorf("Mysql写入组件(%s)未能正常写入数据库: %s", currentNode.Id, csvToSqlErr.Error())
		return map[string]interface{}{}, err
	}
	return map[string]interface{}{"out1": "success"}, nil
}

func ReadCsvToMySql(r io.Reader, currentNode Node) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if currentNode.Config["mysqlConfigFail"].(bool) {
		err := rebuildMysqlConnection(currentNode)
		if err != nil {
			return err
		}
	}
	currentNode.Config["mysqlDB"].(*mysqlDB).l.Lock()
	defer currentNode.Config["mysqlDB"].(*mysqlDB).l.Unlock()
	db := currentNode.Config["mysqlDB"].(*mysqlDB).db

	tablename := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
	dbname := currentNode.Config["dbname"].(string)
	// schema := currentNode.Config["databaseChoose"].(string)
	chunksizeRaw := currentNode.Config["chunksize"].(string)
	mode := currentNode.Config["mode"].(string)
	chunksize, err := strconv.Atoi(chunksizeRaw)
	if err != nil {
		return err
	}

	if strings.Compare(mode, "replace") == 0 {
		//新建表
		columns := records[0]
		tableSchemaArr := make([]string, 0)
		for i := 1; i < len(columns); i++ {
			tableSchemaArr = append(tableSchemaArr, "`"+string(columns[i])+"`"+" "+"varchar(255)")

		}
		tableSchemaStr := strings.Join(tableSchemaArr, ",")
		tableCreateStr := fmt.Sprintf("Create Table `%s` (%s);", tablename, tableSchemaStr)
		tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename)
		_, err := db.Exec(tableDropStr)
		if err != nil {
			return err
		}
		_, err = db.Exec(tableCreateStr)
		if err != nil {
			return err
		}
		//插入数据
		l := len(records) - 1
		n := l/chunksize + 1

		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter*chunksize + 1; i < chunksize*(iter+1)+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)
					for colIdx, col := range records[i] {
						if colIdx != 0 {
							recordsArr = append(recordsArr, "'"+strings.ReplaceAll(col, "'", "''")+"'")
						}
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			} else {
				for i := iter*chunksize + 1; i < l+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)
					for colIdx, col := range records[i] {
						if colIdx != 0 {
							recordsArr = append(recordsArr, "'"+strings.ReplaceAll(col, "'", "''")+"'")
						}
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			}
			if len(tableInsertArr) > 0 {
				tableInsertValues = strings.Join(tableInsertArr, ",")
				tableColumns := make([]string, 0)
				for i := 1; i < len(columns); i++ {
					tableColumns = append(tableColumns, "`"+string(columns[i])+"`")

				}
				tableInsertStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", tablename, strings.Join(tableColumns, ","), tableInsertValues)
				_, err := db.Exec(tableInsertStr)
				if err != nil {
					return err
				}
			}
		}
	} else {
		//判断表是否存在并获取表头信息
		tableColumnStr := fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s' and table_schema = '%s';", tablename, dbname)
		colRows, err := db.Query(tableColumnStr)
		if err != nil {
			return err
		}
		tableCols := make([]mysqlDataCol, 0)
		defer colRows.Close()
		for colRows.Next() {
			var tableCol mysqlDataCol
			err = colRows.Scan(&tableCol.Name, &tableCol.Type)
			if err != nil {
				return err
			}
			tableCols = append(tableCols, tableCol)
		}
		if len(tableCols) == 0 {
			log.Debug("数据表检索失败, 开始自动创建数据表")
			//新建表
			columns := records[0]
			tableSchemaArr := make([]string, 0)
			for i := 1; i < len(columns); i++ {
				tableSchemaArr = append(tableSchemaArr, "`"+string(columns[i])+"`"+" "+"varchar(255)")

			}
			tableSchemaStr := strings.Join(tableSchemaArr, ",")
			tableCreateStr := fmt.Sprintf("Create Table %s (%s);", tablename, tableSchemaStr)
			tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s", tablename)
			_, err := db.Exec(tableDropStr)
			if err != nil {
				return err
			}
			_, err = db.Exec(tableCreateStr)
			if err != nil {
				return err
			}
			tableColumnStr = fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s' and table_schema = '%s';", tablename, dbname)
			colRows, err := db.Query(tableColumnStr)
			if err != nil {
				return err
			}
			defer colRows.Close()
			for colRows.Next() {
				var tableCol mysqlDataCol
				err = colRows.Scan(&tableCol.Name, &tableCol.Type)
				if err != nil {
					return err
				}
				tableCols = append(tableCols, tableCol)
			}
		}
		headers := make([]string, 0)
		for _, col := range tableCols {
			headers = append(headers, "`"+col.Name+"`")
		}
		headersTypes := make([]string, 0)
		for _, col := range tableCols {
			headersTypes = append(headersTypes, col.Type)
		}
		headerToRecords := make(map[string]int)
		for _, header := range headers {
			colIdx := -1
			for colNum, col := range records[0] {
				if "`"+col+"`" == header {
					colIdx = colNum
				}
			}
			headerToRecords[header] = colIdx
		}
		if strings.Compare(mode, "clearAndAppend") == 0 {
			log.Debug("开始清空并追加")
			tableClearStr := fmt.Sprintf("TRUNCATE TABLE %s", tablename)
			_, err := db.Exec(tableClearStr)
			if err != nil {
				return err
			}
		}
		//插入数据
		l := len(records) - 1
		n := l/chunksize + 1
		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter*chunksize + 1; i < chunksize*(iter+1)+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)
					for ctype := 0; ctype < len(headers); ctype++ {
						if headerToRecords[headers[ctype]] != -1 {
							recordIdx := headerToRecords[headers[ctype]]
							if len(records[i][recordIdx]) == 0 && strings.Compare(headersTypes[ctype], "character varying") != 0 {
								recordsArr = append(recordsArr, "NULL")
							} else if len(records[i][recordIdx]) > 0 && strings.Compare(headersTypes[ctype], "integer") == 0 {
								recordsArr = append(recordsArr, "'"+strings.Split(records[i][recordIdx], ".")[0]+"'")
							} else {
								recordsArr = append(recordsArr, "'"+strings.ReplaceAll(records[i][recordIdx], "'", "''")+"'")
							}
						} else {
							recordsArr = append(recordsArr, "NULL")
						}
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			} else {
				for i := iter*chunksize + 1; i < l+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)
					for ctype := 0; ctype < len(headers); ctype++ {
						if headerToRecords[headers[ctype]] != -1 {
							recordIdx := headerToRecords[headers[ctype]]
							if len(records[i][recordIdx]) == 0 && strings.Compare(headersTypes[ctype], "character varying") != 0 {
								recordsArr = append(recordsArr, "NULL")
							} else if len(records[i][recordIdx]) > 0 && strings.Compare(headersTypes[ctype], "integer") == 0 {
								recordsArr = append(recordsArr, "'"+strings.Split(records[i][recordIdx], ".")[0]+"'")
							} else {
								recordsArr = append(recordsArr, "'"+strings.ReplaceAll(records[i][recordIdx], "'", "''")+"'")
							}
						} else {
							recordsArr = append(recordsArr, "NULL")
						}
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			}
			if len(tableInsertArr) > 0 {
				tableInsertValues = strings.Join(tableInsertArr, ",")
				tableInsertStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", tablename, strings.Join(headers, ","), tableInsertValues)
				_, err := db.Exec(tableInsertStr)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
