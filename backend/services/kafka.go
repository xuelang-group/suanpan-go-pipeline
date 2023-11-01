package services

import (
	"context"
	"goPipeline/graph"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

type KafkaService struct {
	DefaultService
	Key     string
	Id      string
	Address string
	Topic   string
	GroupId string
	Topics  string
	// Graph     graph.Graph
	// IsDeploy bool
	// StopChan    chan bool
	kafkaReader *kafka.Reader
}

func getKafkaReader(kafkaURL string, topic string, groupId string, topics string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	if len(topics) == 0 {
		return kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
			// Partition: partition,
			GroupID: groupId,
		})
	} else {
		groupTopics := strings.Split(topics, ",")
		return kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			// Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
			// Partition: partition,
			GroupID:     groupId,
			GroupTopics: groupTopics,
		})
	}

}

func (h *KafkaService) Deploy(g *graph.Graph) {
	// get kafka reader using environment variables.
	kafkaURL := h.Address
	topic := h.Topic
	groupId := h.GroupId
	groupTopics := h.Topics

	h.kafkaReader = getKafkaReader(kafkaURL, topic, groupId, groupTopics)

	// defer h.kafkaReader.Close()

	log.Infof("节点%s(%s)开始消费 address:%s groupid: %s 中的消息", h.Key, h.Id, kafkaURL, groupId)
	for {
		// m, err := h.kafkaReader.ReadMessage(context.Background())
		m, err := h.kafkaReader.FetchMessage(context.Background())
		if err != nil {
			log.Errorf("读取Kafka消息失败: %s", err)
			break
		}
		success := false
		inputData := map[string]string{h.Id: string(m.Value)}
		id := util.GenerateUUID()
		extra := ""
		workflowErr := g.Run(inputData, id, extra, nil, false)
		if workflowErr == nil {
			commitErr := h.kafkaReader.CommitMessages(context.Background(), m)
			if commitErr != nil {
				log.Errorf("提交Kafka消息失败: %s", err)
			} else {
				success = true
				log.Infof("在 topic:%v partition:%v offset:%v 中的消息 %s = %s 消费成功", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			}
		} else {
			for i := 0; i < 1; i++ {
				log.Infof("再次尝试消费消息 %s = %s ", string(m.Key), string(m.Value))
				inputData := map[string]string{h.Id: string(m.Value)}
				id := util.GenerateUUID()
				extra := ""
				workflowErr := g.Run(inputData, id, extra, nil, false)
				if workflowErr == nil {
					commitErr := h.kafkaReader.CommitMessages(context.Background(), m)
					if commitErr != nil {
						log.Errorf("提交Kafka消息失败: %s", err)
					} else {
						success = true
						log.Infof("在 topic:%v partition:%v offset:%v 中的消息 %s = %s 消费成功", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
						break
					}
				}
			}
		}
		if !success {
			log.Infof("在 topic:%v partition:%v offset:%v 中的消息 %s = %s 消费失败，消息订阅程序终止", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			break
		}
	}
	log.Infof("节点%s(%s)停止消费 address:%s groupid: %s 中的消息", h.Key, h.Id, kafkaURL, groupId)
	// graph.GraphInst.Run(currentNode.NextNodes[0], inputData)
}

func (h *KafkaService) Release() {
	if h.kafkaReader != nil {
		h.kafkaReader.Close()
		h.kafkaReader = nil
	}
	// h.kafkaReader.Close()
}

func (h *KafkaService) Init() {

}
