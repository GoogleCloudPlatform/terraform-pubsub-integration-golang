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

// Package metrics generate complete metrics from event.
package metrics

import (
	"google/jss/pubsub-integration/avro"
	"google/jss/pubsub-integration/metrics"
	"time"
)

// New generate complete metrics
func New(event map[string]interface{}, publishTime time.Time, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	metrics, err := metrics.New(event, publishTime, ackTime, processingTime)
	if err != nil {
		return nil, err
	}
	var float32Zero float32
	levelStart, err := avro.GetValue(metrics, "battery_level_start", float32Zero)
	if err != nil {
		return nil, err
	}
	chargeRate, err := avro.GetValue(metrics, "avg_charge_rate_kw", float32Zero)
	if err != nil {
		return nil, err
	}
	duration, err := avro.GetValue(metrics, "session_duration_hr", float32Zero)
	if err != nil {
		return nil, err
	}
	capacity, err := avro.GetValue(metrics, "battery_capacity_kwh", float32Zero)
	if err != nil {
		return nil, err
	}
	levelEnd := levelStart + (chargeRate * duration / capacity)
	if levelEnd > 1.0 {
		levelEnd = 1.0
	}
	metrics["battery_level_end"] = floatValue(levelEnd)
	metrics["charged_total_kwh"] = floatValue((levelEnd - levelStart) * capacity)

	return metrics, nil
}

func floatValue(value float32) map[string]interface{} {
	return map[string]interface{}{"float": value}
}
