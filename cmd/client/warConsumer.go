package main

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerConsumeWarMessage(gs *gamelogic.GameState, chn *amqp.Channel) func(row gamelogic.RecognitionOfWar) pb.AckType {
	return func(row gamelogic.RecognitionOfWar) pb.AckType {
		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(row)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pb.NackRequeue

		case gamelogic.WarOutcomeNoUnits:
			return pb.NackDiscard

		case gamelogic.WarOutcomeOpponentWon:
			logMessage := fmt.Sprintf("%s won a war against %s", winner, loser)
			ackType := pb.PublishGameLog(logMessage, row, chn)
			return ackType

		case gamelogic.WarOutcomeYouWon:
			logMessage := fmt.Sprintf("%s won a war against %s", winner, loser)
			ackType := pb.PublishGameLog(logMessage, row, chn)
			return ackType

		case gamelogic.WarOutcomeDraw:
			logMessage := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			ackType := pb.PublishGameLog(logMessage, row, chn)
			return ackType

		default:
			fmt.Println("Unknown outcome")
			return pb.NackDiscard
		}
	}
}
