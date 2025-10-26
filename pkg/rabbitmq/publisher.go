package rabbitmq

import (
	"context"
	"fmt"

	"github.com/fahri/go-tije/internal/config"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
}

func NewPublisher(cfg *config.RabbitMQConfig) (*Publisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}
	
	err = channel.ExchangeDeclare(
		cfg.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}
	
	_, err = channel.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}
	
	err = channel.QueueBind(
		cfg.Queue,
		"geofence.#",
		cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}
	
	return &Publisher{
		conn:     conn,
		channel:  channel,
		exchange: cfg.Exchange,
		queue:    cfg.Queue,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, body []byte) error {
	return p.channel.Publish(
		p.exchange,
		"geofence.alert",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}