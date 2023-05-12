// Package config keeps config for used Globally.
package config

import (
	"google/jss/up12/env"
	"log"
)

type config struct {
	Project string
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	Config = config{
		Project: env.GetEnv("GOOGLE_CLOUD_PROJECT", ""),
	}
	log.Printf("using config: %+v", Config)
}
