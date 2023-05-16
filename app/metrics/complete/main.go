// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is the entry point of MetricsComplete.
package main

import (
	"context"
	"google/jss/up12/metrics"
	"google/jss/up12/pubsub"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	if err := metrics.Start(ctx, metricComplete); err != nil {
		log.Fatalf("fail to start metircs complete, err: %v", err)
	}
}

func metricComplete(message *pubsub.Message, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	metric, err := metrics.NewMetric(message, ackTime, processingTime)
	if err != nil {
		return nil, err
	}
	var float32Zero float32
	levelStart, err := metrics.GetValue(metric, "battery_level_start", float32Zero)
	if err != nil {
		return nil, err
	}
	chargeRate, err := metrics.GetValue(metric, "avg_charge_rate_kw", float32Zero)
	if err != nil {
		return nil, err
	}
	duration, err := metrics.GetValue(metric, "session_duration_hr", float32Zero)
	if err != nil {
		return nil, err
	}
	capacity, err := metrics.GetValue(metric, "battery_capacity_kwh", float32Zero)
	if err != nil {
		return nil, err
	}
	levelEnd := levelStart + (chargeRate * duration / capacity)
	if levelEnd > 1.0 {
		levelEnd = 1.0
	}
	metric["battery_level_end"] = floatValue(levelEnd)
	metric["charged_total_kwh"] = floatValue((levelEnd - levelStart) * capacity)

	return metric, nil
}

func floatValue(value float32) map[string]interface{} {
	return map[string]interface{}{"float": value}
}
