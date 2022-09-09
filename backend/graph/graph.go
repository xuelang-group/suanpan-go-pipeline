package graph

import (
	"encoding/base64"
	"encoding/json"
	"goPipeline/components"
	"goPipeline/utils"
	"goPipeline/web"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	"gopkg.in/yaml.v3"
)

type Graph struct {
	Status     bool
	Nodes      []components.Node
	Components []utils.Component
	Config     utils.GraphConfig
}

func (g *Graph) Init() {
	g.graphInit()
	g.loadComponents()
}

func (g *Graph) graphInit() {
	log.Info("Init function not implement.")
}

func (g *Graph) Update(newGraph utils.GraphConfig) {
	g.Config = newGraph
	os.Remove(web.DataPath)
	dataJson, _ := json.Marshal(g.Config)
	os.WriteFile(web.DataPath, dataJson, 0644)
	storage.FPutObject(web.DataKey, web.DataPath)
}

func (g Graph) Run() {
	log.Info("Init function not implement.")
}

func (g Graph) Stop() {
	log.Info("Init function not implement.")
}

func (g *Graph) loadComponents() {
	files, err := ioutil.ReadDir("configs")
	if err != nil {
		log.Error(err.Error())
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") {
			if f.Name() == "streamConnector.yml" {
				componentConfig := []utils.Component{}
				nodeInfo := utils.NodeInfo{}
				nodeInfoString, _ := base64.StdEncoding.DecodeString(os.Getenv("SP_NODE_INFO"))
				json.Unmarshal(nodeInfoString, &nodeInfo)
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
				for inputName, inputInfo := range nodeInfo.Inputs {
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
				for outputName, outputInfo := range nodeInfo.Outputs {
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
