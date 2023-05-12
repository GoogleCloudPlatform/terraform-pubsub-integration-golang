// Package metrics receives event and generate metrics
package metrics

import (
	"context"
	"fmt"
	"google/jss/up12/metrics/config"
	"google/jss/up12/pubsub"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/linkedin/goavro/v2"
)

// Start starts to receive event and generate metrics
func Start(ctx context.Context, codec *goavro.Codec, metricFactory MetricFactory) error {
	client, err := pubsub.Service.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // nolint: errcheck

	sub := client.NewSubscription(config.Config.EventSubscription, config.Config.EventAvsc, config.Config.SubscriberMaxOutstanding, config.Config.SubscriberNumGoroutines)

	metricsTopic := client.NewTopic(config.Config.MetricsTopic, codec, config.Config.BatchSize, 0, config.Config.PublisherNumGoroutines)
	defer metricsTopic.Stop()

	handler := eventHandler(client, metricsTopic, metricFactory)
	if err := sub.Receive(ctx, handler); err != nil {
		return err
	}
	log.Println("end of event receiving")
	return nil
}

func eventHandler(client pubsub.Client, metricsTopic pubsub.Topic, metricFactory MetricFactory) pubsub.MessageHandler {
	return func(ctx context.Context, message *pubsub.Message) {
		log.Printf("processing event ID: %v, data: %v", message.ID, message.Data)

		processingTime := processingTime()
		time.Sleep(processingTime) // Simulate processing time

		ackTime := time.Now()
		metric, err := metricFactory(message, ackTime, processingTime)
		if err != nil {
			log.Printf("nack the event ID: %v, error: %v", message.ID, err)
			message.Nack()
			return
		}
		log.Printf("event ID: %v converted to metric: %v", message.ID, metric)
		id, err := metricsTopic.Publish(ctx, metric)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("event ID: %v is processed and published to metiric topic as message ID: %v", message.ID, id)
		}
		log.Printf("ack the event ID: %v", message.ID)
		message.Ack()
	}
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func processingTime() time.Duration {
	min := 0.1
	max := 5.0
	mean := (max + min) / 2
	stdDev := (max - mean) / 3
	seconds := random.NormFloat64()*stdDev + mean
	return time.Duration(seconds * float64(time.Second))
}

// MetricFactory is the interface for creating metric from given event message
type MetricFactory func(*pubsub.Message, time.Time, time.Duration) (map[string]interface{}, error)

// NewMetric creates metric from given event message
func NewMetric(message *pubsub.Message, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	event := message.Data
	sessionEnd, err := GetValue(event, "session_end_time", time.Time{})
	if err != nil {
		return nil, err
	}
	sessionStart, err := GetValue(event, "session_start_time", time.Time{})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"session_id":           event["session_id"],
		"station_id":           event["station_id"],
		"location":             event["location"],
		"event_timestamp":      event["session_end_time"],
		"publish_timestamp":    message.PublishTime,
		"processing_time_sec":  processingTime.Seconds(),
		"ack_timestamp":        ackTime,
		"session_duration_hr":  float32(sessionEnd.Sub(sessionStart).Hours()),
		"avg_charge_rate_kw":   event["avg_charge_rate_kw"],
		"battery_capacity_kwh": event["battery_capacity_kwh"],
		"battery_level_start":  event["battery_level_start"],
	}, nil
}

// GetValue gets value from given map using given key and converts to given type
func GetValue[T any](mapdata map[string]interface{}, key string, valueType T) (T, error) {
	val, ok := mapdata[key]
	if !ok {
		return valueType, fmt.Errorf("the key %s does not exist", key)
	}
	return toType(val, valueType)
}

func toType[T any](data interface{}, valueType T) (T, error) {
	value, ok := data.(T)
	if !ok {
		return valueType, fmt.Errorf("the type of %v is %v, but %v is expected", data, reflect.TypeOf(data), reflect.TypeOf(valueType))
	}
	return value, nil
}
