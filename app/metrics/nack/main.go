// Package main is the entry point of MetricsNack.
package main

import (
	"context"
	"errors"
	"google/jss/up12/metrics"
	"google/jss/up12/metrics/config"
	"google/jss/up12/pubsub"
	"time"
)

func main() {
	ctx := context.Background()
	metrics.Start(ctx, config.Config.MetricsAckAvsc, metricNack)
}

func metricNack(message *pubsub.Message, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	// simulate a bug 🐞 and nack the message
	return nil, errors.New("simulate a bug 🐞 and nack the message")
}
