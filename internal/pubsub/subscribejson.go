package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType string

const (
	Ack         AckType = "Ack"         //.Ack(false)
	NackRequeue AckType = "NackRequeue" //.Nack(false, true)
	NackDiscard AckType = "NackDiscard" //.Nack(false, false)
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	amqpChan, amqpQueue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	deliveryChan, err := amqpChan.Consume(amqpQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range deliveryChan {
			var bytes T
			err := json.Unmarshal(msg.Body, &bytes)
			if err != nil {
				fmt.Println(err)
				return
			}
			ackType := handler(bytes)
			switch ackType {
			case Ack:
				fmt.Println("Calling Acknowledge")
				msg.Ack(false)
			case NackRequeue:
				fmt.Println("Calling Requeue")
				msg.Nack(false, true)
			case NackDiscard:
				fmt.Println("Calling Discard")
				msg.Nack(false, false)
			default:
				fmt.Println("Error processing Delivery's 'Ack' status.")
			}
		}
	}()
	return nil
}
