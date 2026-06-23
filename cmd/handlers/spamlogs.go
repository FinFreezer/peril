package handlers

import (
	"time"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	rt "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HelperSpamLogs(chn *amqp.Channel, exchange, key, userName string, n int) {
	for i := 0; i < n; i++ {
		newBadLog := gamelogic.GetMaliciousLog()
		newBadLogStruct := rt.GameLog{CurrentTime: time.Now(), Message: newBadLog, Username: userName}
		pb.PublishGob(chn, exchange, key, newBadLogStruct)
	}
}
