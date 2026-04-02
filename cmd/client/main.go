package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp091.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	_, err = conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get user name: %v", err)
	}

	gs := gamelogic.NewGameState(username)

	pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", username), routing.PauseKey, pubsub.SimpleQueueType{
		Durable:   false,
		Transient: true,
	}, handlerPause(gs))

	for {
		cmd := gamelogic.GetInput()
		if len(cmd) == 0 {
			continue
		}

		switch cmd[0] {
		case "spawn":
			err = gs.CommandSpawn(cmd)
			if err != nil {
				log.Printf("%v", err)
				continue
			}

		case "move":
			_, err := gs.CommandMove(cmd)
			if err != nil {
				log.Printf("%v", err)
				continue
			}

		case "status":
			gs.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			log.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			log.Println("Command not recognized")
			continue
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(routing.PlayingState{})
	}
}
