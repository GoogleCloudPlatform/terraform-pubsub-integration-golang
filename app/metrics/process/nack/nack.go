package nack

import (
	"context"
	"log"
	"time"
	"up12/metrics/config"
	"up12/metrics/process"
	"up12/pubsub/pubsub"
)

type EventProcessor struct {
}

func (*EventProcessor) AvscFile() string {
	return config.Config.MetricsAckAvsc
}

func (*EventProcessor) EventHandler() pubsub.MessageHandler {
	return func(ctx context.Context, message *pubsub.Message) {
		log.Printf("Processing event=%v", message.Data)

		processingTime := process.ProcessingTime()
		time.Sleep(processingTime) // Simulate processing time

		message.Nack() // simulate a bug 🐞 and nack the message
	}

}
