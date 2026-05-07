package components

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/extrame/xls"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
	"github.com/xuri/excelize/v2"
)

type idCounter struct {
	mu       sync.Mutex
	counters map[string]int
}

var idCtr = &idCounter{
	counters: make(map[string]int),
}

func (c *idCounter) nextIndex(id string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.counters[id]; !exists {
		c.counters[id] = 0
		return 0
	}

	c.counters[id]++
	if c.counters[id] >= 10 {
		c.counters[id] = 0
	}

	return c.counters[id]
}

func streamInLoadInput(currentNode Node, inputData RequestData) error {
	currentNode.InputData["in1"] = inputData.Data
	return nil
}

func streamInMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	if len(inputData.Data) == 0 {
		return map[string]interface{}{}, nil
	}
	if len(inputData.Data) > 0 {
		return loadInput(currentNode, inputData.Data), nil
	} else {
		if currentNode.InputData["in1"] == nil {
			return map[string]interface{}{}, nil
		}
		return loadInput(currentNode, currentNode.InputData["in1"].(string)), nil
	}
}

func loadInput(currentNode Node, inputData string) map[string]interface{} {
	switch currentNode.Config["subtype"] {
	case "string":
		return map[string]interface{}{"out1": inputData}
	case "number":
		inputFloat, _ := strconv.ParseFloat(inputData, 32)
		return map[string]interface{}{"out1": inputFloat}
	case "json":
		var v interface{}
		json.Unmarshal([]byte(inputData), &v)
		return map[string]interface{}{"out1": v}
	case "csv":
		return map[string]interface{}{"out1": csvFileDownload(inputData, currentNode.Id)}
	case "file":
		return map[string]interface{}{"out1": fileDownload(inputData, currentNode.Id)}
	case "image":
		log.Infof("not support image")
		fallthrough
	case "bool":
		if inputData == "true" {
			return map[string]interface{}{"out1": true}
		} else {
			return map[string]interface{}{"out1": false}
		}
	case "array":
		var v []interface{}
		json.Unmarshal([]byte(inputData), &v)
		return map[string]interface{}{"out1": v}
	default:
		return map[string]interface{}{"out1": inputData}
	}
}

func streamOutMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	if currentNode.InputData["in1"] == nil {
		return map[string]interface{}{}, nil
	}
	sendOutput(currentNode, inputData)
	return map[string]interface{}{}, nil
}

func sendOutput(currentNode Node, inputData RequestData) {
	outputData := saveOutputData(currentNode, inputData)
	id := inputData.ID
	extra := inputData.Extra
	r := stream.Request{ID: id, Extra: extra}
	r.Send(map[string]string{
		strings.Replace(currentNode.Key, "outputData", "out", -1): outputData,
	})

}

func saveAsString(outputData interface{}) string {
	var outputString string
	switch i := outputData.(type) {
	case int, int16, int32, int8, int64:
		outputString = strconv.FormatInt(i.(int64), 10)
	case float32, float64:
		outputString = strconv.FormatFloat(i.(float64), 'g', 12, 64)
	default:
		outputString = outputData.(string)
	}
	return outputString
}

func saveOutputData(currentNode Node, inputData RequestData) string {
	switch currentNode.Config["subtype"] {
	case "string":
		return saveAsString(currentNode.InputData["in1"])
	case "number":
		return currentNode.InputData["in1"].(string)
	case "json":
		output, _ := json.Marshal(currentNode.InputData["in1"])
		return string(output)
	case "csv":
		return csvFileUpload(currentNode, inputData)
	case "image":
		log.Infof("not support image")
		fallthrough
	case "bool":
		output, _ := json.Marshal(currentNode.InputData["in1"])
		return string(output)
	case "array":
		output, _ := json.Marshal(currentNode.InputData["in1"])
		return string(output)
	default:
		return saveAsString(currentNode.InputData["in1"])
	}
}

func csvFileUpload(currentNode Node, inputData RequestData) string {
	tmpKey := fmt.Sprintf("studio/%s/tmp/%s/%s/%s", config.GetEnv().SpUserId, config.GetEnv().SpAppId, config.GetEnv().SpNodeId, currentNode.Id)
	storage.FPutObject(fmt.Sprintf("%s/data.csv", tmpKey), currentNode.InputData["in1"].(string))
	os.Remove(currentNode.InputData["in1"].(string))
	return tmpKey
}

