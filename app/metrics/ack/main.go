// Package main is the entry point of MetricsAck.
package main

import (
	"context"
	"google/jss/up12/metrics"
	"google/jss/up12/metrics/config"
)

func main() {
	ctx := context.Background()
	metrics.Start(ctx, config.Config.MetricsAckAvsc, metrics.NewMetric)
}
