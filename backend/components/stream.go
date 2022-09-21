package components

import (
	"strconv"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type StreamInNode struct {
	Node
}

func streamInUpdateInput(currentNode Node, inputData RequestData) {
	loadedInputData := loadInput(currentNode, inputData.Data)
	currentNode.InputData["in1"] = loadedInputData
}

func streamInMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	return loadInput(currentNode, inputData.Data), nil
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
	sendOutput(currentNode, inputData)
	return map[string]interface{}{}, nil
}

func sendOutput(currentNode Node, inputData RequestData) {
	outputData := saveOutputData(currentNode)
	id := inputData.ID
	extra := inputData.Extra
	r := stream.Request{ID: id, Extra: extra}
	r.Send(map[string]string{
		currentNode.Key: outputData,
	})
}

func saveOutputData(currentNode Node) string {
	switch currentNode.Config["subtype"] {
	case "string":
		return currentNode.InputData["in1"].(string)
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
		return currentNode.InputData["in1"].(string)
	}
}
