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
	client, err := pubsub.Service.NewClient(ctx, nil)
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
	err := Start(NewEvent, 2, timeout, 0, time.Second)
	assert.Nil(t, err)
	running.publishers.WaitFinish()
	elapsed := time.Since(now).Seconds()
	log.Printf("elapsed: %v", elapsed)
	assert.Equal(t, int(timeout.Seconds()), int(elapsed))
	time.Sleep(1 * time.Second)
}

func TestGeneratorTimes(t *testing.T) {
	times := 3
	threads := 2
	err := Start(countNewEvent, threads, 0, times, time.Second)
	assert.Nil(t, err)
	running.publishers.WaitFinish()
	assert.Equal(t, int32(times*threads), counter)
	time.Sleep(1 * time.Second)
}

var counter int32

func countNewEvent() map[string]interface{} {
	atomic.AddInt32(&counter, 1)
	return NewEvent()
}
