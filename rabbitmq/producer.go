package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Producer interface {
	ProduceMessages(queue string, message []byte) error
	Close() error
}

type rabbitMQProducer struct {
	channel *amqp.Channel
}

func NewRabbitMQProducer(channel *amqp.Channel) Producer {
	return &rabbitMQProducer{channel: channel}
}

func (p *rabbitMQProducer) ProduceMessages(queue string, message []byte) error {
	err := p.channel.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *rabbitMQProducer) Close() error {
	return p.channel.Close()
}
