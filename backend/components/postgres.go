package components

import (
	"database/sql"
	"fmt"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/lib/pq"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
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
	return map[string]interface{}{"out1": df}, nil
}

func postgresWriterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func postgresExecutorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
