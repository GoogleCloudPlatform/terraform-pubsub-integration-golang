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

// Package generator creates event message and publish to event topic
package generator

import (
	"google/jss/pubsub-integration/avro"
	"google/jss/pubsub-integration/eventgen/config"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Generate 100 messages and show they are within the bounds and have varying values.
func TestNewEvent(t *testing.T) {
	size := 100
	sessionIDs := make(map[interface{}]interface{})
	stationIDS := make(map[int32]interface{})
	chargeRateIDs := make(map[int]interface{})
	capacityIDs := make(map[int]interface{})
	for i := 0; i < size; i++ {
		event := NewEvent()
		sessionIDs[event["session_id"]] = nil

		// Validate station_id
		stationID, err := avro.GetValue(event, "station_id", int32(0))
		assert.Nil(t, err)
		assert.True(t, stationID >= 0 && stationID <= 100)
		stationIDS[stationID] = nil

		// Validate session_start_time and session_end_time are within 5-90 minutes
		sessionStart, err := avro.GetValue(event, "session_start_time", time.Time{})
		assert.Nil(t, err)
		sessionEnd, err := avro.GetValue(event, "session_end_time", time.Time{})
		assert.Nil(t, err)
		duration := sessionEnd.Sub(sessionStart)
		assert.True(t, duration >= 5*time.Minute && duration <= 90*time.Minute)

		// Validate avg_charge_rate_kw
		chargeRate, err := avro.GetValue(event, "avg_charge_rate_kw", float32(0))
		assert.Nil(t, err)
		avgChargeRateIdx, found := sort.Find(len(avgChargeRateKWValues),
			func(i int) int {
				return int(chargeRate - avgChargeRateKWValues[i])
			})
		assert.True(t, found)
		chargeRateIDs[avgChargeRateIdx] = nil // added for checking variety of values

		// Validate battery_capacity_kwh
		capacity, err := avro.GetValue(event, "battery_capacity_kwh", float32(0))
		assert.Nil(t, err)
		capacityIdx, found := sort.Find(len(batteryCapacityKWH),
			func(i int) int {
				return int(capacity - batteryCapacityKWH[i])
			})
		assert.True(t, found)
		capacityIDs[capacityIdx] = nil

		// Validate battery_level_start is between 0.05 ~ 0.8
		levelStart, err := avro.GetValue(event, "battery_level_start", float32(0))
		assert.Nil(t, err)
		assert.True(t, levelStart >= 0.05)
		assert.True(t, levelStart <= 0.8)
	}
	assert.Equal(t, size, len(sessionIDs))
	assert.True(t, len(stationIDS) > 1)
	assert.True(t, len(capacityIDs) == 10)
	assert.True(t, len(chargeRateIDs) == 5)
}

// Create a message and make sure it’s valid JSON matching the Pub/Sub schema.
func TestNewEventAvroCodec(t *testing.T) {
	event := NewEvent()
	json, err := avro.EncodeToJSON(config.Config.EventCodec, event)
	assert.Nil(t, err)
	native, err := avro.DecodeFromJSON(config.Config.EventCodec, json)
	assert.Nil(t, err)
	assert.Equal(t, event, native)
}
