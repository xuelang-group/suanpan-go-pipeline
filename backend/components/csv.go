package components

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
)

var pathId int

func csvDownloaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	needBasename := currentNode.Config["needBasename"].(bool)
	args := config.GetArgs()
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], "tmp", currentNode.Id, "input", strconv.Itoa(pathId), "data.csv")
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
	pathId = (pathId + 1) % 20
	return map[string]interface{}{"out1": tmpPath}, nil
}

func csvUploaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	needBasename := currentNode.Config["needBasename"].(bool)
	tmpKey := fmt.Sprintf("studio/%s/tmp/%s/%s/%s/%s/data.csv", config.GetEnv().SpUserId, config.GetEnv().SpAppId, strings.Join(strings.Split(inputData.ID, "-"), ""), config.GetEnv().SpNodeId, currentNode.Id)
	tmpPath := currentNode.InputData["in1"].(string)
	storageErr := storage.FPutObject(tmpKey, tmpPath)
	if storageErr != nil {
		log.Errorf("Can not download file: %s, with error: %s", tmpKey, storageErr.Error())
		return map[string]interface{}{}, nil
	}
	if needBasename {
		return map[string]interface{}{"out1": tmpPath[:len(tmpPath)-9]}, nil
	} else {
		return map[string]interface{}{"out1": tmpPath}, nil
	}
}

func CsvToDataFrameMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	//直接传递给下游组件dataframe
	tmpPath := currentNode.InputData["in1"].(string)
	if _, err := os.Stat(tmpPath); errors.Is(err, os.ErrNotExist) {
		log.Errorf("Can not find file: %s", tmpPath)
		return map[string]interface{}{}, nil
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
	return map[string]interface{}{"out1": df}, nil
}
func DataFrameToCsvMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	df := currentNode.InputData["in1"].(dataframe.DataFrame)
	colNames := make([]string, 0, len(df.Names()))
	dataCols := df.Names()[0:len(df.Names())]

	colNames = append(colNames, dataCols...)
	df = df.Select(colNames)
	tmpPath := fmt.Sprintf("%s/data.csv", currentNode.Id)
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	os.Remove(tmpPath)
	file, err := os.Create(tmpPath)
	if err != nil {
		log.Error("无法创建临时文件")
		errors.New("无法创建临时文件")
	}
	defer file.Close()
	df.WriteCSV(file)

	return map[string]interface{}{"out1": tmpPath}, nil
}
