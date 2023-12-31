// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package publishers create group of publishers to publish message
package publishers

import (
	"context"
	"google/jss/pubsub-integration/pubsub"
	"log"
	"strconv"
	"sync"
	"time"
)

// NewMessage is the function to generate new message
type NewMessage func() map[string]interface{}

// Publishers is the group of publishers
type Publishers struct {
	pubsub.Topic
	newMessage NewMessage
	publishers []*publisher
	timeout    time.Duration
	sync.Locker
	waitFinish *sync.Cond
}

// NewPublishers creates the publishers group for publishing message concurrently.
// The publishers that have been added will publish messages generated by newMessage function continuously until timeout.
func NewPublishers(topic pubsub.Topic, newMessage NewMessage, timeout time.Duration) *Publishers {
	var mux sync.Mutex
	return &Publishers{
		Topic:      topic,
		newMessage: newMessage,
		timeout:    timeout,
		Locker:     &mux,
		waitFinish: sync.NewCond(&mux),
	}
}

// Adds or removes the publishers based on the given number.
// The publishers will start publishing messages when it is added.
func (pbrs *Publishers) Add(ctx context.Context, number int) {
	pbrs.Lock()
	defer pbrs.Unlock()

	if number < 0 {
		// Remove publishers
		newLen := len(pbrs.publishers) + number
		if newLen < 0 {
			newLen = 0
		}
		stopPubs := pbrs.publishers[newLen:]
		log.Printf("stopping %v publishers", len(stopPubs))
		for _, p := range stopPubs {
			p.Stop()
		}
		pbrs.publishers = pbrs.publishers[:newLen]
	} else {
		// Add publishers
		log.Printf("starting %v publishers", number)
		for i := 0; i < number; i++ {
			name := pbrs.Topic.GetID() + "-publisher-" + strconv.Itoa(len(pbrs.publishers))
			pbrs.addOne(ctx, name)
		}
	}
}

func (pbrs *Publishers) addOne(ctx context.Context, name string) {
	pbrs.publishers = append(pbrs.publishers, runPublisher(ctx, name, pbrs))
}

func (pbrs *Publishers) remove(pbr *publisher) {
	pbrs.Lock()
	defer pbrs.Unlock()

	pbrs.publishers = remove(pbrs.publishers, pbr)
	if len(pbrs.publishers) == 0 {
		pbrs.waitFinish.Broadcast() // All publisers are stopped, release all waiting routings
	}
}

func remove(slice []*publisher, element *publisher) []*publisher {
	for i, e := range slice {
		if e == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Stops all publishers
func (pbrs *Publishers) Stop() {
	pbrs.Lock()
	defer pbrs.Unlock()

	for _, pbr := range pbrs.publishers {
		pbr.Stop()
	}
}

// WaitFinish waits until all publishers are stopped
func (pbrs *Publishers) WaitFinish() {
	pbrs.Lock()
	defer pbrs.Unlock()
	for len(pbrs.publishers) > 0 {
		pbrs.waitFinish.Wait() // Waiting until no running publishers
	}
}

type publisher struct {
	*Publishers
	name   string
	cancel context.CancelFunc
}

func runPublisher(ctx context.Context, name string, publishers *Publishers) *publisher {
	pbr := &publisher{
		Publishers: publishers,
		name:       name,
	}
	pbr.run(ctx)
	return pbr
}

// Starts to run the publisher unitl ctx done
func (pbr *publisher) run(ctx context.Context) {
	var pbrCtx context.Context
	if pbr.timeout > 0 {
		pbrCtx, pbr.cancel = context.WithTimeout(ctx, pbr.timeout)
	} else {
		pbrCtx, pbr.cancel = context.WithCancel(ctx)
	}

	// Create new thread to publish until pbrCtx done
	go func() {
		defer pbr.finish()
		log.Printf("%v: started", pbr.name)
		for {
			select {
			case <-pbrCtx.Done():
				log.Printf("%v: context done, stopped", pbr.name)
				return
			default:
				msg := pbr.newMessage()
				if id, err := pbr.Publish(ctx, msg); err != nil {
					log.Printf("%v: err: %v", pbr.name, err)
				} else {
					log.Printf("%v: published message ID: %v", pbr.name, id)
				}
			}
		}
	}()
}

// Removes itself from publishers
func (pbr *publisher) finish() {
	pbr.Publishers.remove(pbr)
}

// Stops the publisher gracefully
func (pbr *publisher) Stop() {
	pbr.cancel()
}
