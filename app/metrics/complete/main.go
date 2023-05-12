// Package main is the entry point of MetricsComplete.
package main

import (
	"context"
	"google/jss/up12/metrics"
	"google/jss/up12/metrics/config"
	"google/jss/up12/pubsub"
	"time"
)

func main() {
	ctx := context.Background()
	metrics.Start(ctx, config.Config.MetricsCompleteAvsc, metricComplete)
}

func metricComplete(message *pubsub.Message, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	metric, err := metrics.NewMetric(message, ackTime, processingTime)
	if err != nil {
		return nil, err
	}
	levelStart, err := metrics.GetValue(metric, "battery_level_start", float32(0.0))
	if err != nil {
		return nil, err
	}
	chargeRate, err := metrics.GetValue(metric, "avg_charge_rate_kw", float32(0.0))
	if err != nil {
		return nil, err
	}
	duration, err := metrics.GetValue(metric, "session_duration_hr", float32(0.0))
	if err != nil {
		return nil, err
	}
	capacity, err := metrics.GetValue(metric, "battery_capacity_kwh", float32(0.0))
	if err != nil {
		return nil, err
	}
	levelEnd := levelStart + (chargeRate * duration / capacity)
	if levelEnd > 1.0 {
		levelEnd = 1.0
	}
	metric["battery_level_end"] = levelEnd
	metric["charged_total_kwh"] = (levelEnd - levelStart) * capacity

	return metric, nil
}
