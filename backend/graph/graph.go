package graph

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goPipeline/components"
	"goPipeline/utils"
	"goPipeline/variables"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	"gopkg.in/yaml.v3"
)

type Graph struct {
	Status         uint // 0: edit 1: deploy
	PipelineStatus uint // 0: stop 1: running
	// ProcessMode	   uint // 1: 全部运行 2: 运行单个节点 3：停止运行
	Nodes          []components.Node
	Components     []utils.Component
	Config         utils.GraphConfig
	NodeInfo       utils.NodeInfo
	stopChan       chan bool
	wg             sync.WaitGroup
	path           string
	key            string
}

func (g *Graph) Init(appType string) {
	// 获取环境变量
	e := config.GetEnv()
	// 获取命令行参数
	args := config.GetArgs()
	g.path = path.Join(args["--storage-oss-temp-store"], "studio", e.SpUserId, "configs", e.SpAppId, e.SpNodeId, "graph.json")
	g.key = strings.Join([]string{"studio", e.SpUserId, "configs", e.SpAppId, e.SpNodeId, "graph.json"}, "/")
	g.componentsInit(appType)
	g.graphInit()
	g.nodesInit()
	variables.GlobalVariables = make(map[string]interface{})
}

func (g *Graph) graphInit() {
	err := storage.FGetObject(g.key, g.path)
	if err != nil {
		log.Info("Fail to Load Config File, init with default value...")
		g.Config = utils.GraphConfig{Scale: 1, Connectors: []utils.Connector{}, Nodes: []utils.NodeConfig{}}
	} else {
		jsonFile, err := os.Open(g.path)
		if err != nil {
			log.Info(err.Error())
			g.Config = utils.GraphConfig{}
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &g.Config)
		log.Info(fmt.Sprintf("Successfully Loaded Config File %s.", g.path))
	}
}

