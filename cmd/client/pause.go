package main

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(st routing.PlayingState) pb.AckType {
	return func(st routing.PlayingState) pb.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(st)
		return pb.Ack
	}
}
