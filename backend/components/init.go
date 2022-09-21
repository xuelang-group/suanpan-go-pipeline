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
}

type RequestData struct {
	Data  string
	ID    string
	Extra string
}

func (c *Node) Run(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool) {
	defer wg.Done()
	select {
	case <-stopChan:
		log.Info("Recive stop event")
	default:
		outputData, err := c.Main(inputData)
		if err != nil {
			log.Infof("Error occur when running node: %s, error info: %s", c.Key, err.Error())
		} else {
			c.dumpOutput(outputData)
			if len(c.PortConnects["out1"]) > 0 {
				for _, node := range c.NextNodes {
					wg.Add(1)
					go node.Run(RequestData{ID: inputData.ID, Extra: inputData.Extra}, wg, stopChan)
				}
			}
		}
	}
}

func (c *Node) Main(inputData RequestData) (map[string]interface{}, error) {
	log.Info("LoadInput function not implement.")
	return map[string]interface{}{}, nil
}

func (c *Node) dumpOutput(outputData map[string]interface{}) {
	for port, data := range outputData {
		for _, tgt := range c.PortConnects[port] {
			tgtInfo := strings.Split(tgt, "-")
			for _, node := range c.NextNodes {
				if node.Id == tgtInfo[0] {
					node.InputData[tgtInfo[1]] = data
				}
			}
		}
	}

}

func (c *Node) UpdateInput(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool) {
	log.Info("UpdateInput function not implement.")
}
