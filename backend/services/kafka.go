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
	Key       string
	Id        string
	Address   string
	Topic     string
	Partition string
	// Graph     graph.Graph
	// IsDeploy bool
	StopChan chan bool
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func (h *KafkaService) Deploy(g *graph.Graph) {
	select {
	case <-h.StopChan:
		log.Infof("节点%s(%s)停止消费 address:%s topic:%s partition: %s 中的消息", h.Key, h.Id, h.Address, h.Topic, h.Partition)
	default:
		// get kafka reader using environment variables.
		kafkaURL := h.Address
		topic := h.Topic
		groupID := h.Partition

		reader := getKafkaReader(kafkaURL, topic, groupID)

		defer reader.Close()

		log.Infof("节点%s(%s)开始消费 address:%s topic:%s partition: %s 中的消息", h.Key, h.Id, kafkaURL, topic, groupID)
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Errorf("读取Kafka消息失败: %s", err)
			}
			inputData := map[string]string{h.Id: string(m.Value)}
			id := util.GenerateUUID()
			extra := ""
			g.Run(inputData, id, extra, nil, false)
			log.Infof("在 topic:%v partition:%v offset:%v 中的消息 %s = %s 消费成功", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
	}
	// graph.GraphInst.Run(currentNode.NextNodes[0], inputData)
}

func (h *KafkaService) Release() {
	close(h.StopChan)
}

func (h *KafkaService) Init() {

}
