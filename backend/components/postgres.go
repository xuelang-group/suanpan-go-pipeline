package components

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
)

type pgDataCol struct {
	Name string
	Type string
}

func postgresInit(currentNode Node) error {
	postgresDataType := map[string]string{"bigint": "int64", "bigserial": "int64",
		"boolean": "bool", "bytea": "[]uint8", "date": "time.Time",
		"integer": "int32", "smallint": "int16", "smallserial": "int16",
		"serial": "int32", "text": "string", "time without time zone": "time.Time",
		"time with time zone": "time.Time", "timestamp without time zone": "time.Time",
		"timestamp with time zone": "time.Time", "double precision": "float64", "numeric": "float64"}
	currentNode.Config["postgresDataType"] = postgresDataType
	return nil
}
func postgresReaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["dbname"].(string))
	config, err := pgxpool.ParseConfig(psqlconn)
	if err != nil {
		log.Infof("数据库配置解析失败，请检查配置：%s", err.Error())
		return map[string]interface{}{}, nil
	}

	// config.MaxConns = 25                      // 最大连接数
	// config.MinConns = 5                       // 最小连接数
	config.MaxConnIdleTime = 30 * time.Minute // 连接的最大空闲时间
	config.MaxConnLifetime = 1 * time.Hour    // 连接的最大存活时间
	config.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Infof("数据库创建连接池失败，请检查配置：%s", err.Error())
		return map[string]interface{}{}, nil
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		log.Infof("数据库测试连接失败，请检查配置, 具体原因为: %s", err.Error())
		return map[string]interface{}{}, nil
	}
	tableCols := make([]string, 0)
	tableQueryStr := ""
	if len(currentNode.Config["sql"].(string)) == 0 {
		tablename := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
		tableQueryStr = fmt.Sprintf("SELECT * FROM %s.%s", currentNode.Config["schema"].(string), tablename)
	} else {
		tableQueryStr = loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	}

	ctx := context.Background()

	rows, err := pool.Query(ctx, tableQueryStr)
	if err != nil {
		log.Infof("数据表检索失败：%s", err.Error())
		return map[string]interface{}{}, nil
	}
	columnNames := rows.FieldDescriptions()

	for _, col := range columnNames {
		tableCols = append(tableCols, col.Name)
	}
	headers := make([]string, 0)
	headers = append(headers, "indexCol")
	headers = append(headers, tableCols...)

	os.Mkdir(currentNode.Id, os.ModePerm)
	tmpPath := fmt.Sprintf("%s/data.csv", currentNode.Id)
	os.Remove(tmpPath)
	file, err := os.Create(tmpPath)
	if err != nil {
		log.Error("无法创建临时文件")
		return map[string]interface{}{}, nil
	}
	defer file.Close()
	w := csv.NewWriter(file)
	err = w.Write(headers)
	if err != nil {
		log.Error("无法写入csv数据")
		return map[string]interface{}{}, nil
	}

	recordNum := 0
	defer rows.Close()
	for rows.Next() {
		record := make([]interface{}, len(tableCols))
		recordP := make([]interface{}, len(tableCols))
		for i := range record {
			recordP[i] = &record[i]
		}
		err = rows.Scan(recordP...)
		if err != nil {
			log.Infof("数据表数据检索失败")
			return map[string]interface{}{}, nil
		}
		data := make([]string, 0)
		data = append(data, strconv.FormatInt(int64(recordNum), 10))
		for i := range record {
			switch v := record[i].(type) {
			case int64:
				data = append(data, strconv.FormatInt(v, 10))
			case int32:
				data = append(data, strconv.FormatInt(int64(v), 10))
			case int16:
				data = append(data, strconv.FormatInt(int64(v), 10))
			case int8:
				data = append(data, strconv.FormatInt(int64(v), 10))
			case int:
				data = append(data, strconv.FormatInt(int64(v), 10))
			case uint64:
				data = append(data, strconv.FormatUint(v, 10))
			case uint32:
				data = append(data, strconv.FormatUint(uint64(v), 10))
			case uint16:
				data = append(data, strconv.FormatUint(uint64(v), 10))
			case uint8: // 通常为 byte
				data = append(data, strconv.FormatUint(uint64(v), 10))
			case uint:
				data = append(data, strconv.FormatUint(uint64(v), 10))
			case bool:
				data = append(data, strconv.FormatBool(v))
			case float32:
				data = append(data, strconv.FormatFloat(float64(v), 'E', -1, 32))
			case float64:
				data = append(data, strconv.FormatFloat(v, 'E', -1, 32))
			case time.Time:
				if columnNames[i].DataTypeOID == 1082 {
					data = append(data, v.Format("2006-01-02"))
				} else {
					data = append(data, v.Format("2006-01-02 15:04:05"))
				}
			case nil:
				data = append(data, "")
			case []uint8:
				data = append(data, string([]byte(v)))
			case pgtype.Numeric:
				value, _ := v.MarshalJSON()
				data = append(data, string(value))
			case pgtype.Date:
				value, _ := v.MarshalJSON()
				data = append(data, string(value))
			case pgtype.Time:
				duration := time.Duration(v.Microseconds * 1000)
				hours := duration / time.Hour
				duration -= hours * time.Hour
				minutes := duration / time.Minute
				duration -= minutes * time.Minute
				seconds := duration / time.Second
				data = append(data, fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds))
			case pgtype.Timestamp:
				value, _ := v.MarshalJSON()
				data = append(data, string(value))
			case pgtype.Timestamptz:
				value, _ := v.MarshalJSON()
				data = append(data, string(value))
			default:
				data = append(data, fmt.Sprintf("%v", v))
			}
		}
		recordNum += 1
		// records = append(records, data)
		err := w.Write(data)
		if err != nil {
			log.Error("无法写入csv数据")
			return map[string]interface{}{}, nil
		}
	}
	w.Flush()

	return map[string]interface{}{"out1": tmpPath}, nil
}

func postgresExecutorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["dbname"].(string))

	config, err := pgxpool.ParseConfig(psqlconn)
	if err != nil {
		log.Infof("数据库配置解析失败，请检查配置：%s", err.Error())
		return map[string]interface{}{}, nil
	}

	// config.MaxConns = 25                      // 最大连接数
	// config.MinConns = 5                       // 最小连接数
	config.MaxConnIdleTime = 30 * time.Minute // 连接的最大空闲时间
	config.MaxConnLifetime = 1 * time.Hour    // 连接的最大存活时间
	config.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Infof("数据库创建连接池失败，请检查配置：%s", err.Error())
		return map[string]interface{}{}, nil
	}
	defer pool.Close()
	if err = pool.Ping(context.Background()); err != nil {
		log.Infof("数据库测试连接失败，请检查配置, 具体原因为: %s", err.Error())
		return map[string]interface{}{}, nil
	}
	tableQueryStr := loadParameter(currentNode.Config["sql"].(string), currentNode.InputData)
	ctx := context.Background()
	_, err = pool.Exec(ctx, tableQueryStr)
	if err != nil {
		log.Infof("数据库执行sql语句失败")
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": "success"}, nil
}

func postgresWriterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
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
	csvToSqlErr := ReadCsvToSql(csvFile, currentNode)
	if csvToSqlErr != nil {
		log.Errorf("未能正常写入数据库：%s", csvToSqlErr.Error())
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": "success"}, nil
}

func readBatch(csvReader *csv.Reader, batch int) ([][]string, error) {
	records := make([][]string, 0)
	for i := 0; i < batch; i++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			if len(record) > 0 {
				records = append(records, record)
			}
			return records, err
		}
		if err != nil {
			return [][]string{}, err
		}
		records = append(records, record)
	}
	return records, nil
}

