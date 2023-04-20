// Package config keeps config for used Globally.
package config

import (
	"log"
	"os"
)

type config struct {
	Project string
}

// Config is the global configuration parsed from environment variables.
var Config config

func init() {
	Config = config{
		Project: os.Getenv("PROJECT"),
	}
	log.Println("using config:", Config)
}
