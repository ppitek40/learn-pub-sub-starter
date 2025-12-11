package main

import (
	"fmt"
	"github.com/ppitek40/learn-pub-sub-starter/internal/pubsub"
	"github.com/ppitek40/learn-pub-sub-starter/internal/gamelogic"
	"github.com/ppitek40/learn-pub-sub-starter/internal/routing"

)

func handlerLog() func(routing.GameLog) pubsub.AckType {
	return func (log routing.GameLog) pubsub.AckType {
		defer fmt.Println("> ")
		if err := gamelogic.WriteLog(log); err != nil {
			return pubsub.NackDiscard
		}
		return pubsub.Ack
	}
}