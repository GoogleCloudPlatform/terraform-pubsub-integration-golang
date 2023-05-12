// Package config keeps config for used Globally.
package config

import (
	"google/jss/up12/avro"
	"google/jss/up12/env"
	"log"
	"time"

	"github.com/linkedin/goavro/v2"
)

type config struct {
	RESTPort                string
	Location                string
	EventTopic              string
	EventAvsc               *goavro.Codec // codec is thread safe
	PublisherNumGoroutines  int
	PublisherMaxOutstanding int
	Threads                 int
	Timeout                 time.Duration
	Sleep                   time.Duration
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	Config = config{
		RESTPort:                env.GetEnv("REST_PORT", "8001"),
		Location:                env.GetEnv("GOOGLE_CLOUD_LOCATION", "west"),
		EventTopic:              env.GetEnv("EVENT_TOPIC", "EventTopic"),
		EventAvsc:               avro.NewCodedecFromFile(env.GetEnv("EVENT_AVSC", "Event.avsc")),
		PublisherNumGoroutines:  env.GetEnvInt("PUBLISHER_THREADS", 0), // use default 25 * GOMAXPROCS
		PublisherMaxOutstanding: env.GetEnvInt("PUBLISHER_FLOW_CONTROL_MAX_OUTSTANDING_MESSAGES", 100),
		Threads:                 env.GetEnvInt("EVENT_GENERATOR_THREADS", 200),
		Timeout:                 time.Duration(env.GetEnvFloat64("EVENT_GENERATOR_RUNTIME", 5) * float64(time.Minute)),
		Sleep:                   time.Duration(env.GetEnvFloat64("EVENT_GENERATOR_SLEEP_TIME", 0.2) * float64(time.Second)),
	}
	log.Printf("using config: %+v", Config)
}
