package event_test

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/segmentio/kafka-go"
)

// 创建Topic
func TestCreateTopic(t *testing.T) {
	// 1. 连上kafka
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// contoller 管理, 获取zk地址后，管理集群状态
	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	// contoller 集群的维护
	err = controllerConn.CreateTopics(kafka.TopicConfig{Topic: "maudit", NumPartitions: 6, ReplicationFactor: 1})
	if err != nil {
		t.Fatal(err)
	}
}

// 查询Topic列表
func TestListTopic(t *testing.T) {
	// 1. 连上kafka
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		t.Fatal(err)
	}

	topics := map[string]int{}
	for _, p := range partitions {
		topics[p.Topic]++
	}
	t.Log(topics)
}

func TestKafkaWirteMessage(t *testing.T) {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "maudit",
		Balancer: &kafka.LeastBytes{},
	}
	err := w.WriteMessages(context.Background(), kafka.Message{
		// 支持 Writing to multiple topics
		//  NOTE: Each Message has Topic defined, otherwise an error is returned.
		// Topic: "topic-A",
		Key:   []byte("Key-A"),
		Value: []byte("Hello World!"),
	},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		})
	if err != nil {
		t.Fatal(err)
	}
}

func TestKafkaReadMessage(t *testing.T) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "maudit",
		GroupID: "maudit",
	})
	//  自动化 1. 读取消息, 读出来 就代表已经被处理,  FetchMessage, Commit(OK)
	// r.ReadMessage(context.Background())
	//  手动操作: 2. 获取消息, commit(OK), 对消息可靠度有要求，自己严格控制，避免消息丢失
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 标记消息已处理
		r.CommitMessages(context.Background(), m)
	}
}
