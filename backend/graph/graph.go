package main

import (
	"goPipeline/components"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type Graph struct {
	Status bool
	Nodes  []components.Node
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
