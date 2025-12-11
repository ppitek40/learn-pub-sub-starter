package main

import (
	"fmt"
	"os/signal"
	"os"
	"github.com/ppitek40/learn-pub-sub-starter/internal/gamelogic"
)

func main() {
	ch := ServerStartup()

	gamelogic.PrintServerHelp()
	
	ListenForCommands(ch)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Console is shutting down")
}
