// Package generator generates event and publish to event topic
package generator

import (
	"context"
	"errors"
	"google/jss/up12/eventgen/config"
	"google/jss/up12/eventgen/generator/publishers"
	"google/jss/up12/pubsub"
	"log"
	"sync"
	"time"

	"github.com/linkedin/goavro/v2"
)

type generator struct {
	client     pubsub.Client
	topic      pubsub.Topic
	publishers *publishers.Publishers
	cancel     context.CancelFunc
}

func newGenerator(topicID string, codec *goavro.Codec, batchSize int, maxOutstanding int, numGoroutines int) (*generator, error) {
	var g generator

	client, err := pubsub.Service.NewClient(context.Background())
	if err != nil {
		log.Printf("fail to connect to pubsub, err: %v", err)
		return nil, err
	}
	g.client = client

	g.topic = client.NewTopic(topicID, codec, batchSize, maxOutstanding, numGoroutines)
	return &g, nil
}

func (g *generator) Run(event publishers.NewMessage, numPublishers int, timeout time.Duration, times int, sleep time.Duration) {
	log.Printf("run event generator with numPublishers: %v, timeout: %v, times: %v, sleep: %v", numPublishers, timeout, times, sleep)
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel

	pbrs := publishers.NewPublishers(g.topic, event, timeout, times, sleep)
	pbrs.Add(ctx, numPublishers)
	g.publishers = pbrs

	go func() {
		pbrs.WaitFinish()
		g.release()
	}()
}

// Stop stops the generator gracefully
func (g *generator) Stop() {
	if g.publishers != nil {
		g.publishers.Stop() // Resources will be released after all publishers have finished
		// g.cancel()	TBD force stop
	} else {
		g.release()
	}
}

func (g *generator) release() {
	g.topic.Stop()
	g.client.Close()
	mux.Lock()
	defer mux.Unlock()
	if running == g {
		running = nil
	}
}

var mux sync.Mutex     // Protect running singleon
var running *generator // Singleton, only one generator at the same time.

// Start generates event and publish to event Topic
func Start(event publishers.NewMessage, numPublishers int, timeout time.Duration, times int, sleep time.Duration) error {
	mux.Lock()
	defer mux.Unlock()

	if running != nil {
		return errors.New("there is already an running generator")
	}
	g, err := newGenerator(config.Config.EventTopic, config.Config.EventAvsc, 1, config.Config.PublisherMaxOutstanding, config.Config.PublisherNumGoroutines)
	if err != nil {
		return err
	}
	g.Run(event, numPublishers, timeout, times, sleep)
	running = g
	return nil
}

// Stop stops the event generating
func Stop() {
	mux.Lock()
	defer mux.Unlock()

	if running == nil {
		log.Printf("there is no running generator")
		return
	}
	running.Stop()
	running = nil
}
