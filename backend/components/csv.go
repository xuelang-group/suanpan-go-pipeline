package components

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
)

func csvDownloaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	needBasename := currentNode.Config["needBasename"].(bool)
	args := config.GetArgs()
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], currentNode.InputData["in1"].(string), currentNode.Id, "data.csv")
	tmpKey := currentNode.InputData["in1"].(string)
	if needBasename {
		tmpKey = path.Join(currentNode.InputData["in1"].(string), "data.csv")
	}
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	storageErr := storage.FGetObject(tmpKey, tmpPath)
	if storageErr != nil {
		log.Errorf("Can not download file: %s, with error: %s", tmpKey, storageErr.Error())
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": tmpPath}, nil
}

func CsvToDataFrameMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	//直接传递给下游组件dataframe
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
	df := dataframe.ReadCSV(csvFile)

	// log.Infof("ly---read table %s", df)
	return map[string]interface{}{"out1": df}, nil
}
func DataFrameToCsvMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	//dataframe转成string传递给流输出组件
	df := currentNode.InputData["in1"].(dataframe.DataFrame)
	// index := 0
	// idxSeries := df.Rapply(func(s series.Series) series.Series {
	// 	index++
	// 	return series.Ints(index)

	// })
	// df = df.Mutate(idxSeries.Col("X0")).
	// 	Rename("index", "X0")
	// df = df.Drop(0)

	colNames := make([]string, 0, len(df.Names()))
	// colNames = append(colNames, df.Names()[len(df.Names())-1])
	log.Infof("ly---before1 table %s", colNames)
	// dataCols := df.Names()[0 : len(df.Names())-1]
	dataCols := df.Names()[0:len(df.Names())]

	colNames = append(colNames, dataCols...)
	log.Infof("ly---before2 table %s", colNames)
	df = df.Select(colNames)
	tmpPath := "data.csv"
	tmpKey := fmt.Sprintf("studio/%s/tmp/%s/%s/%s/out1", config.GetEnv().SpUserId, config.GetEnv().SpAppId, strings.Join(strings.Split(inputData.ID, "-"), ""), config.GetEnv().SpNodeId)
	os.Remove(tmpPath)
	file, err := os.Create(tmpPath)
	if err != nil {
		log.Error("无法创建临时文件")
		errors.New("无法创建临时文件")
	}
	df.WriteCSV(file)

	log.Infof("ly---write table %s", df)
	storage.FPutObject(fmt.Sprintf("%s/data.csv", tmpKey), tmpPath)

	return map[string]interface{}{"out1": tmpKey}, nil
}