func csvFileDownload(data string, id string) string {
	pathIndex := idCtr.nextIndex(id)
	tmpPath, tmpKey := buildTempDownloadPathWithIndex(data, id, pathIndex, "data.csv")
	log.Infof(tmpPath)
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	storage.FGetObject(tmpKey, tmpPath)
	return tmpPath
}

func fileDownload(data string, id string) string {
	objects, err := storage.ListObjects(data, true, 100)
	if err != nil {
		log.Errorf("Can not list files from %s, with error: %s", data, err.Error())
		return ""
	}

	suffixPriority := []string{".csv", ".xlsx", ".xls"}
	for _, suffix := range suffixPriority {
		fileName := "data" + suffix
		for _, object := range objects {
			if path.Base(object.Name) != fileName {
				continue
			}
			pathIndex := idCtr.nextIndex(id)
			sourcePath, _ := buildTempDownloadPathWithIndex(data, id, pathIndex, fileName)
			csvPath, _ := buildTempDownloadPathWithIndex(data, id, pathIndex, "data.csv")
			log.Infof(sourcePath)
			os.MkdirAll(filepath.Dir(sourcePath), os.ModePerm)
			storageErr := storage.FGetObject(object.Name, sourcePath)
			if storageErr != nil {
				log.Errorf("Can not download file: %s, with error: %s", object.Name, storageErr.Error())
				return ""
			}
			finalPath, convertErr := convertFileToCSV(sourcePath, csvPath)
			if convertErr != nil {
				log.Errorf("Can not convert file: %s, with error: %s", object.Name, convertErr.Error())
				return ""
			}
			log.Infof("Download file %s from %s successfully", object.Name, data)
			return finalPath
		}
	}

	log.Errorf("Can not find supported data file under %s", data)
	return ""
}

func buildTempDownloadPath(data string, id string, fileName string) (string, string) {
	pathIndex := idCtr.nextIndex(id)
	return buildTempDownloadPathWithIndex(data, id, pathIndex, fileName)
}

func buildTempDownloadPathWithIndex(data string, id string, pathIndex int, fileName string) (string, string) {
	args := config.GetArgs()
	pathWithIndex := fmt.Sprintf("%s/%d", id, pathIndex)
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], pathWithIndex, fileName)
	tmpKey := path.Join(data, fileName)
	return tmpPath, tmpKey
}

func convertFileToCSV(sourcePath string, csvPath string) (string, error) {
	switch strings.ToLower(filepath.Ext(sourcePath)) {
	case ".csv":
		if sourcePath == csvPath {
			return csvPath, nil
		}
		if err := os.Rename(sourcePath, csvPath); err != nil {
			return "", err
		}
		return csvPath, nil
	case ".xlsx":
		if err := xlsxToCSV(sourcePath, csvPath); err != nil {
			return "", err
		}
	case ".xls":
		if err := xlsToCSV(sourcePath, csvPath); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported file extension: %s", filepath.Ext(sourcePath))
	}

	if sourcePath != csvPath {
		if err := os.Remove(sourcePath); err != nil && !os.IsNotExist(err) {
			log.Errorf("Can not remove temp source file: %s, with error: %s", sourcePath, err.Error())
		}
	}
	return csvPath, nil
}

func xlsxToCSV(sourcePath string, csvPath string) error {
	file, err := excelize.OpenFile(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("xlsx file has no sheets")
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return err
	}

	return writeCSVRows(csvPath, rows)
}

func xlsToCSV(sourcePath string, csvPath string) error {
	workbook, err := xls.Open(sourcePath, "utf-8")
	if err != nil {
		return err
	}
	if workbook.NumSheets() == 0 {
		return fmt.Errorf("xls file has no sheets")
	}

	sheet := workbook.GetSheet(0)
	if sheet == nil {
		return fmt.Errorf("can not read first sheet from xls file")
	}

	rows := make([][]string, 0, int(sheet.MaxRow)+1)
	for rowIndex := 0; rowIndex <= int(sheet.MaxRow); rowIndex++ {
		row := sheet.Row(rowIndex)
		if row == nil {
			rows = append(rows, []string{})
			continue
		}

		cols := make([]string, row.LastCol())
		for colIndex := 0; colIndex < row.LastCol(); colIndex++ {
			cols[colIndex] = row.Col(colIndex)
		}
		rows = append(rows, cols)
	}

	return writeCSVRows(csvPath, rows)
}

func writeCSVRows(csvPath string, rows [][]string) error {
	os.Remove(csvPath)
	csvFile, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return writer.Error()
}
