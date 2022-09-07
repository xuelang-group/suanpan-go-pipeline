package components

import "goPipeline/utils"

type NodeAction interface {
	Init()
	BeforeExit()
	Run()
	LoadInput(data utils.Data)
	DumpOutput(data utils.Data)
	Send([]Node)
}

type Node struct {
	PreviousNode *Node
	NextNode     *Node
	InputData    utils.Data
	OutputData   utils.Data
	Config       map[string]interface{}
	Id           string
}
