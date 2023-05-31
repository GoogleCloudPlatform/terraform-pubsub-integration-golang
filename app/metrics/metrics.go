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

// Package metrics generate metrics from event.
package metrics

import (
	"google/jss/pubsub-integration/avro"
	"google/jss/pubsub-integration/metrics/config"
	"time"
)

// Factory is the interface for creating metrics from given event message
type Factory func(map[string]interface{}, time.Time, time.Time, time.Duration) (map[string]interface{}, error)

// New creates metrics from given event message
func New(event map[string]interface{}, publishTime time.Time, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	sessionEnd, err := avro.GetValue(event, "session_end_time", time.Time{})
	if err != nil {
		return nil, err
	}
	sessionStart, err := avro.GetValue(event, "session_start_time", time.Time{})
	if err != nil {
		return nil, err
	}
	duration := float32(sessionEnd.Sub(sessionStart).Hours())

	return map[string]interface{}{
		"session_id":           event["session_id"],
		"station_id":           event["station_id"],
		"location":             event["location"],
		"event_timestamp":      event["session_end_time"],
		"publish_timestamp":    publishTime.Truncate(time.Microsecond).UTC(),
		"processing_time_sec":  float32(processingTime.Seconds()),
		"ack_timestamp":        ackTime.Truncate(time.Microsecond).UTC(),
		"session_duration_hr":  duration,
		"avg_charge_rate_kw":   event["avg_charge_rate_kw"],
		"battery_capacity_kwh": event["battery_capacity_kwh"],
		"battery_level_start":  event["battery_level_start"],
		"event_node":           event["event_node"],
		"metrics_node":         config.Config.Node,
	}, nil
}
