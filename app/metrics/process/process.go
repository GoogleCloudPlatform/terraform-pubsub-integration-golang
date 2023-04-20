package process

import (
	"math/rand"
	"time"
	"up12/pubsub/pubsub"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type EventProcessor interface {
	AvscFile() string
	EventHandler() pubsub.MessageHandler
}

func ProcessingTime() time.Duration {
	min := 0.1
	max := 5.0
	mean := (max + min) / 2
	stdDev := (max - mean) / 3
	seconds := random.NormFloat64()*stdDev + mean
	return time.Duration(seconds * float64(time.Second))
}
