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

// Package metrics generate ack metrics from event.
package metrics

import (
	"google/jss/pubsub-integration/metrics"
	"time"
)

// New creates ack metrics from given event message, It is just the same with the metrics.New().
// Redeclared here for clarity and test purpose.
func New(event map[string]interface{}, publishTime time.Time, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	return metrics.New(event, publishTime, ackTime, processingTime)
}
