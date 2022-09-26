package components

import (
	"strings"
	"sync"

	socketio "github.com/googollee/go-socket.io"
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
	Run           func(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool, server *socketio.Server)
	dumpOutput    func(currentNode Node, outputData map[string]interface{})
	UpdateInput   func(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	loadInput     func(currentNode Node, inputData RequestData) error
	main          func(currentNode Node, inputData RequestData) (map[string]interface{}, error)
	Status        int // 0: stoped 1： running 2： finished -1：error
}

type RequestData struct {
	Data  string
	ID    string
	Extra string
}

func (c *Node) Init(nodeType string) {
	c.Run = Run
	c.UpdateInput = UpdateInput
	c.dumpOutput = dumpOutput
	switch nodeType {
	case "StreamIn":
		c.main = streamInMain
		c.loadInput = streamInLoadInput
	case "StreamOut":
		c.main = streamOutMain
	case "JsonExtractor":
		c.main = jsonExtractorMain
	default:
	}
}

func Run(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool, server *socketio.Server) {
	defer wg.Done()
	select {
	case <-stopChan:
		log.Info("Recive stop event")
	default:
		currentNode.Status = 1
		outputData, err := currentNode.main(currentNode, inputData)
		if err != nil {
			log.Infof("Error occur when running node: %s, error info: %s", currentNode.Key, err.Error())
			currentNode.Status = -1
			if server != nil {
				server.BroadcastToNamespace("/", "notify.process.status", map[string]int{currentNode.Id: -1})
				server.BroadcastToNamespace("/", "notify.process.error", map[string]string{currentNode.Id: err.Error()})
			}
		} else {
			currentNode.dumpOutput(currentNode, outputData)
			currentNode.Status = 2
			if server != nil {
				server.BroadcastToNamespace("/", "notify.process.status", map[string]int{currentNode.Id: 2})
			}
			if len(currentNode.PortConnects["out1"]) > 0 {
				for _, node := range currentNode.NextNodes {
					wg.Add(1)
					go node.Run(*node, RequestData{ID: inputData.ID, Extra: inputData.Extra}, wg, stopChan, server)
				}
			}
		}
	}
}

func UpdateInput(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool) {
	defer wg.Done()
	select {
	case <-stopChan:
		log.Info("Recive stop event")
	default:
		err := currentNode.loadInput(currentNode, inputData)
		if err != nil {
			log.Infof("Error occur when running node: %s, error info: %s", currentNode.Key, err.Error())
		}
	}
}

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
