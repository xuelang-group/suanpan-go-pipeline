package services

import (
	"goPipeline/graph"
	"goPipeline/utils"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type Services struct {
	services []Service
}

func (h Services) Init() {
	log.Info("Init function not implement.")
}

func (h *Services) Update(newGraph utils.GraphConfig) {
	for _, service := range h.services {
		service.Release()
	}
	h.services = []Service{}
	for _, node := range newGraph.Nodes {
		if node.Key == "KafkaConsumer" {
			var address string
			var topic string
			var groupId string
			var topics string
			for _, param := range node.Parameters {
				switch param.Key {
				case "address":
					address = param.Value.(string)
				case "topic":
					topic = param.Value.(string)
				case "groupId":
					groupId = param.Value.(string)
				case "topics":
					topics = param.Value.(string)
				}
			}
			h.services = append(h.services, &KafkaService{Key: node.Key, Id: node.Uuid, Address: address, Topic: topic, GroupId: groupId, Topics: topics})
		}
	}
	for _, service := range h.services {
		service.Init()
	}
}

func (h Services) Deploy(g *graph.Graph) {
	for _, service := range h.services {
		go service.Deploy(g)
	}
}

func (h Services) Release() {
	for _, service := range h.services {
		service.Release()
	}
}

var ServicesManager Services
