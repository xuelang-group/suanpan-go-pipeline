package stream

import (
	"goPipeline/graph"
	"goPipeline/services"
	"goPipeline/web"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type DataConnectorComponent struct {
	DefaultComponents
	Type string
	Mode string
}

func (c *DataConnectorComponent) CallHandler(r stream.Request) {
	inputData := r.Input
	for key, value := range inputData {
		log.Infof("输入端口： %s 收到请求数据： %s", key, value)
	}
	if graph.GraphInst.Status == 1 {
		graph.GraphInst.Run(inputData, r.ID, r.Extra, nil, false)
	} else {
		graph.GraphInst.UpdateInputs(inputData, r.ID, r.Extra)
	}
}

func (c *DataConnectorComponent) InitHandler() {
	graph.GraphInst.Init(c.Type, c.Mode)
	services.ServicesManager.Update(graph.GraphInst.Config)
}

func (c *DataConnectorComponent) SioHandler() {
	if c.Mode == "edit" {
		go web.RunWeb(c.Type)
	}
}
