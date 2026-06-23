package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	lg "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	rt "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const moveExchangeKey = "army_moves."
	fmt.Println("Starting Peril client...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-c
		fmt.Println()
		fmt.Println("Peril client closed...")
		os.Exit(0)
	}()

	connStr := "amqp://guest:guest@localhost:5672/"
	amqpConn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer amqpConn.Close()
	fmt.Println("Peril client succesfully started.")
	usrName, err := lg.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		return
	}
	newGame := lg.NewGameState(usrName)
	cliChann, _, err := pubsub.DeclareAndBind(amqpConn, rt.ExchangePerilTopic, moveExchangeKey+usrName, moveExchangeKey+"*", pb.Transient)
	if err != nil {
		fmt.Println(err)
		return
	}
	queueName := rt.PauseKey + "." + usrName
	err = pubsub.SubscribeJSON(amqpConn, rt.ExchangePerilDirect, queueName, rt.PauseKey, pb.Transient, handlerPause(newGame))
	if err != nil {
		fmt.Println(err)
		return
	}
	queueName = moveExchangeKey + usrName
	err = pb.SubscribeJSON(amqpConn, rt.ExchangePerilTopic, queueName, moveExchangeKey+"*", pb.Transient, handlerArmyMove(newGame, cliChann))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pb.SubscribeJSON(amqpConn, rt.ExchangePerilTopic, "war", rt.WarRecognitionsPrefix+".*", pb.Durable, handlerConsumeWarMessage(newGame, cliChann))
	if err != nil {
		fmt.Println(err)
		return
	}
	/*queueName := rt.PauseKey + "." + usrName

	_, _, err = pb.DeclareAndBind(amqpConn, rt.ExchangePerilDirect, queueName, rt.PauseKey, pb.Transient)
	if err != nil {
		fmt.Println(err)
		return
	}*/

gameLoop:
	for {
		cmds := lg.GetInput()
		switch cmds[0] {
		case "spawn":
			err := newGame.CommandSpawn(cmds)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			armyMove, err := newGame.CommandMove(cmds)
			if err != nil {
				fmt.Println(err)
			}
			err = pb.PublishJSON(cliChann, rt.ExchangePerilTopic, moveExchangeKey, armyMove)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Message published succesfully.")
		case "status":
			newGame.CommandStatus()
		case "help":
			lg.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			lg.PrintQuit()
			break gameLoop
		default:
			fmt.Println("Unknown command...")
		}
	}
	fmt.Println()
	fmt.Println("Peril client closed...")
}
