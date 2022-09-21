package components

import (
	"strings"
	"sync"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type NodeAction interface {
	Run(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	UpdateInput(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	Main(inputData RequestData) (map[string]interface{}, error)
}

type Node struct {
	PreviousNodes []*Node
	NextNodes     []*Node
	InputData     map[string]interface{}
	OutputData    map[string]interface{}
	PortConnects  map[string][]string
	Config        map[string]interface{}
	Id            string
	Key           string
	Run           func(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	dumpOutput    func(currentNode Node, outputData map[string]interface{})
	UpdateInput   func(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	Main          func(currentNode Node, inputData RequestData) (map[string]interface{}, error)
}

type RequestData struct {
	Data  string
	ID    string
	Extra string
}

func (c *Node) Init(nodeType string) {
	switch nodeType {
	case "StreamIn":
		c.Main = streamInMain
		c.UpdateInput = streamInUpdateInput
	case "StreamOut":
		c.Main = streamOutMain
	case "JsonExtractor":
		c.Main = jsonExtractorMain
	default:
	}
}

func Run(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool) {
	defer wg.Done()
	select {
	case <-stopChan:
		log.Info("Recive stop event")
	default:
		outputData, err := currentNode.Main(currentNode, inputData)
		if err != nil {
			log.Infof("Error occur when running node: %s, error info: %s", currentNode.Key, err.Error())
		} else {
			currentNode.dumpOutput(currentNode, outputData)
			if len(currentNode.PortConnects["out1"]) > 0 {
				for _, node := range currentNode.NextNodes {
					wg.Add(1)
					go node.Run(currentNode, RequestData{ID: inputData.ID, Extra: inputData.Extra}, wg, stopChan)
				}
			}
		}
	}
}

// func (c *Node) Main(inputData RequestData) (map[string]interface{}, error) {
// 	log.Info("LoadInput function not implement.")
// 	return map[string]interface{}{}, nil
// }

func dumpOutput(currentNode Node, outputData map[string]interface{}) {
	for port, data := range outputData {
		for _, tgt := range currentNode.PortConnects[port] {
			tgtInfo := strings.Split(tgt, "-")
			for _, node := range currentNode.NextNodes {
				if node.Id == tgtInfo[0] {
					node.InputData[tgtInfo[1]] = data
				}
			}
		}
	}

}

// func (c *Node) UpdateInput(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool) {
// 	log.Info("UpdateInput function not implement.")
// }
