package components

import "github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"

type ProcessNode struct {
	Node
}

func (c *ProcessNode) Init() {
	log.Info("Init function not implement.")
}

func (c *ProcessNode) BeforeExit() {
	log.Info("BeforeExit function not implement.")
}

func (c *ProcessNode) Run() {
	log.Info("Run function not implement.")
}

func (c *ProcessNode) LoadInput() {
	log.Info("LoadInput function not implement.")
}

func (c *ProcessNode) DumpOutput() {
	log.Info("DumpOutput function not implement.")
}

func (c *ProcessNode) Send() {
	log.Info("Send function not implement.")
}
