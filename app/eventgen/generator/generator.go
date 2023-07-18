// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package generator creates event message and publish to event topic
package generator

import (
	"context"
	"errors"
	"google/jss/pubsub-integration/eventgen/config"
	"google/jss/pubsub-integration/eventgen/generator/publishers"
	"google/jss/pubsub-integration/pubsub"
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

// newGenerator initializes the Cloud Pub/Sub client and topic of the generator
func newGenerator(topicID string, codec *goavro.Codec, batchSize int, numGoroutines int, maxOutstanding int) (*generator, error) {
	var g generator

	backoff := pubsub.NewClientBackoffConfig(config.Config.PublisherRetryInit, config.Config.PublisherRetryTotal)
	client, err := pubsub.Service.NewClient(context.Background(), backoff)
	if err != nil {
		log.Printf("fail to connect to Cloud Pub/Sub, err: %v", err)
		return nil, err
	}
	g.client = client

	g.topic = client.NewTopic(topicID, codec, batchSize, numGoroutines, maxOutstanding)
	return &g, nil
}

// Run creates the publisher group and starts to publish events
func (g *generator) Run(event publishers.NewMessage, numPublishers int, timeout time.Duration) {
	log.Printf("run event generator with numPublishers: %v, timeout: %v", numPublishers, timeout)
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel

	pbrs := publishers.NewPublishers(g.topic, event, timeout)
	pbrs.Add(ctx, numPublishers)
	g.publishers = pbrs

	// Wait for all publishers to finish and release resources in another thread
	go func() {
		pbrs.WaitFinish()
		g.release()
	}()
}

// Stop stops the generator and release resources
func (g *generator) Stop() {
	if g.publishers != nil {
		g.publishers.Stop() // Resources will be released after all publishers have finished
		g.cancel()          // force stop generator
	} else {
		g.release()
	}
}

// release stops the topic and close the Cloud Pub/Sub client
func (g *generator) release() {
	g.topic.Stop()
	if err := g.client.Close(); err != nil {
		log.Printf("fail to close Cloud Pub/Sub client, err: %v", err)
	}
	mux.Lock()
	defer mux.Unlock()
	if running == g {
		running = nil
	}
}

var mux sync.Mutex     // Protect running singleon
var running *generator // Singleton, only one generator at the same time.

// Start generates event and publish to event Topic
func Start(event publishers.NewMessage, numPublishers int, timeout time.Duration) error {
	mux.Lock()
	defer mux.Unlock()

	if running != nil {
		return errors.New("there is already an running generator")
	}
	g, err := newGenerator(config.Config.EventTopic, config.Config.EventCodec, config.Config.PublisherBatchSize, config.Config.PublisherNumGoroutines, config.Config.PublisherMaxOutstanding)
	if err != nil {
		return err
	}
	g.Run(event, numPublishers, timeout)
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
