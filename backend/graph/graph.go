package graph

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goPipeline/components"
	"goPipeline/utils"
	"goPipeline/web"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	"gopkg.in/yaml.v3"
)

type Graph struct {
	Status       bool
	Nodes        []reflect.Value
	Components   []utils.Component
	Config       utils.GraphConfig
	NodeInfo     utils.NodeInfo
	stopChan     chan bool
	wg           sync.WaitGroup
	typeRegistry map[string]reflect.Type
}

func (g *Graph) Init() {
	g.componentsInit()
	g.typesInit()
	g.graphInit()
	g.nodesInit()
}

func (g *Graph) typesInit() {
	myTypes := []interface{}{components.StreamInNode{}, components.StreamOutNode{}}
	for _, v := range myTypes {
		g.typeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}
}

func (g *Graph) graphInit() {
	log.Info("Init function not implement.")
	err := storage.FGetObject(web.GraphKey, web.GraphPath)
	if err != nil {
		log.Info("Fail to Load Config File, init with default value...")
		g.Config = utils.GraphConfig{}
	} else {
		jsonFile, err := os.Open(web.GraphPath)
		if err != nil {
			log.Info(err.Error())
			g.Config = utils.GraphConfig{}
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &g.Config)
		log.Info(fmt.Sprintf("Successfully Loaded Config File %s.", web.GraphPath))
	}
}

func (g *Graph) nodesInit() {
	for _, nodeConfig := range g.Config.Nodes {
		var realNode reflect.Value
		if strings.HasPrefix(nodeConfig.Key, "in") {
			realNode = reflect.New(g.typeRegistry["components.StreamInNode"])
		} else if strings.HasPrefix(nodeConfig.Key, "in") {
			realNode = reflect.New(g.typeRegistry["components.StreamOutNode"])
		} else {
			realNode = reflect.New(g.typeRegistry[fmt.Sprintf("components.%s", nodeConfig.Key)])
		}

		node := components.Node{Id: nodeConfig.Uuid, Key: nodeConfig.Key}
		params := make(map[string]interface{})
		for _, param := range nodeConfig.Parameters {
			params[param["key"]] = param["value"]
		}
		if strings.HasPrefix(node.Key, "in") || strings.HasSuffix(node.Key, "in") {
			var subtype string
			if strings.HasPrefix(node.Key, "in") {
				subtype = g.NodeInfo.Inputs[node.Key].Subtype
			}
			if strings.HasPrefix(node.Key, "out") {
				subtype = g.NodeInfo.Outputs[node.Key].Subtype
			}
			params["subtype"] = subtype
		}
		node.Config = params
		for _, component := range g.Components {
			if component.Key == node.Key {
				for _, port := range component.Ports.In {
					node.InputData[port.Id] = nil
				}
				for _, port := range component.Ports.Out {
					node.InputData[port.Id] = nil
				}
			}
		}

		nodeJson, _ := json.Marshal(node)
		json.Unmarshal(nodeJson, &realNode)

		g.Nodes = append(g.Nodes, realNode)
	}
	for _, connection := range g.Config.Connectors {
		for _, node := range g.Nodes {
			if node.FieldByName("Id").String() == connection.Src["uuid"] {
				if !g.checkNode(connection.Tgt["uuid"], node.FieldByName("NextNodes")) {
					node.FieldByName("NextNodes").Set(reflect.Append(node.FieldByName("NextNodes").Elem(), reflect.ValueOf(g.findNode(connection.Tgt["uuid"]))))
					if !utils.SlicesContain(node.FieldByName("PortConnects").Interface()[connection.Src["port"]], connection.Tgt["uuid"]+"-"+connection.Tgt["port"]) {
						node.PortConnects[connection.Src["port"]] = append(node.PortConnects[connection.Src["port"]], connection.Tgt["uuid"]+"-"+connection.Tgt["port"])
					}
				}
			}
			if node.FieldByName("Id").String() == connection.Tgt["uuid"] {
				if !g.checkNode(connection.Src["uuid"], node.NextNodes) {
					node.PreviousNodes = append(node.PreviousNodes, g.findNode(connection.Tgt["uuid"]))
				}
			}
		}
	}
}

func (g *Graph) findNode(uuid string) *reflect.Value {
	for _, node := range g.Nodes {
		if node.FieldByName("Id").String() == uuid {
			return &node
		}
	}
	return nil
}

