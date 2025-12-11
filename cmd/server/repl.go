package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ppitek40/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ppitek40/learn-pub-sub-starter/internal/pubsub"
	"github.com/ppitek40/learn-pub-sub-starter/internal/routing"
)

func ListenForCommands(ch *amqp.Channel){
	for {
		input := gamelogic.GetInput()
		if input == nil || len(input) == 0 {
			continue
		}
		
		switch input[0]{
		case "pause":
			fmt.Println("Pausing game")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{ IsPaused: true})
		case "resume":
			fmt.Println("Resuming game")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{ IsPaused: false})
		case "quit":
			fmt.Println("Quitting game")
			return;
		default:
			fmt.Println("Unrecognized command")
		}
	}
}