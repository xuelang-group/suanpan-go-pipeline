package components

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func kafkaProducer(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	topic := currentNode.Config["topic"].(string)
	partition := currentNode.Config["partition"].(int)
	address := currentNode.Config["address"].(string)
	// 连接至Kafka集群的Leader节点
	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		log.Fatal("与Kafka主节点连接失败:", err)
	}

	// 设置发送消息的超时时间
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 发送消息
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(currentNode.InputData["in1"].(string))},
	)
	if err != nil {
		log.Fatal("Kafka写入消息失败:", err)
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("Kafka写入关闭失败:", err)
	}
	return map[string]interface{}{}, nil
}

func kafkaConsumerMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	return map[string]interface{}{"out1": inputData.Data}, nil
}
