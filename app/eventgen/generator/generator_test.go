// Package generator generates event and publish to event topic
package generator

import (
	"context"
	"google/jss/up12/eventgen/config"
	"google/jss/up12/eventgen/generator/publishers"
	"google/jss/up12/pubsub"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	for i := 0; i < 10; i++ {
		event := NewEvent()
		log.Println("event=", event)
	}
}

func TestAddPublishers(t *testing.T) {
	ctx := context.Background()
	client, err := pubsub.Service.NewClient(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close() // nolint:errcheck

	topic := client.NewTopic(config.Config.EventTopic, config.Config.EventAvsc, 1, 0, 0)
	defer topic.Stop()

	publishers := publishers.NewPublishers(topic, NewEvent, 10*time.Second, 0, time.Second)
	defer publishers.Stop()

	publishers.Add(ctx, 2)
	time.Sleep(3 * time.Second)
	publishers.Add(ctx, -1)
	time.Sleep(5 * time.Second)
	publishers.Add(ctx, 1)
	publishers.WaitFinish()
}

func TestGeneratorTimeout(t *testing.T) {
	timeout := 2 * time.Second
	now := time.Now()
	Start(NewEvent, 2, timeout, 0, time.Second)
	running.publishers.WaitFinish()
	elapsed := time.Since(now).Seconds()
	log.Printf("elapsed: %v", elapsed)
	assert.Equal(t, int(timeout.Seconds()), int(elapsed))
	time.Sleep(1 * time.Second)
}

func TestGeneratorTimes(t *testing.T) {
	times := 3
	threads := 2
	Start(countNewEvent, threads, 0, times, time.Second)
	running.publishers.WaitFinish()
	assert.Equal(t, int32(times*threads), counter)
	time.Sleep(1 * time.Second)
}

var counter int32

func countNewEvent() map[string]interface{} {
	atomic.AddInt32(&counter, 1)
	return NewEvent()
}
