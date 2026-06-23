package main

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerArmyMove(gs *gamelogic.GameState) func(st gamelogic.ArmyMove) pb.AckType {
	return func(st gamelogic.ArmyMove) pb.AckType {
		defer fmt.Print("> ")
		mvOutcome := gs.HandleMove(st)

		switch mvOutcome {
		case gamelogic.MoveOutComeSafe:
			return pb.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pb.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pb.NackDiscard
		default:
			return pb.NackDiscard
		}
	}
}
