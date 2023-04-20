// Package main the entry point of event generator
package main

import (
	"context"
	"up12/eventgen/config"
	"up12/eventgen/generator"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), config.Config.Timeout)
	defer cancel()

	generator.Start(ctx)
}
