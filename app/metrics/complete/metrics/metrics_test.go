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

// Package metrics generate complete metrics
package metrics

import (
	"google/jss/pubsub-integration/avro"
	"google/jss/pubsub-integration/eventgen/generator"
	"google/jss/pubsub-integration/metrics"
	"google/jss/pubsub-integration/metrics/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MetricsComplete calculates fields correctly
func TestMetricsComplete(t *testing.T) {
	event := generator.NewEvent()
	publishTime := time.Now()
	processingTime := metrics.ProcessingTime()
	ackTime := publishTime.Add(processingTime)
	m, err := New(event, publishTime, ackTime, processingTime)
	assert.Nil(t, err)

	// valiate event_timestamp
	assert.Equal(t, event["session_end_time"], m["event_timestamp"])
	// valiate publish_timestamp
	assert.Equal(t, publishTime.Truncate(time.Microsecond).UTC(), m["publish_timestamp"])
	//valiate processing_time_sec
	assert.Equal(t, float32(processingTime.Seconds()), m["processing_time_sec"])
	// valiate ack_timestamp
	assert.Equal(t, ackTime.Truncate(time.Microsecond).UTC(), m["ack_timestamp"])
	// valiate session_duration_hr
	sessionEnd, err := avro.GetValue(event, "session_end_time", time.Time{})
	assert.Nil(t, err)
	sessionStart, err := avro.GetValue(event, "session_start_time", time.Time{})
	assert.Nil(t, err)
	duration := float32(sessionEnd.Sub(sessionStart).Hours())
	assert.Equal(t, duration, m["session_duration_hr"])

	// valiate battery_level_end
	var float32Zero float32
	chargeRate, err := avro.GetValue(event, "avg_charge_rate_kw", float32Zero)
	assert.Nil(t, err)
	capacity, err := avro.GetValue(event, "battery_capacity_kwh", float32Zero)
	assert.Nil(t, err)
	levelStart, err := avro.GetValue(event, "battery_level_start", float32Zero)
	assert.Nil(t, err)
	levelEnd := levelStart + (chargeRate * duration / capacity)
	if levelEnd > 1.0 {
		levelEnd = 1.0
	}
	mLevelEnd, err := avro.GetFloatTypeValue(m, "battery_level_end")
	assert.Nil(t, err)
	assert.Equal(t, levelEnd, mLevelEnd)

	// validate charged_total_kwh
	chargedTotal := (levelEnd - levelStart) * capacity
	mChargedTotal, err := avro.GetFloatTypeValue(m, "charged_total_kwh")
	assert.Nil(t, err)
	assert.Equal(t, chargedTotal, mChargedTotal)
}

// Create a metrics and make sure it’s valid JSON matching the Pub/Sub schema.
func TestMetricsAckAvroCodec(t *testing.T) {
	event := generator.NewEvent()
	publishTime := time.Now()
	processingTime := metrics.ProcessingTime()
	ackTime := publishTime.Add(processingTime)
	m, err := New(event, publishTime, ackTime, processingTime)
	assert.Nil(t, err)

	json, err := avro.EncodeToJSON(config.Config.MetricsAvsc, m)
	assert.Nil(t, err)
	native, err := avro.DecodeFromJSON(config.Config.MetricsAvsc, json)
	assert.Nil(t, err)
	assert.Equal(t, m, native)
}
