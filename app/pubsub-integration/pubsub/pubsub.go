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

// Package pubsub provides API to publish and receive message
package pubsub

import (
	"context"
	"fmt"
	"google/jss/pubsub-integration/avro"
	"google/jss/pubsub-integration/pubsub/config"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	vkit "cloud.google.com/go/pubsub/apiv1"
	"github.com/googleapis/gax-go/v2"
	"github.com/linkedin/goavro/v2"
	"google.golang.org/grpc/codes"
)

type service interface {
	NewClient(context.Context, *pubsub.ClientConfig) (Client, error)
}

// Service used to creates client for bucket handling.
var Service service = new(pubsubService)

type pubsubService struct {
}

// NewClientBackoffConfig creates the default backoff config for pubsub client
func NewClientBackoffConfig(initial time.Duration, max time.Duration) *pubsub.ClientConfig {
	retryer := func() gax.Retryer {
		return gax.OnCodes([]codes.Code{
			codes.Aborted,
			codes.Canceled,
			codes.Internal,
			codes.ResourceExhausted,
			codes.Unknown,
			codes.Unavailable,
			codes.DeadlineExceeded,
		}, gax.Backoff{
			Initial: initial,
			Max:     max,
		})
	}

	return &pubsub.ClientConfig{
		PublisherCallOptions: &vkit.PublisherCallOptions{
			Publish: []gax.CallOption{gax.WithRetry(retryer)},
		},
	}
}

// NewClient creates the client for bucket handling. Using the default backoff config if clientCfg is nil
func (*pubsubService) NewClient(ctx context.Context, clientCfg *pubsub.ClientConfig) (Client, error) {
	client, err := pubsub.NewClientWithConfig(ctx, config.Config.Project, clientCfg)
	if err != nil {
		return nil, err
	}
	return &pubsubClient{client: client}, err
}

// Client is the interface of the pubsub client for pubsub handling.
type Client interface {
	NewTopic(string, *goavro.Codec, int, int, int) Topic
	NewSubscription(string, *goavro.Codec, int, int) *Subscription
	Close() error
}

type pubsubClient struct {
	client *pubsub.Client
}

// NewTopic get the topic for publishing message. Using the default value if batchSize, numGoroutines, maxOutstanding <= 0
func (c *pubsubClient) NewTopic(topicID string, codec *goavro.Codec, batchSize int, numGoroutines int, maxOutstanding int) Topic {
	topic := c.client.Topic(topicID)

	if batchSize > 0 {
		topic.PublishSettings.CountThreshold = batchSize
	}
	if numGoroutines > 0 {
		topic.PublishSettings.NumGoroutines = numGoroutines // default is 25 * GOMAXPROCS
	}
	topic.PublishSettings.FlowControlSettings.LimitExceededBehavior = pubsub.FlowControlBlock
	if maxOutstanding > 0 {
		topic.PublishSettings.FlowControlSettings.MaxOutstandingMessages = maxOutstanding
	}
	return &pubsubTopic{
		id:    topicID,
		topic: topic,
		codec: codec,
	}
}

// NewSubscription gets the subscription for receiving message. Using the default value if maxOutstanding, numGoroutines <= 0
func (c *pubsubClient) NewSubscription(ID string, codec *goavro.Codec, numGoroutines int, maxOutstanding int) *Subscription {
	sub := c.client.Subscription(ID)

	if numGoroutines > 0 {
		sub.ReceiveSettings.NumGoroutines = numGoroutines // default is 10
	}
	if maxOutstanding > 0 {
		sub.ReceiveSettings.MaxOutstandingMessages = maxOutstanding
	}
	return &Subscription{
		ID:           ID,
		subscription: sub,
		codec:        codec,
	}
}

// Close close the underlying client.
func (c *pubsubClient) Close() error {
	log.Printf("close client: %v", c.client)
	return c.client.Close()
}

// Topic is used to publish message to topic
type Topic interface {
	Publish(context.Context, map[string]interface{}) (string, error)
	GetID() string
	Stop()
}

type pubsubTopic struct {
	id    string
	topic *pubsub.Topic
	codec *goavro.Codec
}

// Publish publishing message to the topic
func (t *pubsubTopic) Publish(ctx context.Context, data map[string]interface{}) (string, error) {
	json, err := avro.EncodeToJSON(t.codec, data)
	if err != nil {
		return "", fmt.Errorf("ignore invalid message: %v", data)
	}
	msg := &pubsub.Message{
		Data: json,
	}
	result := t.topic.Publish(ctx, msg)
	id, err := result.Get(ctx)
	if err != nil {
		return id, fmt.Errorf("fail to publish message: %v to topic: %v, err: %w", json, t.topic, err)
	}
	return id, nil
}

func (t *pubsubTopic) GetID() string {
	return t.id
}

func (t *pubsubTopic) Stop() {
	log.Printf("stop topic: %v", t.GetID())
	t.topic.Stop()
}

// Subscription is used to receive message
type Subscription struct {
	ID           string
	subscription *pubsub.Subscription
	codec        *goavro.Codec
}

// MessageHandler is function to handle the received message
type MessageHandler func(context.Context, *Message)

// Message contains the message content decoded by avro schema
type Message struct {
	*pubsub.Message
	Data map[string]interface{}
}

// Receive starts to receive messages and handle them by given message handler
func (sub *Subscription) Receive(ctx context.Context, handler MessageHandler) error {
	return sub.subscription.Receive(ctx, func(ctx context.Context, pubsubMessage *pubsub.Message) {
		log.Printf("got pubsub message ID: %v", pubsubMessage.ID)
		data, err := avro.DecodeFromJSON(sub.codec, pubsubMessage.Data)
		if err != nil {
			log.Printf("failed to check schema, message: %v, ", pubsubMessage.ID)
			pubsubMessage.Nack()
			return
		}
		message := &Message{
			Message: pubsubMessage,
			Data:    data,
		}
		handler(ctx, message)
	})
}
