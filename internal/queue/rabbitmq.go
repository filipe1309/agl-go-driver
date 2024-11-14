package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Time
}

type RabbitMQConnection struct {
	cfg  RabbitMQConfig
	conn *amqp.Connection
}

func (rc *RabbitMQConnection) Publish(msg []byte) error {
	conn, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return conn.PublishWithContext(ctx, "", rc.cfg.TopicName, false, false, message)
}

func (rc *RabbitMQConnection) Consume() error {
	conn, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	q, err := conn.QueueDeclare(rc.cfg.TopicName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := conn.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for deliveryChan := range msgs {
		// Do something with the message
		_ = deliveryChan
	}

	return nil
}
