// Package main the entry point of metrics
package main

import (
	"context"
	"log"
	"time"
	"up12/metrics/config"
	"up12/metrics/process"
	"up12/metrics/process/ack"
	"up12/metrics/process/nack"
	"up12/pubsub/pubsub"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgChnn := make(chan map[string]interface{})
	topic := pubsub.NewTopic(config.Config.MetricsTopicID, config.Config.MetricsAckAvsc, msgChnn, config.Config.BatchSize)
	topic.AddPublishers(ctx, 1)

	sub := pubsub.NewSubscription(config.Config.SubscriptionID, config.Config.EventAvsc)

	processor := GetProcessor(msgChnn)

	if err := sub.Receive(ctx, processor.EventHandler()); err != nil {
		log.Fatal(err)
	}

}

func GetProcessor(msgChnn chan map[string]interface{}) process.EventProcessor {
	switch config.Config.MetricsAppType {
	case config.MetricsAck:
		return &ack.EventProcessor{MsgChnn: msgChnn}
	case config.MetricsNack:
		return &nack.EventProcessor{}
	// case MetricsComplete:
	default:
		log.Panicf("Unknown metrics app type=")
		return nil
	}
}