func (g *Graph) nodesInit() {
	for _, nodeConfig := range g.Config.Nodes {

		node := components.Node{Id: nodeConfig.Uuid, Key: nodeConfig.Key}
		if strings.HasPrefix(node.Key, "in") {
			node.Init("StreamIn")
		} else if strings.HasPrefix(node.Key, "out") {
			node.Init("StreamOut")
		} else {
			node.Init(nodeConfig.Key)
		}

		params := make(map[string]interface{})
		for _, param := range nodeConfig.Parameters {
			params[param.Key] = param.Value
		}
		if strings.HasPrefix(node.Key, "in") || strings.HasSuffix(node.Key, "out") {
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
		node.InputData = make(map[string]interface{})
		node.OutputData = make(map[string]interface{})
		supportPortConfig := []string{"ExecutePythonScript", "DataSync"}
		nodeSupportPortConfig := false
		for _, key := range supportPortConfig {
			if key == node.Key {
				nodeSupportPortConfig = true
			}
		}
		if nodeSupportPortConfig {
			for _, port := range node.Config["inPorts"].([]interface{}) {
				port := port.(map[string]interface{})
				node.InputData[port["id"].(string)] = nil
			}
			for _, port := range node.Config["outPorts"].([]interface{}) {
				port := port.(map[string]interface{})
				node.OutputData[port["id"].(string)] = nil
			}
		} else {
			for _, component := range g.Components {

				if component.Key == node.Key {
					for _, port := range component.Ports.In {
						node.InputData[port.Id] = nil
					}
					for _, port := range component.Ports.Out {
						node.OutputData[port.Id] = nil
					}
				}
			}
		}
		g.Nodes = append(g.Nodes, node)
	}
	for _, connection := range g.Config.Connectors {
		for i := range g.Nodes {
			node := &g.Nodes[i]
			if node.PortConnects == nil {
				node.PortConnects = make(map[string][]string)
			}
			if node.Id == connection.Src["uuid"] {
				if !g.checkNode(connection.Tgt["uuid"], node.NextNodes) {
					node.NextNodes = append(node.NextNodes, g.findNode(connection.Tgt["uuid"]))
				}
				if !utils.SlicesContain(node.PortConnects[connection.Src["port"]], connection.Tgt["uuid"]+"-"+connection.Tgt["port"]) {
					node.PortConnects[connection.Src["port"]] = append(node.PortConnects[connection.Src["port"]], connection.Tgt["uuid"]+"-"+connection.Tgt["port"])
				}
			}
			if node.Id == connection.Tgt["uuid"] {
				if !g.checkNode(connection.Src["uuid"], node.NextNodes) {
					node.PreviousNodes = append(node.PreviousNodes, g.findNode(connection.Tgt["uuid"]))
				}
			}
		}
	}
}

func (g *Graph) findNode(uuid string) *components.Node {
	for i := range g.Nodes {
		if g.Nodes[i].Id == uuid {
			return &g.Nodes[i]
		}
	}
	return nil
}

func (g *Graph) checkNode(uuid string, nodes []*components.Node) bool {
	for _, node := range nodes {
		if node.Id == uuid {
			return true
		}
	}
	return false
}

func (g *Graph) Update(newGraph utils.GraphConfig) {
	g.Config = newGraph
	os.Remove(g.path)
	dataJson, _ := json.Marshal(g.Config)
	os.WriteFile(g.path, dataJson, 0644)
	storage.FPutObject(g.key, g.path)
	g.Nodes = []components.Node{}
	g.nodesInit()
}

func (g *Graph) Run(inputData map[string]string, id string, extra string, server *socketio.Server, useCache bool) {
	log.Info("流程图开始运行")
	g.PipelineStatus = 1
	g.wg = sync.WaitGroup{}
	g.stopChan = make(chan bool)
	for _, node := range g.Nodes {
		if len(node.PreviousNodes) == 0 {
			if strings.HasPrefix(node.Key, "in") {
				if data, ok := inputData[strings.Replace(node.Key, "inputData", "in", -1)]; ok {
					g.wg.Add(1)
					go node.Run(node, components.RequestData{Data: data, ID: id, Extra: extra}, &g.wg, g.stopChan, server)
				} else {
					if useCache {
						g.wg.Add(1)
						go node.Run(node, components.RequestData{ID: id, Extra: extra}, &g.wg, g.stopChan, server)
					}
				}
			} else {
				if len(node.InputData) == 0 {
					g.wg.Add(1)
					go node.Run(node, components.RequestData{ID: id, Extra: extra}, &g.wg, g.stopChan, server)
				}
			}
		}
	}
	g.wg.Wait()
	g.PipelineStatus = 0
	close(g.stopChan)
	log.Info("流程图运行结束")
}

func (g *Graph) UpdateInputs(inputData map[string]string, id string, extra string) {
	log.Info("输入数据开始更新")
	g.wg = sync.WaitGroup{}
	g.stopChan = make(chan bool)
	for _, node := range g.Nodes {
		if len(node.PreviousNodes) == 0 {
			g.wg.Add(1)
			if strings.HasPrefix(node.Key, "in") {
				go node.UpdateInput(node, components.RequestData{Data: inputData[node.Key], ID: id, Extra: extra}, &g.wg, g.stopChan)
			}
		}
	}
	g.wg.Wait()
	close(g.stopChan)
	log.Info("输入数据结束更新")
}

func (g *Graph) Stop() {
	log.Info("Stop Graph")
	close(g.stopChan)
}

func (g *Graph) componentsInit(appType string) {
	files, err := os.ReadDir("configs")
	if err != nil {
		log.Error(err.Error())
	}
	componentsToLoad := make(map[string][]string)
	componentsToLoad["DataConnector"] = []string{"streamConnector.yml", "postgres.yml", "script.yml", "dataProcess.yml", "csv.yml"}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") && utils.SlicesContain(componentsToLoad[appType], f.Name()) {
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
					inPortConfig.Parameters = []utils.Parameter{}
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
					outPortConfig.Parameters = []utils.Parameter{}
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