func (g *Graph) checkNode(uuid string, nodes reflect.Value) bool {
	for _, node := range nodes {
		if node.Id == uuid {
			return true
		}
	}
	return false
}

func (g *Graph) Update(newGraph utils.GraphConfig) {
	g.Config = newGraph
	os.Remove(web.DataPath)
	dataJson, _ := json.Marshal(g.Config)
	os.WriteFile(web.DataPath, dataJson, 0644)
	storage.FPutObject(web.DataKey, web.DataPath)
	g.Nodes = []components.Node{}
	g.nodesInit()
}

func (g *Graph) Run(inputData map[string]string, id string, extra string) {
	log.Info("Start To Run Graph.")
	g.wg = sync.WaitGroup{}
	g.stopChan = make(chan bool)
	for _, node := range g.Nodes {
		if len(node.PreviousNodes) == 0 {
			g.wg.Add(1)
			if strings.HasPrefix(node.Key, "in") {
				go node.Run(components.RequestData{Data: inputData[node.Key], ID: id, Extra: extra}, &g.wg, g.stopChan)
			} else {
				go node.Run(components.RequestData{ID: id, Extra: extra}, &g.wg, g.stopChan)
			}
		}
	}
	g.wg.Wait()
	_, ok := (<-g.stopChan)
	if ok {
		close(g.stopChan)
	}
	log.Info("Graph Run Done.")
}

func (g *Graph) UpdateInputs(inputData map[string]string, id string, extra string) {
	log.Info("Start To Update Inputs.")
	g.wg = sync.WaitGroup{}
	g.stopChan = make(chan bool)
	for _, node := range g.Nodes {
		if len(node.PreviousNodes) == 0 {
			g.wg.Add(1)
			if strings.HasPrefix(node.Key, "in") {
				go node.UpdateInput(components.RequestData{Data: inputData[node.Key], ID: id, Extra: extra}, g.stopChan)
			}
		}
	}
	g.wg.Wait()
	_, ok := (<-g.stopChan)
	if ok {
		close(g.stopChan)
	}
	log.Info("Update Inputs Done.")
}

func (g Graph) Stop() {
	log.Info("Stop Graph")
	close(g.stopChan)
}

func (g *Graph) componentsInit() {
	files, err := ioutil.ReadDir("configs")
	if err != nil {
		log.Error(err.Error())
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") {
			if f.Name() == "streamConnector.yml" {
				componentConfig := []utils.Component{}
				nodeInfoString, _ := base64.StdEncoding.DecodeString(os.Getenv("SP_NODE_INFO"))
				json.Unmarshal(nodeInfoString, &g.NodeInfo)
				ymlFile, err := os.Open(path.Join("configs", f.Name()))
				if err != nil {
					log.Info(err.Error())
				}
				defer ymlFile.Close()
				byteValue, _ := io.ReadAll(ymlFile)
				yaml.Unmarshal(byteValue, &componentConfig)
				inPortConfig := new(utils.Component)
				for _, v := range componentConfig {
					if v.Category == "inPorts" {
						inPortConfig = &v
						break
					}
				}
				for inputName, inputInfo := range g.NodeInfo.Inputs {
					inPortConfig.Name = inputInfo.Description["zh_CN"]
					inPortConfig.Key = inputName
					g.Components = append(g.Components, *inPortConfig)
				}
				outPortConfig := new(utils.Component)
				for _, v := range componentConfig {
					if v.Category == "outPorts" {
						outPortConfig = &v
						break
					}
				}
				for outputName, outputInfo := range g.NodeInfo.Outputs {
					outPortConfig.Name = outputInfo.Description["zh_CN"]
					outPortConfig.Key = outputName
					g.Components = append(g.Components, *outPortConfig)
				}
			} else {
				componentConfig := []utils.Component{}
				ymlFile, err := os.Open(path.Join("configs", f.Name()))
				if err != nil {
					log.Info(err.Error())
				}
				defer ymlFile.Close()
				byteValue, _ := io.ReadAll(ymlFile)
				yaml.Unmarshal(byteValue, &componentConfig)
				g.Components = append(g.Components, componentConfig...)
			}
		}
	}
}

var GraphInst Graph
