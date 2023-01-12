package components

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	_ "sqlflow.org/gohive"
)

type hiveDataCol struct {
	Name    string
	Type    string
	Comment string
}

func hiveReaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	hiveConn := fmt.Sprintf("%s:%s@%s:%s/%s?auth=PLAIN",
		currentNode.Config["user"].(string),
		url.QueryEscape(currentNode.Config["password"].(string)),
		currentNode.Config["host"].(string),
		currentNode.Config["port"].(string),
		currentNode.Config["dbname"].(string))
	db, err := sql.Open("hive", hiveConn)
	if err != nil {
		log.Info("数据库连接失败，请检查配置")
		return map[string]interface{}{}, nil
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Info("数据库测试连接失败，请检查配置")
		return map[string]interface{}{}, nil
	}
	tableCols := make([]hiveDataCol, 0)
	tableQueryStr := ""
	if len(currentNode.Config["sql"].(string)) == 0 {
		tableName := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
		tableQueryStr = fmt.Sprintf("SELECT * FROM %s", tableName)
	} else {
		tableQueryStr = loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	}
	rows, err := db.Query(tableQueryStr)
	if err != nil {
		log.Info("数据表检索失败")
		return map[string]interface{}{}, nil
	}
	columnNames, err := rows.Columns()
	if err != nil {
		log.Info("查询数据表结构失败")
		return map[string]interface{}{}, nil
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Info("查询数据表类型失败")
		return map[string]interface{}{}, nil
	}
	for i, col := range columnNames {
		tableCol := hiveDataCol{Name: col, Type: columnTypes[i].DatabaseTypeName()}
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
			log.Info("数据表数据检索失败")
			return map[string]interface{}{}, nil
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
	tmpKey := fmt.Sprintf("studio/%s/tmp/%s/%s/%s/%s", config.GetEnv().SpUserId, config.GetEnv().SpAppId, strings.Join(strings.Split(inputData.ID, "-"), ""), config.GetEnv().SpNodeId, currentNode.Id)
	os.Remove(tmpPath)
	file, err := os.Create(tmpPath)
	if err != nil {
		log.Error("无法创建临时文件")
		return map[string]interface{}{}, nil
	}
	w := csv.NewWriter(file)
	err = w.WriteAll(records)
	if err != nil {
		log.Error("无法写入csv数据")
		return map[string]interface{}{}, nil
	}
	storage.FPutObject(fmt.Sprintf("%s/data.csv", tmpKey), tmpPath)

	return map[string]interface{}{"out1": tmpKey}, nil
}

func hiveExecutorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	hiveConn := fmt.Sprintf("%s:%s@%s:%s/%s?auth=PLAIN",
		currentNode.Config["user"].(string),
		url.QueryEscape(currentNode.Config["password"].(string)),
		currentNode.Config["host"].(string),
		currentNode.Config["port"].(string),
		currentNode.Config["dbname"].(string))
	db, err := sql.Open("hive", hiveConn)
	if err != nil {
		log.Infof("数据库连接失败，请检查配置")
		return map[string]interface{}{}, nil
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Infof("数据库测试连接失败，请检查配置")
		return map[string]interface{}{}, nil
	}
	tableQueryStr := loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	_, err = db.Exec(tableQueryStr)
	if err != nil {
		log.Infof("数据表执行sql语句失败")
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": "success"}, nil
}

func hiveWriterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	args := config.GetArgs()
	tmpPath := currentNode.InputData["in1"].(string)
	if _, err := os.Stat(tmpPath); errors.Is(err, os.ErrNotExist) {
		tmpPath = path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], currentNode.InputData["in1"].(string), currentNode.Id, "data.csv")
		tmpKey := path.Join(currentNode.InputData["in1"].(string), "data.csv")
		os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
		storageErr := storage.FGetObject(tmpKey, tmpPath)
		if storageErr != nil {
			log.Errorf("Can not download file: %s, with error: %s", tmpKey, storageErr.Error())
			return map[string]interface{}{}, nil
		}
	}
	csvFile, err := os.Open(tmpPath)
	if err != nil {
		log.Errorf("Can not open csv file: %s, with error: %s", tmpPath, err.Error())
		return map[string]interface{}{}, nil
	}
	defer func() {
		csvFile.Close()
		err = os.Remove(tmpPath)
		if err != nil {
			log.Errorf("Can not remove csv file: %s, with error: %s", tmpPath, err.Error())
		}
	}()
	csvToSqlErr := ReadCsvSaveToHive(csvFile, currentNode)
	if csvToSqlErr != nil {
		log.Error("未能正常写入数据库")
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": "success"}, nil
}

