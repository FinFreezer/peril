package pubsub

import (
	"fmt"

	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	rt "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(msg string, row gamelogic.RecognitionOfWar, chn *amqp.Channel) AckType {
	newLog := rt.GameLog{CurrentTime: time.Now(), Message: msg, Username: row.Attacker.Username}
	exchange := rt.ExchangePerilTopic
	rtKey := rt.GameLogSlug + "." + row.Attacker.Username
	err := PublishGob(chn, exchange, rtKey, newLog)
	if err != nil {
		fmt.Println(err)
		return NackRequeue
	}
	return Ack
}
