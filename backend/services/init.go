package services

import (
	"goPipeline/graph"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type Service interface {
	Init()
	Deploy(g *graph.Graph)
	Release()
}

type DefaultService struct {
}

func (h DefaultService) Init() {
	log.Info("Init function not implement.")
}

func (h DefaultService) Deploy(g *graph.Graph) {
	log.Info("Deploy function not implement.")
}

func (h DefaultService) Release() {
	log.Info("Release function not implement.")
}
