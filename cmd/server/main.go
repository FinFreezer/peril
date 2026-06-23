package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	lg "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pb "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	rt "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		fmt.Println()
		fmt.Println("Peril server closed...")
		os.Exit(0)
	}()

	connStr := "amqp://guest:guest@localhost:5672/"
	amqpConn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer amqpConn.Close()

	queueName := rt.GameLogSlug
	routingKey := rt.GameLogSlug + "." + "*"
	err = pb.SubscribeGob(amqpConn, rt.ExchangePerilTopic, queueName, routingKey, pb.Durable, handlerWriteLogs())
	if err != nil {
		fmt.Println(err)
		return
	}

	gameChann, err := amqpConn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gameChann.Close()

	exch := rt.ExchangePerilDirect
	pauseKey := rt.PauseKey
	//val := rt.PlayingState{IsPaused: true}
	//err = pb.PublishJSON(gameChann, exch, pauseKey, val)

	/*if err != nil {
		fmt.Println(err)
		return
	}*/
	fmt.Println("Peril server succesfully started.")
	lg.PrintServerHelp()

serverLoop:
	for {
		cmds := lg.GetInput()
		switch cmds[0] {
		case "pause":
			log.Println("Sending pause message...")
			err := pb.PublishJSON(gameChann, exch, pauseKey, rt.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Println(err)
			}
		case "resume":
			err = pb.PublishJSON(gameChann, exch, pauseKey, rt.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Println(err)
			}
		case "quit":
			log.Println("Exiting server...")
			break serverLoop
		case "help":
			lg.PrintServerHelp()
		default:
			log.Println("Unknown command.")
		}
	}

	fmt.Println()
	fmt.Println("Program shutting down...")
	fmt.Println("Connection closing...")
}
