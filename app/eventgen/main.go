// Package main is the entry point of event generator.
package main

import (
	"google/jss/up12/eventgen/api"
	"google/jss/up12/eventgen/config"
	"google/jss/up12/eventgen/generator"
	"log"
)

func main() {
	if err := generator.Start(generator.NewEvent, config.Config.Threads, config.Config.Timeout, 0, config.Config.Sleep); err != nil {
		log.Fatalf("fail to start generator, err=%v", err)
	}
	api.StartRESTServer()
}
