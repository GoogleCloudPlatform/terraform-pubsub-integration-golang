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
	"time"

	"github.com/linkedin/goavro/v2"
)

type config struct {
	Node                    string
	RESTPort                string
	Location                string
	EventTopic              string
	EventCodec              *goavro.Codec // codec is thread safe
	PublisherBatchSize      int
	PublisherNumGoroutines  int
	PublisherMaxOutstanding int
	PublisherRetryInit      time.Duration
	PublisherRetryTotal     time.Duration
	Threads                 int
	Timeout                 time.Duration
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

	Config = config{
		Node:                    hostName,
		RESTPort:                env.GetEnv("REST_PORT", "8001"),
		Location:                env.GetEnv("GOOGLE_CLOUD_LOCATION", "west"),
		EventTopic:              env.GetEnv("EVENT_TOPIC", "EventTopic"),
		EventCodec:              eventCodec,
		PublisherBatchSize:      env.GetEnvInt("PUBLISHER_BATCH_SIZE", 100),
		PublisherNumGoroutines:  env.GetEnvInt("PUBLISHER_THREADS", 0), // use default 25 * GOMAXPROCS
		PublisherMaxOutstanding: env.GetEnvInt("PUBLISHER_FLOW_CONTROL_MAX_OUTSTANDING_MESSAGES", 100),
		PublisherRetryInit:      time.Duration(env.GetEnvFloat64("PUBLISHER_RETRY_INITIAL_TIMEOUT", 5) * float64(time.Second)),
		PublisherRetryTotal:     time.Duration(env.GetEnvFloat64("PUBLISHER_RETRY_TOTAL_TIMEOUT", 600) * float64(time.Second)),
		Threads:                 env.GetEnvInt("EVENT_GENERATOR_THREADS", 200),
		Timeout:                 time.Duration(env.GetEnvFloat64("EVENT_GENERATOR_RUNTIME", 5) * float64(time.Minute)),
	}
	log.Printf("using config: %+v", Config)
}
