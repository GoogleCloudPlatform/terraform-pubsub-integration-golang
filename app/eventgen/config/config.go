// Package config keeps config for used Globally.
package config

import (
	"log"
	"os"
	"time"
)

type config struct {
	Location     string
	EventAvsc    string
	EventTopicID string
	Publishers   int
	Timeout      time.Duration
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	Config = config{
		Location:     os.Getenv("GOOGLE_CLOUD_LOCATION"),
		EventAvsc:    os.Getenv("EVENT_AVSC"),
		EventTopicID: os.Getenv("EVENT_TOPIC_ID"),
		Publishers:   2,               //TBD
		Timeout:      2 * time.Second, //TBD
	}
	log.Println("using config:", Config)
}
