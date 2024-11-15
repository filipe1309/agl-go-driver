package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Duration
}

func newRabbitMQConnection(cfg RabbitMQConfig) (rc *RabbitMQConnection, err error) {
	rc.cfg = cfg
	rc.conn, err = amqp.Dial(cfg.URL)
	return rc, err
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

func (rc *RabbitMQConnection) Consume(chanDTO chan<- QueueDTO) error {
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
		dto := QueueDTO{}
		dto.Unmarshal(deliveryChan.Body)
		chanDTO <- dto
	}

	return nil
}
