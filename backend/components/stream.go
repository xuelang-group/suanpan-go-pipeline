package components

import (
	"strconv"
	"strings"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
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
	//studio/100026/tmp/55149/b18cba70697a11edbc2631b746db181e/2e9df810697811edb633ab10346ad070/out1
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
		log.Infof("not support json")
		fallthrough
	case "csv":
		log.Infof("not support json")
		fallthrough
	case "image":
		log.Infof("not support json")
		fallthrough
	case "bool":
		log.Infof("not support json")
		fallthrough
	case "array":
		log.Infof("not support json")
		fallthrough
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
	outputData := saveOutputData(currentNode)
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

func saveOutputData(currentNode Node) string {
	switch currentNode.Config["subtype"] {
	case "string":
		return saveAsString(currentNode.InputData["in1"])
	case "number":
		return currentNode.InputData["in1"].(string)
	case "json":
		log.Infof("not support json")
		fallthrough
	case "csv":
		log.Infof("not support json")
		fallthrough
	case "image":
		log.Infof("not support json")
		fallthrough
	case "bool":
		log.Infof("not support json")
		fallthrough
	case "array":
		log.Infof("not support json")
		fallthrough
	default:
		return saveAsString(currentNode.InputData["in1"])
	}
}
