package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](conn *amqp.Connection,
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

	err = amqpChan.Qos(10, 0, false)
	if err != nil {
		return err
	}

	deliveryChan, err := amqpChan.Consume(amqpQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		defer amqpChan.Close()
		for msg := range deliveryChan {
			data := bytes.NewBuffer(msg.Body)
			var newBuf T
			dec := gob.NewDecoder(data)
			err := dec.Decode(&newBuf)
			if err != nil {
				fmt.Println(err)
				return
			}
			ackType := handler(newBuf)
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
