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

// Package config keeps config for used Globally.
package config

import (
	"google/jss/pubsub-integration/avro"
	"google/jss/pubsub-integration/env"
	"log"
	"os"

	"github.com/linkedin/goavro/v2"
)

type config struct {
	Node                     string
	EventCodec               *goavro.Codec
	EventSubscription        string
	MetricsTopic             string
	MetricsCodec             *goavro.Codec
	SubscriberNumGoroutines  int
	SubscriberMaxOutstanding int
	PublisherBatchSize       int
	PublisherNumGoroutines   int
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	log.SetOutput(os.Stdout)

	hostName, err := os.Hostname()
	if err != nil {
		log.Fatalf("fail to get hostname, err: %v", err)
	}

	eventCodec, err := avro.NewCodedecFromFile(env.GetEnv("EVENT_AVSC", "Event.avsc"))
	if err != nil {
		log.Fatalf("fail to create event avro codec, err: %v", err)
	}
	metricsCodec, err := avro.NewCodedecFromFile(env.GetEnv("METRICS_AVSC", "MetricsAck.avsc"))
	if err != nil {
		log.Fatalf("fail to create metrics avro codec, err: %v", err)
	}

	Config = config{
		Node:                     hostName,
		EventSubscription:        env.GetEnv("EVENT_SUBSCRIPTION", "EventSubscription"),
		EventCodec:               eventCodec,
		MetricsTopic:             env.GetEnv("METRICS_TOPIC", "MetricsTopic"),
		MetricsCodec:             metricsCodec,
		SubscriberNumGoroutines:  env.GetEnvInt("SUBSCRIBER_THREADS", 0), // use default 10
		SubscriberMaxOutstanding: env.GetEnvInt("SUBSCRIBER_FLOW_CONTROL_MAX_OUTSTANDING_MESSAGES", 100),
		PublisherBatchSize:       env.GetEnvInt("PUBLISHER_BATCH_SIZE", 100),
		PublisherNumGoroutines:   env.GetEnvInt("PUBLISHER_THREADS", 0), // use default 25 * GOMAXPROCS
	}
	log.Printf("using config: %+v", Config)
}
