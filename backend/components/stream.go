package components

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

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
	tmpKey := fmt.Sprintf("studio/%s/tmp/%s/%s/%s/%s", config.GetEnv().SpUserId, config.GetEnv().SpAppId, strings.Join(strings.Split(inputData.ID, "-"), ""), config.GetEnv().SpNodeId, strings.Replace(currentNode.Key, "outputData", "out", -1))
	storage.FPutObject(fmt.Sprintf("%s/data.csv", tmpKey), currentNode.InputData["in1"].(string))
	os.Remove(currentNode.InputData["in1"].(string))
	return tmpKey
}

func csvFileDownload(data string, id string) string {
	args := config.GetArgs()
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], data, id, "data.csv")
	tmpKey := path.Join(data, "data.csv")
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	storage.FGetObject(tmpKey, tmpPath)
	return tmpPath
}
