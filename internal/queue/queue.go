package queue

import (
	"fmt"
	"log"
	"reflect"
)

const (
	RabbitMQ QueueType = iota
	MockQueueProvider
)

type QueueType int

func New(qt QueueType, cfg any) (q *Queue, err error) {
	q = new(Queue)

	rt := reflect.TypeOf(cfg)

	switch qt {
	case RabbitMQ:
		if rt.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("invalid rabbitmq config type: %s", rt.Name())
		}
		conn, err := newRabbitMQConnection(cfg.(RabbitMQConfig))
		if err != nil {
			return nil, err
		}
		q.qc = conn
	case MockQueueProvider:
		q.qc = &MockQueue{make([]*QueueDTO, 0)}
	default:
		log.Fatal("type not supported")
	}
	return
}

type QueueConnection interface {
	Publish([]byte) error
	Consume(chan<- QueueDTO) error
}

type Queue struct {
	qc  QueueConnection
}


func (q *Queue) Publish(msg []byte) error {
	return q.qc.Publish(msg)
}

func (q *Queue) Consume(chanDTO chan<- QueueDTO) error {
	return q.qc.Consume(chanDTO)
}
