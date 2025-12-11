package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	gameState, ch := ClientStartup()

	ListenForCommands(gameState, ch)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Console is shutting down")
}

