package stream

import (
	"goPipeline/web"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type DataConnectorComponent struct {
	DefaultComponents
}

func (c *DataConnectorComponent) CallHandler(r stream.Request) {

}

func (c *DataConnectorComponent) InitHandler() {

}

func (c *DataConnectorComponent) SioHandler() {
	go web.RunWeb()
}
