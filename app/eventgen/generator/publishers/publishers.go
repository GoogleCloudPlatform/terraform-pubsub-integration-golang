// Package publishers create group of publishers to publish message
package publishers

import (
	"context"
	"google/jss/up12/pubsub"
	"log"
	"strconv"
	"sync"
	"time"
)

type NewMessage func() map[string]interface{}

type Publishers struct {
	pubsub.Topic
	newMessage NewMessage
	publishers []*publisher
	sleep      time.Duration
	times      int
	timeout    time.Duration
	sync.Locker
	waitFinish *sync.Cond
}

// NewPublishers gets the publishers group for publishing message concurrently
func NewPublishers(topic pubsub.Topic, newMessage NewMessage, timeout time.Duration, times int, sleep time.Duration) *Publishers {
	var mux sync.Mutex
	return &Publishers{
		Topic:      topic,
		newMessage: newMessage,
		sleep:      sleep,
		times:      times,
		timeout:    timeout,
		Locker:     &mux,
		waitFinish: sync.NewCond(&mux),
	}
}

// Add adds or removes given number of publishers, publishers will start to publish when added.
func (pbrs *Publishers) Add(ctx context.Context, number int) {
	pbrs.Lock()
	defer pbrs.Unlock()

	if number < 0 {
		// remove
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
		// add
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

func (pbrs *Publishers) Stop() {
	pbrs.Lock()
	defer pbrs.Unlock()

	for _, pbr := range pbrs.publishers {
		pbr.Stop()
	}
}

func (pbrs *Publishers) WaitFinish() {
	pbrs.Lock()
	defer pbrs.Unlock()
	for len(pbrs.publishers) > 0 {
		pbrs.waitFinish.Wait() // waiting until no running publishers
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

// run publisher unitl ctx done
func (pbr *publisher) run(ctx context.Context) {
	var pbrCtx context.Context
	if pbr.timeout > 0 {
		pbrCtx, pbr.cancel = context.WithTimeout(ctx, pbr.timeout)
	} else {
		pbrCtx, pbr.cancel = context.WithCancel(ctx)
	}

	go func() {
		defer pbr.finish()
		log.Printf("%v: started", pbr.name)
		var count int
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
			if pbr.sleep > 0 {
				time.Sleep(pbr.sleep)
			}
			if pbr.times > 0 {
				count++
				if count >= pbr.times {
					break
				}
			}
		}
	}()
}

// finish removes itself from publishers
func (pbr *publisher) finish() {
	pbr.Lock()
	defer pbr.Unlock()

	pbr.publishers = remove(pbr.publishers, pbr)
	if len(pbr.publishers) == 0 {
		pbr.waitFinish.Broadcast() // all publisers are stopped, release all waiting routings
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

// Stop stops the publisher gracefully
func (pbr *publisher) Stop() {
	pbr.cancel()
}