func ReadCsvToSql(r io.Reader, currentNode Node) error {
	chunksizeRaw := currentNode.Config["chunksize"].(string)
	chunksize, err := strconv.Atoi(chunksizeRaw)
	if err != nil {
		log.Infof("chunksize设置非数值")
		return err
	}
	csvReader := csv.NewReader(r)
	columns, err := csvReader.Read()
	if err != nil {
		return err
	}
	//链接数据库
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["dbname"].(string))
	config, err := pgxpool.ParseConfig(psqlconn)
	if err != nil {
		log.Infof("数据库配置解析失败，请检查配置：%s", err.Error())
		return err
	}

	// config.MaxConns = 25                      // 最大连接数
	// config.MinConns = 5                       // 最小连接数
	config.MaxConnIdleTime = 30 * time.Minute // 连接的最大空闲时间
	config.MaxConnLifetime = 1 * time.Hour    // 连接的最大存活时间
	config.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Infof("数据库创建连接池失败，请检查配置：%s", err.Error())
		return err
	}
	defer pool.Close()
	if err = pool.Ping(context.Background()); err != nil {
		log.Infof("数据库测试连接失败，请检查配置, 具体原因为: %s", err.Error())
		return err
	}

	tablename := loadParameter(currentNode.Config["table"].(string), currentNode.InputData)
	schema := currentNode.Config["databaseChoose"].(string)
	mode := currentNode.Config["mode"].(string)

	ctx := context.Background()
	if strings.Compare(mode, "replace") == 0 {
		//新建表
		tableSchemaArr := make([]string, 0)
		for i := 1; i < len(columns); i++ {
			tableSchemaArr = append(tableSchemaArr, "\""+string(columns[i])+"\""+" "+"varchar")

		}
		tableSchemaStr := strings.Join(tableSchemaArr, ",")
		tableCreateStr := fmt.Sprintf("Create Table %s.%s (%s);", schema, tablename, tableSchemaStr)
		tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", schema, tablename)

		_, err := pool.Exec(ctx, tableDropStr)
		if err != nil {
			log.Infof("删除原表失败")
			return err
		}
		_, err = pool.Exec(ctx, tableCreateStr)
		if err != nil {
			log.Infof("创建表失败")
			return err
		}

		for {
			records, err := readBatch(csvReader, chunksize)
			if err != nil && err != io.EOF {
				log.Infof("读取csv文件失败")
				return err
			}
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			for i := 0; i < len(records); i++ {
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

			if len(tableInsertArr) > 0 {
				tableInsertValues = strings.Join(tableInsertArr, ",")
				tableColumns := make([]string, 0)
				for i := 1; i < len(columns); i++ {
					tableColumns = append(tableColumns, "\""+string(columns[i])+"\"")

				}
				tableInsertStr := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, tablename, strings.Join(tableColumns, ","), tableInsertValues)
				_, err := pool.Exec(ctx, tableInsertStr)
				if err != nil {
					log.Infof("覆盖写入表失败")
					return err
				}
			}
			if err == io.EOF {
				return nil
			}
		}

	} else {
		//判断表是否存在并获取表头信息
		tableColumnStr := fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s' and table_schema = '%s';", tablename, schema)
		colRows, err := pool.Query(ctx, tableColumnStr)
		if err != nil {
			log.Infof("数据表检索失败, 请确认要写入的表是否存在")
			return err
		}
		tableCols := make([]pgDataCol, 0)
		defer colRows.Close()
		for colRows.Next() {
			var tableCol pgDataCol
			err = colRows.Scan(&tableCol.Name, &tableCol.Type)
			if err != nil {
				log.Infof("数据表检索失败, 请确认要写入的表是否存在")
				return err
			}
			tableCols = append(tableCols, tableCol)
		}
		if len(tableCols) == 0 {
			log.Infof("数据表检索失败, 开始自动创建数据表")
			//新建表
			tableSchemaArr := make([]string, 0)
			for i := 1; i < len(columns); i++ {
				tableSchemaArr = append(tableSchemaArr, "\""+string(columns[i])+"\""+" "+"varchar")

			}
			tableSchemaStr := strings.Join(tableSchemaArr, ",")
			tableCreateStr := fmt.Sprintf("Create Table %s.%s (%s);", schema, tablename, tableSchemaStr)
			tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", schema, tablename)
			_, err := pool.Exec(ctx, tableDropStr)
			if err != nil {
				log.Infof("删除原表失败")
				return err
			}
			_, err = pool.Exec(ctx, tableCreateStr)
			if err != nil {
				log.Infof("创建表失败")
				return err
			}
			tableColumnStr = fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s' and table_schema = '%s';", tablename, schema)
			colRows, err := pool.Query(ctx, tableColumnStr)
			if err != nil {
				log.Infof("数据表检索失败, 请确认要写入的表是否存在")
				return err
			}
			defer colRows.Close()
			for colRows.Next() {
				var tableCol pgDataCol
				err = colRows.Scan(&tableCol.Name, &tableCol.Type)
				if err != nil {
					log.Infof("数据表检索失败, 请确认要写入的表是否存在")
					return err
				}
				tableCols = append(tableCols, tableCol)
			}
		}
		headers := make([]string, 0)
		for _, col := range tableCols {
			headers = append(headers, "\""+col.Name+"\"")
		}
		headersTypes := make([]string, 0)
		for _, col := range tableCols {
			headersTypes = append(headersTypes, col.Type)
		}
		headerToRecords := make(map[string]int)
		for _, header := range headers {
			colIdx := -1
			for colNum, col := range columns {
				if "\""+col+"\"" == header {
					colIdx = colNum
				}
			}
			headerToRecords[header] = colIdx
		}
		if strings.Compare(mode, "clearAndAppend") == 0 {
			log.Infof("开始清空并追加")
			tableClearStr := fmt.Sprintf("TRUNCATE TABLE %s.%s", schema, tablename)
			_, err := pool.Exec(ctx, tableClearStr)
			if err != nil {
				log.Infof("清空表失败")
				return err
			}
		}
		for {
			records, err := readBatch(csvReader, chunksize)
			if err != nil && err != io.EOF {
				log.Infof("读取csv文件失败")
				return err
			}
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			for i := 0; i < len(records); i++ {
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

			if len(tableInsertArr) > 0 {
				tableInsertValues = strings.Join(tableInsertArr, ",")
				tableInsertStr := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, tablename, strings.Join(headers, ","), tableInsertValues)
				_, err := pool.Exec(ctx, tableInsertStr)
				if err != nil {
					log.Infof("追加写入表失败\n执行SQL为：%s\n具体报错为：%s", tableInsertStr, err.Error())
					return err
				}
			}
			if err == io.EOF {
				return nil
			}
		}
	}
}
