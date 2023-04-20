// Package generator generates event and publish to event topic
package generator

import (
	"context"
	"log"
	"testing"
	"time"
	"up12/eventgen/config"
	"up12/pubsub/pubsub"
)

func TestNewEvent(t *testing.T) {
	for i := 0; i < 10; i++ {
		event := newEvent()
		log.Println("event=", event)
	}
}

func TestAddPublishers(t *testing.T) {
	msgChan := make(chan map[string]interface{})

	topic := pubsub.NewTopic(config.Config.EventTopicID, config.Config.EventAvsc, msgChan, 1)
	defer topic.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go slowGenerate(ctx, msgChan)
	topic.AddPublishers(ctx, 2)
	time.Sleep(10 * time.Second)
	topic.AddPublishers(ctx, -1)
	time.Sleep(10 * time.Second)
	topic.AddPublishers(ctx, 1)
	time.Sleep(10 * time.Second)
}

func slowGenerate(ctx context.Context, msgChan chan map[string]interface{}) {
	log.Println("Start to generate events!")
	for {
		select {
		case <-ctx.Done():
			log.Println("Time is up or cancnelled")
			return
		default:
			time.Sleep(1 * time.Second)
			msgChan <- newEvent()
		}
	}
}
