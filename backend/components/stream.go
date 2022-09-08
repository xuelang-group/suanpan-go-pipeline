package components

import "github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"

type StreamNode struct {
	Node
}

func (c *StreamNode) Init() {
	log.Info("Init function not implement.")
}

func (c *StreamNode) BeforeExit() {
	log.Info("BeforeExit function not implement.")
}

func (c *StreamNode) Run() {
	c.DumpOutput()
}

func (c *StreamNode) LoadInput() {
	log.Info("LoadInput function not implement.")
}

func (c *StreamNode) DumpOutput() {
	log.Info("DumpOutput function not implement.")
}

func (c *StreamNode) Send() {
	log.Info("Send function not implement.")
}
