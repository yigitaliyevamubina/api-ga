package consumer

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	ConsumeMessages(handler func(message []byte)) error
	Close() error
}

type Consumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) (KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: reader}, nil
}

func (c *Consumer) ConsumeMessages(handler func(message []byte)) error {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		handler(msg.Value)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
