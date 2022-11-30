package components

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

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

func postgresWriterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
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
	log.Infof(" csv records !!!")
	ReadCsvToSql(csvFile, currentNode)
	return map[string]interface{}{"out1": "true"}, nil
}
func ReadCsvToSql(r io.Reader, currentNode Node) {
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	// log.Infof("ly---- csv records  is %s", records) //[[number name] [12 23] [3 2] [3 7] [4 6]]
	if err != nil {
		return
	}
	//链接数据库
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

	tablename := currentNode.Config["table"].(string)
	schema := currentNode.Config["databaseChoose"].(string)
	chunksize := currentNode.Config["chunksize"].(string)
	mode := currentNode.Config["mode"].(string)

	if strings.Compare(mode, "replace") == 0 {
		//新建表
		columns := records[0]
		// columns_type := make([]string, 0)
		tableScheamArr := make([]string, 0)
		for i := 0; i < len(columns); i++ {
			// columns_type[i] = "varchar"
			tableScheamArr = append(tableScheamArr, string(columns[i])+" "+"varchar")

		}
		tableScheamStr := strings.Join(tableScheamArr, ",")
		tableCreateStr := fmt.Sprintf("Create Table %s.%s (%s);", schema, tablename, tableScheamStr)
		// log.Infof("ly----tableCreateStr ： %s", tableCreateStr)

		tableDropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", schema, tablename)
		drop_rows, err := db.Query(tableDropStr)
		defer drop_rows.Close()
		if err != nil {
			log.Infof("删除原表失败")
			return
		}
		create_rows, err := db.Query(tableCreateStr)
		defer create_rows.Close()
		if err != nil {
			log.Infof("创建表失败")
			return
		}
		//插入数据

		l := len(records) - 1
		chunksize, err := strconv.Atoi(chunksize)
		n := l/chunksize + 1

		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter*chunksize + 1; i < chunksize*(iter+1)+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)
					for _, col := range records[i] {
						recordsArr = append(recordsArr, "'"+col+"'")
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			} else {
				for i := iter*chunksize + 1; i < l+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)
					for _, col := range records[i] {
						recordsArr = append(recordsArr, "'"+col+"'")
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			}
			if len(tableInsertArr) > 0 {
				tableInsertValues = strings.Join(tableInsertArr, ",")
				tableInsertStr := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, tablename, strings.Join(columns, ","), tableInsertValues)
				// log.Infof("ly----tableInsertStr ： %s", tableInsertStr)
				log.Infof("ly----iter ： %s", iter)
				rows, err := db.Query(tableInsertStr)
				defer rows.Close()
				if err != nil {
					log.Infof("覆盖写入表失败")
					return
				}
			}
		}

	} else {
		//判断表是否存在并获取表头信息
		tableColumnStr := fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s' and table_schema = '%s';", tablename, schema)
		colRows, err := db.Query(tableColumnStr)
		if err != nil {
			log.Infof("数据表检索失败, 请确认要写入的表是否存在")
			return
		}
		tableCols := make([]pgDataCol, 0)
		defer colRows.Close()
		for colRows.Next() {
			var tableCol pgDataCol
			err = colRows.Scan(&tableCol.Name, &tableCol.Type)
			if err != nil {
				log.Infof("数据表检索失败, 请确认要写入的表是否存在")
				return
			}
			tableCols = append(tableCols, tableCol)
		}
		headers := make([]string, 0)
		for _, col := range tableCols {
			headers = append(headers, col.Name)
		}
		// log.Infof("ly----len headers  ： %s", len(headers))
		headersTypes := make([]string, 0)
		for _, col := range tableCols {
			headersTypes = append(headersTypes, col.Type)
		}
		// log.Infof("ly----headersTypes  ： %s", headersTypes)
		if strings.Compare(mode, "clearAndAppend") == 0 {
			log.Infof("开始清空并追加")
			tableClearStr := fmt.Sprintf("TRUNCATE TABLE %s.%s", schema, tablename)
			rows, err := db.Query(tableClearStr)
			defer rows.Close()
			if err != nil {
				log.Infof("清空表失败")
				return
			}
		}
		//插入数据
		l := len(records) - 1
		chunksize, err := strconv.Atoi(chunksize)
		n := l/chunksize + 1
		for iter := 0; iter < n; iter++ {
			var tableInsertValues string
			tableInsertArr := make([]string, 0)
			if iter < n-1 {
				for i := iter*chunksize + 1; i < chunksize*(iter+1)+1; i++ {
					var rowTmpStr string
					recordsArr := make([]string, 0)

					for ctype := 0; ctype < len(headers); ctype++ {
						if len(records[i][ctype]) == 0 && strings.Compare(headersTypes[ctype], "character varying") != 0 {
							recordsArr = append(recordsArr, "NULL")
						} else if len(records[i][ctype]) > 0 && strings.Compare(headersTypes[ctype], "integer") == 0 {
							recordsArr = append(recordsArr, "'"+strings.Split(records[i][ctype], ".")[0]+"'")
						} else {
							recordsArr = append(recordsArr, "'"+records[i][ctype]+"'")
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
						if len(records[i][ctype]) == 0 && strings.Compare(headersTypes[ctype], "character varying") != 0 {
							recordsArr = append(recordsArr, "NULL")
						} else if len(records[i][ctype]) > 0 && strings.Compare(headersTypes[ctype], "integer") == 0 {
							recordsArr = append(recordsArr, "'"+strings.Split(records[i][ctype], ".")[0]+"'")
						} else {
							recordsArr = append(recordsArr, "'"+records[i][ctype]+"'")
						}
					}
					rowTmpStr = "(" + strings.Join(recordsArr, ",") + ")"
					tableInsertArr = append(tableInsertArr, rowTmpStr)
				}
			}
			if len(tableInsertArr) > 0 {
				tableInsertValues = strings.Join(tableInsertArr, ",")
				tableInsertStr := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, tablename, strings.Join(headers, ","), tableInsertValues)
				//log.Infof("ly----tableInsertStr ： %s", tableInsertStr)
				log.Infof("ly----iter ： %s", iter)
				rows, err := db.Query(tableInsertStr)
				defer rows.Close()
				if err != nil {
					log.Infof("追加写入表失败")
					return
				}
			}
		}
	}
	return
}
