package main

import (
	"fmt"
	"strconv"
	"time"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ppitek40/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ppitek40/learn-pub-sub-starter/internal/pubsub"
	"github.com/ppitek40/learn-pub-sub-starter/internal/routing"
)

func ListenForCommands(gameState *gamelogic.GameState, ch *amqp.Channel){
	name := gameState.GetUsername()

	for {
		input := gamelogic.GetInput()
		if input == nil || len(input) == 0 {
			continue
		}
		
		switch input[0]{
		case "spawn":
			fmt.Println("Spawning unit")
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println("Error commanding move")
			}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+name, move)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Army move published successfully")
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(input) != 2 {
				fmt.Println("Wrong format use: spam 10")
				continue
			}
			result, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Errorf("%w", err)
			}
			for i := 0; i < result; i++ {
				log := gamelogic.GetMaliciousLog()
				publishGameLog(log, gamelogic.RecognitionOfWar{Attacker: gamelogic.Player{Username: name}}, ch)
				//pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+name, log)
			}
		case "quit":
			gamelogic.PrintQuit()
			break;
		default:
			fmt.Println("Unrecognized command")
		}
	}
}

func publishGameLog(mess string, war gamelogic.RecognitionOfWar, ch *amqp.Channel) error {
	gameLog := routing.GameLog{ CurrentTime: time.Now(), Message: mess, Username: war.Attacker.Username}
	err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+war.Attacker.Username, gameLog)
	return err
}