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

// Package processor handles the received event and generates metrics
package processor

import (
	"context"
	"google/jss/pubsub-integration/metrics"
	"google/jss/pubsub-integration/metrics/config"
	"google/jss/pubsub-integration/pubsub"
	"log"
	"math/rand"
	"time"
)

// Start starts to receive event and generate metrics
func Start(ctx context.Context, factory metrics.Factory) error {
	client, err := pubsub.Service.NewClient(ctx, nil)
	if err != nil {
		return err
	}
	defer client.Close() // nolint: errcheck

	// The subscription to receive event
	sub := client.NewSubscription(config.Config.EventSubscription, config.Config.EventCodec, config.Config.SubscriberNumGoroutines, config.Config.SubscriberMaxOutstanding)

	// The topic to publish the metrics converted from received event
	metricsTopic := client.NewTopic(config.Config.MetricsTopic, config.Config.MetricsCodec, config.Config.PublisherBatchSize, config.Config.PublisherNumGoroutines, 0)
	defer metricsTopic.Stop()

	// The handler to handles the received event, generate and publish metrics to the metrics topic
	handler := eventHandler(metricsTopic, factory)

	// Start to handle received event using given handler.
	// It does not return until the context is done
	for {
		if err := sub.Receive(ctx, handler); err != nil {
			log.Printf("sub.Receive: %v", err)
		}
		select {
		case <-ctx.Done():
			log.Printf("context done, subscriber stopped")
			return nil
		default:
			waitTime := 30 * time.Second
			log.Printf("waiting %v for retry", waitTime)
			time.Sleep(waitTime)
		}
	}
}

// eventHandler creates the event message handler for subscriber to handle the received event
// The handler receives event message and generates metrics using the given metrics factory
// It acks the message and publishes the metrics to the metrics topic if it generates metrics successfully or nacks if it does not
func eventHandler(metricsTopic pubsub.Topic, factory metrics.Factory) pubsub.MessageHandler {
	// factory: the metrics factory to generate metrics from the received event

	return func(ctx context.Context, message *pubsub.Message) {
		log.Printf("processing event ID: %v, data: %v", message.ID, message.Data)

		processingTime := ProcessingTime()
		time.Sleep(processingTime) // Simulate processing time

		ackTime := time.Now()
		metrics, err := factory(message.Data, message.PublishTime, ackTime, processingTime)
		if err != nil {
			log.Printf("nack the event ID: %v, error: %v", message.ID, err)
			message.Nack()
			return
		}
		log.Printf("event ID: %v converted to metrics: %v", message.ID, metrics)
		id, err := metricsTopic.Publish(ctx, metrics)
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

const proesssTimeMin = 0.1
const processTimeMax = 0.3
const processTimeMean = (proesssTimeMin + processTimeMax) / 2
const processTimeStdDev = (processTimeMax - processTimeMean) / 3.29 // 3.29 is the z-score for 99.9% confidence interval

// ProcessingTime returns a normal distributed random processing time to simulate the time used to process an event
// It is between 0.1 and 0.5 seconds and 99.9% of the time between 0.1 and 0.3 seconds
func ProcessingTime() time.Duration {
	for {
		seconds := random.NormFloat64()*processTimeStdDev + processTimeMean
		if seconds >= proesssTimeMin && seconds <= 5.0 {
			return time.Duration(seconds * float64(time.Second))
		}
	}
}
