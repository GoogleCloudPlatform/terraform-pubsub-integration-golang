// Package config keeps config for used Globally.
package config

import (
	"log"
	"os"
	"time"
)

type MetricsAppType string

const (
	MetricsAck      MetricsAppType = "MetricsAck"
	MetricsNack     MetricsAppType = "MetricsNack"
	MetricsComplete MetricsAppType = "MetricsComplete"
)

type config struct {
	MetricsAppType      MetricsAppType
	Location            string
	EventAvsc           string
	EventTopicID        string
	SubscriptionID      string
	MetricsTopicID      string
	MetricsAckAvsc      string
	MetricsCompleteAvsc string
	BatchSize           int
	Timeout             time.Duration
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	Config = config{
		MetricsAppType:      MetricsAppType(os.Getenv("METRICS_APP_TYPE")),
		Location:            os.Getenv("LOCATION"),
		EventAvsc:           os.Getenv("EVENT_AVSC"),
		EventTopicID:        os.Getenv("EVENT_TOPIC_ID"),
		SubscriptionID:      os.Getenv("SUBSCRIPTION_ID"),
		MetricsTopicID:      os.Getenv("METRICS_TOPIC_ID"),
		MetricsAckAvsc:      os.Getenv("METRICS_ACK_AVSC"),
		MetricsCompleteAvsc: os.Getenv("METRICS_COMPLETE_AVSC"),
		BatchSize:           100,
		Timeout:             30 * time.Second,
	}
	log.Println("using config:", Config)
}
