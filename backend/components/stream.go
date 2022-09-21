package components

import (
	"strconv"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type StreamInNode struct {
	Node
}

func (c *StreamInNode) UpdateInput(inputData RequestData) {
	loadedInputData := c.loadInput(inputData.Data)
	c.InputData["in1"] = loadedInputData
}

func (c *StreamInNode) Main(inputData RequestData) (map[string]interface{}, error) {
	return c.loadInput(inputData.Data), nil
}

func (c *StreamInNode) loadInput(inputData string) map[string]interface{} {
	switch c.Config["subtype"] {
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

type StreamOutNode struct {
	Node
}

func (c *StreamOutNode) Main(inputData RequestData) (map[string]interface{}, error) {
	c.sendOutput(inputData)
	return map[string]interface{}{}, nil
}

func (c *StreamOutNode) sendOutput(inputData RequestData) {
	outputData := c.saveOutputData()
	id := inputData.ID
	extra := inputData.Extra
	r := stream.Request{ID: id, Extra: extra}
	r.Send(map[string]string{
		c.Key: outputData,
	})
}

func (c *StreamOutNode) saveOutputData() string {
	switch c.Config["subtype"] {
	case "string":
		return c.InputData["in1"].(string)
	case "number":
		return c.InputData["in1"].(string)
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
		return c.InputData["in1"].(string)
	}
}
