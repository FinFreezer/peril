package handlers

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	rt "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandlerArmyMove(gs *gamelogic.GameState, chn *amqp.Channel) func(st gamelogic.ArmyMove) pb.AckType {
	return func(st gamelogic.ArmyMove) pb.AckType {
		defer fmt.Print("> ")
		mvOutcome := gs.HandleMove(st)

		switch mvOutcome {
		case gamelogic.MoveOutcomeSafe:
			return pb.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pb.PublishJSON(
				chn,
				rt.ExchangePerilTopic,
				rt.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: st.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				return pb.NackRequeue
			}
			return pb.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pb.NackDiscard
		default:
			return pb.NackDiscard
		}
	}
}
