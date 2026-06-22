package main

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(st routing.PlayingState) {
	return func(st routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(st)
	}
}
