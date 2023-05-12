// Package config keeps config for used Globally.
package config

import (
	"google/jss/up12/avro"
	"google/jss/up12/env"
	"log"

	"github.com/linkedin/goavro/v2"
)

type config struct {
	EventAvsc                *goavro.Codec
	EventSubscription        string
	MetricsTopic             string
	MetricsAckAvsc           *goavro.Codec
	MetricsCompleteAvsc      *goavro.Codec
	SubscriberMaxOutstanding int
	SubscriberNumGoroutines  int
	PublisherNumGoroutines   int
	BatchSize                int
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	Config = config{
		EventSubscription:        env.GetEnv("EVENT_SUBSCRIPTION", "EventSubscription"),
		EventAvsc:                avro.NewCodedecFromFile(env.GetEnv("EVENT_AVSC", "Event.avsc")),
		MetricsTopic:             env.GetEnv("METRICS_TOPIC", "MetricsTopic"),
		MetricsAckAvsc:           avro.NewCodedecFromFile(env.GetEnv("METRICS_ACK_AVSC", "MetricsAck.avsc")),
		MetricsCompleteAvsc:      avro.NewCodedecFromFile(env.GetEnv("METRICS_COMPLETE_AVSC", "MetricsComplete.avsc")),
		SubscriberMaxOutstanding: env.GetEnvInt("SUBSCRIBER_FLOW_CONTROL_MAX_OUTSTANDING_MESSAGES", 100),
		SubscriberNumGoroutines:  env.GetEnvInt("SUBSCRIBER_THREADS", 0), // use default 10
		PublisherNumGoroutines:   env.GetEnvInt("PUBLISHER_THREADS", 0),  // use default 25 * GOMAXPROCS
		BatchSize:                env.GetEnvInt("PUBLISHER_BATCH_SIZE", 100),
	}
	log.Printf("using config: %+v", Config)
}
