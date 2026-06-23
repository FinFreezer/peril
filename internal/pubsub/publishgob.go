package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var newBuf bytes.Buffer
	enc := gob.NewEncoder(&newBuf)
	err := enc.Encode(val)
	if err != nil {
		return err
	}
	newPub := amqp.Publishing{ContentType: "application/gob", Body: newBuf.Bytes()}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, newPub)
	if err != nil {
		return err
	}
	return nil
}
