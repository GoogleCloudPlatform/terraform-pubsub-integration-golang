// Package pubsub provides API to publish and receive message
package pubsub

import (
	"context"
	"fmt"
	"google/jss/up12/avro"
	"google/jss/up12/pubsub/config"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/linkedin/goavro/v2"
)

type service interface {
	NewClient(context.Context) (Client, error)
}

// Service used to creates client for bucket handling.
var Service service = new(pubsubService)

type pubsubService struct {
}

// NewClient creates the client for bucket handling.
func (*pubsubService) NewClient(ctx context.Context) (Client, error) {
	client, err := pubsub.NewClient(ctx, config.Config.Project)
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

// NewTopic get the topic for publishing message
func (c *pubsubClient) NewTopic(topicID string, codec *goavro.Codec, batchSize int, maxOutstanding int, numGoroutines int) Topic {
	topic := c.client.Topic(topicID)

	topic.PublishSettings.CountThreshold = batchSize
	topic.PublishSettings.FlowControlSettings.LimitExceededBehavior = pubsub.FlowControlBlock
	if maxOutstanding > 0 {
		topic.PublishSettings.FlowControlSettings.MaxOutstandingMessages = maxOutstanding
	}
	if numGoroutines > 0 {
		topic.PublishSettings.NumGoroutines = numGoroutines // default is 25 * GOMAXPROCS
	}
	return &pubsubTopic{
		id:    topicID,
		topic: topic,
		codec: codec,
	}
}

// NewSubscription gets the subscription for receiving message
func (c *pubsubClient) NewSubscription(ID string, codec *goavro.Codec, maxOutstanding int, numGoroutines int) *Subscription {
	sub := c.client.Subscription(ID)

	if maxOutstanding > 0 {
		sub.ReceiveSettings.MaxOutstandingMessages = maxOutstanding
	}
	if numGoroutines > 0 {
		sub.ReceiveSettings.NumGoroutines = numGoroutines // default is 10
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
