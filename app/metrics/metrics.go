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

// Package metrics receives event and generate metrics
package metrics

import (
	"context"
	"google/jss/pubsub-integration/metrics/config"
	"google/jss/pubsub-integration/pubsub"
	"log"
	"math/rand"
	"time"
)

// Start starts to receive event and generate metrics
func Start(ctx context.Context, metricsFactory Factory) error {
	client, err := pubsub.Service.NewClient(ctx, nil)
	if err != nil {
		return err
	}
	defer client.Close() // nolint: errcheck

	sub := client.NewSubscription(config.Config.EventSubscription, config.Config.EventCodec, config.Config.SubscriberNumGoroutines, config.Config.SubscriberMaxOutstanding)

	metricsTopic := client.NewTopic(config.Config.MetricsTopic, config.Config.MetricsCodec, config.Config.PublisherBatchSize, config.Config.PublisherNumGoroutines, 0)
	defer metricsTopic.Stop()

	handler := eventHandler(metricsTopic, metricsFactory)
	if err := sub.Receive(ctx, handler); err != nil {
		return err
	}
	log.Println("end of event receiving")
	return nil
}

func eventHandler(metricsTopic pubsub.Topic, factory Factory) pubsub.MessageHandler {
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

// ProcessingTime returns a normal distributed random processing time between 0.1 and 0.5 seconds and 99.9% of the time between 0.1 and 0.3 seconds
func ProcessingTime() time.Duration {
	for {
		seconds := random.NormFloat64()*processTimeStdDev + processTimeMean
		if seconds >= proesssTimeMin && seconds <= 5.0 {
			return time.Duration(seconds * float64(time.Second))
		}
	}
}

// Factory is the interface for creating metrics from given event message
type Factory func(map[string]interface{}, time.Time, time.Time, time.Duration) (map[string]interface{}, error)
