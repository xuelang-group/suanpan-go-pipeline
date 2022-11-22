package components

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Infof("数据库连接失败，请检查配置")
		return map[string]interface{}{}, nil
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Infof("数据库测试连接失败，请检查配置")
		return map[string]interface{}{}, nil
	}
	tableColumnStr := fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s' and table_schema = '%s';", currentNode.Config["table"].(string), currentNode.Config["schema"].(string))
	colRows, err := db.Query(tableColumnStr)
	if err != nil {
		log.Infof("数据表检索失败")
		return map[string]interface{}{}, nil
	}
	tableCols := make([]pgDataCol, 0)
	defer colRows.Close()
	for colRows.Next() {
		var tableCol pgDataCol
		err = colRows.Scan(&tableCol.Name, &tableCol.Type)
		if err != nil {
			log.Infof("数据表检索失败")
			return map[string]interface{}{}, nil
		}
		tableCols = append(tableCols, tableCol)
	}

	tableQueryStr := fmt.Sprintf("SELECT * FROM %s.%s", currentNode.Config["schema"].(string), currentNode.Config["table"].(string))
	rows, err := db.Query(tableQueryStr)
	if err != nil {
		log.Infof("数据表检索失败")
		return map[string]interface{}{}, nil
	}
	records := make([][]string, 0)
	headers := make([]string, 0)
	for _, col := range tableCols {
		headers = append(headers, col.Name)
	}
	records = append(records, headers)
	defer rows.Close()
	for rows.Next() {
		record := make([]sql.NullString, len(tableCols))
		recordP := make([]interface{}, len(tableCols))
		for i := range record {
			recordP[i] = &record[i]
		}
		err = rows.Scan(recordP...)
		if err != nil {
			log.Infof("数据表检索失败")
			return map[string]interface{}{}, nil
		}
		data := make([]string, 0)
		for i := range record {
			data = append(data, record[i].String)
		}

		records = append(records, data)
	}
	df := dataframe.LoadRecords(
		records,
		dataframe.DetectTypes(true),
	)
	tmpPath := "data.csv"
	tmpKey := fmt.Sprintf("studio/%s/tmp/%s/%s/%s/out1", config.GetEnv().SpUserId, config.GetEnv().SpAppId, strings.Join(strings.Split(inputData.ID, "-"), ""), config.GetEnv().SpNodeId)
	os.Remove(tmpPath)
	file, err := os.Create(tmpPath)
	if err != nil {
		log.Error("无法创建临时文件")
		errors.New("无法创建临时文件")
	}
	log.Infof("node df path is  %s, tmpKey is %s", file, tmpKey)
	df.WriteCSV(file)
	storage.FPutObject(fmt.Sprintf("%s/data.csv", tmpKey), tmpPath)

	return map[string]interface{}{"out1": tmpKey}, nil
}

func postgresWriterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	// log.Infof("ly---- sql currentNode.Config  is %s", currentNode.Config)
	// log.Infof("ly---- sql inputData  is %s", currentNode.InputData)
	//studio/100026/tmp/55149/35963ba0697d11edbc2631b746db181e/2e9df810697811edb633ab10346ad070/out1
	args := config.GetArgs()
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], currentNode.InputData["in1"].(string), "data.csv")
	tmpKey := path.Join(currentNode.InputData["in1"].(string), "data.csv")
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	storageErr := storage.FGetObject(tmpKey, tmpPath)
	if storageErr != nil {
		log.Errorf("Can not download file: %s, with error: %s", tmpKey, storageErr.Error())
		return map[string]interface{}{"out1": tmpKey}, nil
	}
	csvFile, err := os.Open(tmpPath)
	if err != nil {
		log.Errorf("Can not open csv file: %s, with error: %s", tmpPath, err.Error())
		return map[string]interface{}{"out1": tmpKey}, nil
	}
	defer func() {
		csvFile.Close()
		err = os.Remove(tmpPath)
		if err != nil {
			log.Errorf("Can not remove csv file: %s, with error: %s", tmpPath, err.Error())
		}
	}()
	df := dataframe.ReadCSV(csvFile)

	//log.Infof("ly---- read csv data is %s", df) //"DataFrame\n\n    number name\n 0: 12     23\n 1: 3      2\n 2: 3      7\n 3: 4      6\n    <int>  <int>\n"
	newTableName := currentNode.Config["table"].(string)
	schema := currentNode.Config["databaseChoose"].(string)
	chunksize := currentNode.Config["chunksize"].(string)
	mode := currentNode.Config["mode"].(string)

	csvToSql(currentNode, df, newTableName, schema, mode, chunksize)

	return map[string]interface{}{"out1": "true"}, nil
}

func postgresExecutorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	//log.Infof("ly---- sql currentNode.Config  is %s", currentNode.Config)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["dbname"].(string))

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Infof("数据库连接失败，请检查配置")
		return map[string]interface{}{"out1": "false"}, nil
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Infof("数据库测试连接失败，请检查配置")
		return map[string]interface{}{"out1": "false"}, nil
	}
	tableQueryStr := currentNode.Config["sql"].(string)
	// log.Infof("ly---- sql tableQueryStr  is %s", tableQueryStr)
	rows, err := db.Query(tableQueryStr)
	log.Infof("ly--- execute success ")
	defer rows.Close()
	if err != nil {
		log.Infof("数据表执行sql语句失败")
		return map[string]interface{}{"out1": "false"}, nil
	}
	return map[string]interface{}{"out1": "true"}, nil
}

func csvToSql(currentNode Node, df dataframe.DataFrame, tablename string, schema string, mode string, chunksize string) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", currentNode.Config["host"].(string), currentNode.Config["port"].(string), currentNode.Config["user"].(string), currentNode.Config["password"].(string), currentNode.Config["dbname"].(string))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Infof("数据库连接失败，请检查配置")
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Infof("数据库测试连接失败，请检查配置")
		return
	}

	if strings.Compare(mode, "replace") == 0 {
		//新建表
		columns := df.Names()
		columns_type := df.Types()
		tableScheamArr := make([]string, 0)
		for i := 0; i < len(columns); i++ {
			if strings.Compare(string(columns_type[i]), "int") == 0 || strings.Compare(string(columns_type[i]), "float") == 0 || strings.Compare(string(columns_type[i]), "boolean") == 0 {
				tableScheamArr = append(tableScheamArr, string(columns[i])+" "+string(columns_type[i]))
			} else {
				columns_type[i] = "varchar"
				tableScheamArr = append(tableScheamArr, string(columns[i])+" "+string(columns_type[i]))
			}

		}
		// log.Infof("ly---- csvtosql tableScheamArr is %s ", tableScheamArr)
		tableScheamStr := strings.Join(tableScheamArr, ",")
		// log.Infof("ly---- csvtosql tableScheamStr is %s ", tableScheamStr)
		tableCreateStr := fmt.Sprintf("Create Table %s.%s (%s);", schema, tablename, tableScheamStr)
		// log.Infof("ly---- sql tableCreateStr  is %s", tableCreateStr)

		tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", schema, tablename)
		drop_rows, err := db.Query(tableDropStr)
		defer drop_rows.Close()
		if err != nil {
			log.Infof("删除原表失败")
			return
		}
		create_rows, err := db.Query(tableCreateStr)
		log.Infof("ly--- create table success ")
		defer create_rows.Close()
		if err != nil {
			log.Infof("创建表失败")
			return
		}
		log.Infof("ly---- dataframe map is %s", df.Maps())
		//插入数据
		dfToMaps := df.Maps()
		l := len(dfToMaps)
		chunksize, err := strconv.Atoi(chunksize)
		n := l/chunksize + 1
		//
		//
		var tmpStr string
		var rowTmpStr string

		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter * chunksize; i < chunksize*(iter+1); i++ {
					tmpStr = ""
					rowTmpStr = ""
					row := dfToMaps[i]
					for _, v := range row {
						var value string
						switch vtype := v.(type) {
						case int:
							value = strconv.Itoa(v.(int))
						case int64:
							//fmt.Println(k, "is int", vv)
							value = strconv.FormatInt(v.(int64), 10)
						case float32, float64:
							value = strconv.FormatFloat(vtype.(float64), 'g', 12, 64)
						default:
							value = v.(string)

						}
						tmpStr = tmpStr + value + ","
					}
					log.Infof("ly---- brfore tmpStr  is %s", tmpStr)
					rowTmpStr = "(" + tmpStr[0:len(tmpStr)-1] + ")"
					log.Infof("ly---- rowTmpStr  is %s", rowTmpStr)
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			} else {
				for i := iter * chunksize; i < l; i++ {
					tmpStr = ""
					rowTmpStr = ""
					row := dfToMaps[i]
					for _, v := range row {
						var value string
						switch vtype := v.(type) {
						case int:
							value = strconv.Itoa(v.(int))
						case int64:
							//fmt.Println(k, "is int", vv)
							value = strconv.FormatInt(v.(int64), 10)
						case float32, float64:
							value = strconv.FormatFloat(vtype.(float64), 'g', 12, 64)
						default:
							value = v.(string)

						}
						tmpStr = tmpStr + value + ","
					}
					log.Infof("ly---- brfore tmpStr  is %s", tmpStr)
					rowTmpStr = "(" + tmpStr[0:len(tmpStr)-1] + ")"
					log.Infof("ly---- rowTmpStr  is %s", rowTmpStr)
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			}
			tableInsertValues = strings.Join(tableInsertArr, ",")
			log.Infof("ly---- tableInsertValues  is %s", tableInsertValues) //(12,23),(3,2),(3,7),(4,6)
			colnames := df.Names()
			tableInsertStr := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, tablename, strings.Join(colnames, ","), tableInsertValues)
			// log.Infof("ly---- tableInsertStr  is %s", tableInsertStr)
			rows, err := db.Query(tableInsertStr)
			log.Infof("ly--- replace wirte table success ")
			defer rows.Close()
			if err != nil {
				log.Infof("覆盖写入表失败")
				return
			}
		}

	}
	return
}