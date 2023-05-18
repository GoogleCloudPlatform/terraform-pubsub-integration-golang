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
	"google/jss/pubsub-integration/eventgen/config"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var avgChargeRateKWValues = [5]float32{20, 72, 100, 120, 250}
var batteryCapacityKWH = [10]float32{40, 50, 58, 62, 75, 77, 82, 100, 129, 131}

func newSessionStartTime(now time.Time) time.Time {
	return now.Add(-1 * time.Duration(random.Intn(86)+5) * time.Minute) // 5 ~ 90 minutes ago
}

func newAvgChargeRateKW() float32 {
	avg := avgChargeRateKWValues[random.Intn(len(avgChargeRateKWValues))]
	avg += (random.Float32() * 2) - 1 // +-1
	return avg
}

func newBatteryCapacityKWH() float32 {
	return batteryCapacityKWH[random.Intn(len(batteryCapacityKWH))]
}

func newBatteryLevelStart() float32 {
	return (float32(random.Intn(76)) + 5) / 100 // 0.05 ~ 0.8
}

// NewEvent creates a new event
func NewEvent() map[string]interface{} {
	now := time.Now().Truncate(time.Microsecond).UTC()
	return map[string]interface{}{
		"session_id":           uuid.New().String(),
		"station_id":           int32(random.Intn(101)),
		"location":             config.Config.Location,
		"session_start_time":   newSessionStartTime(now),
		"session_end_time":     now,
		"avg_charge_rate_kw":   newAvgChargeRateKW(),
		"battery_capacity_kwh": newBatteryCapacityKWH(),
		"battery_level_start":  newBatteryLevelStart(),
	}
}
