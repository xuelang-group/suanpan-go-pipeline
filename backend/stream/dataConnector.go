package stream

import (
	"goPipeline/graph"
	"goPipeline/web"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type DataConnectorComponent struct {
	DefaultComponents
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

}

func (c *DataConnectorComponent) SioHandler() {
	go web.RunWeb()
}
