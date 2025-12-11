package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ppitek40/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ppitek40/learn-pub-sub-starter/internal/routing"
	"github.com/ppitek40/learn-pub-sub-starter/internal/pubsub"
)

func ClientStartup() (*gamelogic.GameState, *amqp.Channel) {
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		panic("Connection couldnt be opened")
	}
	defer connection.Close()

	ch, err := connection.Channel()
	if err != nil {
		panic("Error with creating channel")
	}

	name, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("%w", err)
		panic("Error")
	}

	gameState := gamelogic.NewGameState(name)
	pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix+"."+name, pubsub.Durable, handlerWar(gameState, ch))
	pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, routing.PauseKey+"."+name, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+name, routing.ArmyMovesPrefix+".*", pubsub.Transient, handlerMove(gameState, ch))
	return gameState, ch
}