package messagequeue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

//Producer structure contains the producer object
type Producer struct {
	writer *kafka.Writer
}

//Message structure of traffic information
type Message struct {
	Domain        string    `json:"domain"`
	RemoteAddress string    `json:"remoteAddress"`
	Timestamp     time.Time `json:"timestamp"`
	Blocked       bool      `json:"blocked"`
}

//CreateProducer method returns a struct that allows sending messages to kafka
func CreateProducer(bootstapServers []string, topic string) (*Producer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  bootstapServers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &Producer{writer: writer}, nil
}

//Send method is used to send a message to a kafka topic
func (kafkaProducer *Producer) Send(message *Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return kafkaProducer.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: messageBytes,
		})
}
