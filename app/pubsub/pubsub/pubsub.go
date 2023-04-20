package pubsub

import (
	"context"
	"errors"
	"log"
	"strconv"
	"up12/pubsub/avro"
	"up12/pubsub/config"

	"cloud.google.com/go/pubsub"
	"github.com/linkedin/goavro/v2"
)

// Client is thread safe and can be shared by multiple goroutines
// Close needs not be called at exit, because the Client is available for the lifetime of the program
var Client *pubsub.Client

func init() {
	Client = newClient(context.Background())
}

func newClient(ctx context.Context) *pubsub.Client {
	client, err := pubsub.NewClient(ctx, config.Config.Project)
	if err != nil {
		panic(err)
	}
	return client
}

type topic struct {
	ID         string
	topic      *pubsub.Topic
	codec      *goavro.Codec
	msgChan    chan map[string]interface{}
	publishers []*publisher
}

func NewTopic(topicID string, schemaPath string, msgChan chan map[string]interface{}, batchSize int) *topic {
	t := Client.Topic(topicID)

	t.PublishSettings.CountThreshold = batchSize
	t.PublishSettings.FlowControlSettings.LimitExceededBehavior = pubsub.FlowControlBlock
	// t.PublishSettings.NumGoroutines // TBD, default is 25 * GOMAXPROCS

	return &topic{
		ID:      topicID,
		topic:   t,
		codec:   avro.NewCodedecFromFile(schemaPath),
		msgChan: msgChan,
	}
}

// AddPublishers adds or removes given number of publishers, publishers will start to publish when added.
func (topic *topic) AddPublishers(ctx context.Context, number int) {
	if number < 0 {
		// remove
		newLen := len(topic.publishers) + number
		if newLen < 0 {
			newLen = 0
		}
		stopPubs := topic.publishers[newLen:]
		log.Printf("Stopping %v publishers", len(stopPubs))
		for _, p := range stopPubs {
			p.Stop()
		}
		topic.publishers = topic.publishers[:newLen]
	} else {
		// add
		log.Printf("Starting %v publishers", number)
		for i := 0; i < number; i++ {
			name := topic.ID + "-publisher-" + strconv.Itoa(len(topic.publishers))
			pubCtx, pub := newPublisher(ctx, name, topic)
			pub.run(pubCtx)
			topic.publishers = append(topic.publishers, pub)
		}
	}

}

func (t *topic) Stop() {
	t.topic.Stop()
}

type publisher struct {
	name   string
	topic  *topic
	cancel context.CancelFunc
}

func newPublisher(ctx context.Context, name string, topic *topic) (context.Context, *publisher) {
	pubCtx, cancel := context.WithCancel(ctx)
	return pubCtx, &publisher{
		topic:  topic,
		name:   topic.ID + "-publisher-" + strconv.Itoa(len(topic.publishers)),
		cancel: cancel,
	}
}

func (pub *publisher) run(ctx context.Context) {
	go func() {
		log.Printf("%v: Started", pub.name)
		for {
			select {
			case msg := <-pub.topic.msgChan:
				if json, err := avro.EncodeToJSON(pub.topic.codec, msg); err != nil {
					log.Printf("%v: Ignore invalid message=%v", pub.name, msg)
				} else {
					if err := publish(ctx, pub.name, pub.topic.topic, json); err != nil {
						if errors.Is(err, context.Canceled) {
							log.Printf("%v: Context Canceled, terminated.", pub.name)
							return
						}
					}
				}
			case <-ctx.Done():
				log.Printf("%v: Context done, stopped", pub.name)
				return
			}
		}
	}()
}

func publish(ctx context.Context, name string, topic *pubsub.Topic, data []byte) error {
	msg := &pubsub.Message{
		Data: data,
	}
	result := topic.Publish(ctx, msg)
	if id, err := result.Get(ctx); err != nil {
		log.Printf("%v: Fail to publish to topic=%v message=%v err=%v", name, topic, msg, err)
		return err
	} else {
		log.Printf("%v: Published message ID=%v", name, id)
		return nil
	}
}

func (pub *publisher) Stop() {
	pub.cancel()
}

type Subscription struct {
	ID           string
	subscription *pubsub.Subscription
	codec        *goavro.Codec
}

func NewSubscription(ID string, schemaPath string) *Subscription {
	sub := Client.Subscription(ID)
	// sub.ReceiveSettings.NumGoroutines = 10	// TBD, default is 10

	return &Subscription{
		ID:           ID,
		subscription: sub,
		codec:        avro.NewCodedecFromFile(schemaPath),
	}
}

type MessageHandler func(context.Context, *Message)

type Message struct {
	*pubsub.Message
	Data map[string]interface{}
}

func (sub *Subscription) Receive(ctx context.Context, handler MessageHandler) error {
	return sub.subscription.Receive(ctx, func(ctx context.Context, pubsubMessage *pubsub.Message) {
		log.Printf("Got pubsub message ID=%v", pubsubMessage.ID)
		data, err := avro.DecodeFromJSON(sub.codec, pubsubMessage.Data)
		if err != nil {
			log.Printf("Failed to check schema, message=%v, ", pubsubMessage.ID)
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
