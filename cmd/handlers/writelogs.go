package handlers

import (
	"fmt"

	gl "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func HandlerWriteLogs() func(newLog routing.GameLog) pb.AckType {
	return func(newLog routing.GameLog) pb.AckType {
		defer fmt.Print("> ")
		err := gl.WriteLog(newLog)
		if err != nil {
			fmt.Println(err)
			return pb.NackRequeue
		}
		return pb.Ack
	}
}
