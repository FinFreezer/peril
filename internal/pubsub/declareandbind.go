package pubsub

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	newChann, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	durable, autoDelete, exclusive, err := helperGetQueueparams(queueType)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	newTable := amqp.Table{"x-dead-letter-exchange": "peril_dlx"}
	newQueue, err := newChann.QueueDeclare(queueName, durable, autoDelete, exclusive, false, newTable)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = newChann.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return newChann, newQueue, nil
}

func helperGetQueueparams(queueType SimpleQueueType) (bool, bool, bool, error) {
	switch queueType {
	case Durable:
		return true, false, false, nil
	case Transient:
		return false, true, true, nil
	default:
		return false, false, false, errors.New("Problem converting queueType")
	}
}
