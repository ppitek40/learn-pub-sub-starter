package main


import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ppitek40/learn-pub-sub-starter/internal/pubsub"
	"github.com/ppitek40/learn-pub-sub-starter/internal/routing"
)

func ServerStartup() *amqp.Channel {
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		panic("Connection couldnt be opened")
	}
	defer connection.Close()

	ch, err := connection.Channel()
	if err != nil {
		panic("Error creating channel")
	}

	
	fmt.Println("Starting Peril server...")
	err = pubsub.SubscribeGob(connection, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable, handlerLog())
	if err != nil {
		fmt.Errorf("Error: %w", err)
		panic("Error with subscribing to logs")
	}

	return ch
}