func ReadCsvSaveToHive(r io.Reader, currentNode Node) error {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	//链接数据库
	hiveConn := fmt.Sprintf("%s:%s@%s:%s/%s?auth=PLAIN",
		currentNode.Config["user"].(string),
		url.QueryEscape(currentNode.Config["password"].(string)),
		currentNode.Config["host"].(string),
		currentNode.Config["port"].(string),
		currentNode.Config["dbname"].(string))
	db, err := sql.Open("hive", hiveConn)
	if err != nil {
		log.Info("数据库连接失败，请检查配置")
		return err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Info("数据库测试连接失败，请检查配置")
		return err
	}

	tableName := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
	chunkSizeRaw := currentNode.Config["chunkSize"].(string)
	mode := currentNode.Config["mode"].(string)
	chunkSize, err := strconv.Atoi(chunkSizeRaw)
	if err != nil {
		log.Info("chunkSize设置非数值")
		return err
	}

	if strings.Compare(mode, "replace") == 0 {
		//新建表
		columns := records[0]
		tableSchemaArr := make([]string, 0)
		for i := 1; i < len(columns); i++ {
			tableSchemaArr = append(tableSchemaArr, "\""+string(columns[i])+"\""+" "+"varchar(100)")

		}
		tableSchemaStr := strings.Join(tableSchemaArr, ",")
		tableCreateStr := fmt.Sprintf("Create Table %s (%s)", tableName, tableSchemaStr)
		tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
		_, err := db.Exec(tableDropStr)
		if err != nil {
			log.Infof("删除原表失败%s", err.Error())
			return err
		}
		_, err = db.Exec(tableCreateStr)
		if err != nil {
			log.Info(tableCreateStr)
			log.Infof("创建表失败%s", err.Error())
			return err
		}
		//插入数据
		l := len(records) - 1
		n := l/chunkSize + 1

		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter*chunkSize + 1; i < chunkSize*(iter+1)+1; i++ {
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
				for i := iter*chunkSize + 1; i < l+1; i++ {
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
					tableColumns = append(tableColumns, "\""+string(columns[i])+"\"")

				}
				tableInsertStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", tableName, strings.Join(tableColumns, ","), tableInsertValues)
				_, err := db.Exec(tableInsertStr)
				if err != nil {
					log.Infof("覆盖写入表失败%s", err.Error())
					return err
				}
			}
		}

	} else {
		//判断表是否存在并获取表头信息
		tableColumnStr := fmt.Sprintf("Describe %s;", tableName)
		colRows, err := db.Query(tableColumnStr)
		if err != nil {
			log.Infof("数据表检索失败, 请确认要写入的表是否存在, %s", err.Error())
			return err
		}
		tableCols := make([]hiveDataCol, 0)
		defer colRows.Close()
		for colRows.Next() {
			var tableCol hiveDataCol
			err = colRows.Scan(&tableCol.Name, &tableCol.Type, &tableCol.Comment)
			if err != nil {
				log.Infof("数据表检索失败, 请确认要写入的表是否存在, %s", err.Error())
				return err
			}
			tableCols = append(tableCols, tableCol)
		}
		if len(tableCols) == 0 {
			log.Info("数据表检索失败, 开始自动创建数据表")
			//新建表
			columns := records[0]
			tableSchemaArr := make([]string, 0)
			for i := 1; i < len(columns); i++ {
				tableSchemaArr = append(tableSchemaArr, "\""+string(columns[i])+"\""+" "+"varchar(100)")

			}
			tableSchemaStr := strings.Join(tableSchemaArr, ",")
			tableCreateStr := fmt.Sprintf("Create Table %s (%s)", tableName, tableSchemaStr)
			tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
			_, err := db.Exec(tableDropStr)
			if err != nil {
				log.Infof("删除原表失败%s", err.Error())
				return err
			}
			_, err = db.Exec(tableCreateStr)
			if err != nil {
				log.Info(tableCreateStr)
				log.Infof("创建表失败%s", err.Error())
				return err
			}
			tableColumnStr = fmt.Sprintf("Describe '%s'", tableName)
			colRows, err := db.Query(tableColumnStr)
			if err != nil {
				log.Infof("数据表检索失败, 请确认要写入的表是否存在, %s", err.Error())
				return err
			}
			defer colRows.Close()
			for colRows.Next() {
				var tableCol hiveDataCol
				err = colRows.Scan(&tableCol.Name, &tableCol.Type)
				if err != nil {
					log.Infof("数据表检索失败, 请确认要写入的表是否存在, %s", err.Error())
					return err
				}
				tableCols = append(tableCols, tableCol)
			}
		}
		headers := make([]string, 0)
		for _, col := range tableCols {
			headers = append(headers, col.Name)
		}
		headersTypes := make([]string, 0)
		for _, col := range tableCols {
			headersTypes = append(headersTypes, col.Type)
		}
		headerToRecords := make(map[string]int)
		for _, header := range headers {
			colIdx := -1
			for colNum, col := range records[0] {
				if col == header {
					colIdx = colNum
				}
			}
			headerToRecords[header] = colIdx
		}
		if strings.Compare(mode, "clearAndAppend") == 0 {
			log.Info("开始清空并追加")
			tableClearStr := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
			_, err := db.Exec(tableClearStr)
			if err != nil {
				log.Info("清空表失败")
				return err
			}
		}
		//插入数据
		l := len(records) - 1
		n := l/chunkSize + 1
		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter*chunkSize + 1; i < chunkSize*(iter+1)+1; i++ {
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
				for i := iter*chunkSize + 1; i < l+1; i++ {
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
				tableInsertStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", tableName, strings.Join(headers, ","), tableInsertValues)
				log.Info(tableInsertStr)
				_, err := db.Exec(tableInsertStr)
				if err != nil {
					log.Infof("追加写入表失败：%s", err.Error())
					return err
				}
			}
		}
	}
	return nil
}
