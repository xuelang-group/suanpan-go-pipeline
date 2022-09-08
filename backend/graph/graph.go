package graph

import (
	"goPipeline/components"
	"goPipeline/utils"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"gopkg.in/yaml.v3"
)

type Graph struct {
	Status     bool
	Nodes      []components.Node
	Components []utils.Component
}

func (g Graph) Init() {
	log.Info("Init function not implement.")
}

func (g Graph) Update() {
	log.Info("Init function not implement.")
}

func (g Graph) Get() {
	log.Info("Init function not implement.")
}

func (g Graph) Run() {
	log.Info("Init function not implement.")
}

func (g Graph) Stop() {
	log.Info("Init function not implement.")
}

func (g *Graph) LoadComponents() {
	files, err := ioutil.ReadDir("configs")
	if err != nil {
		log.Error(err.Error())
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") {
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

var GraphInst Graph
