package main

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func handlerArmyMove(gs *gamelogic.GameState) func(st gamelogic.ArmyMove) {
	return func(st gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(st)
	}
}
