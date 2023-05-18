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
package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProcessingTime tests the processing time is between 0.1 and 5 seconds and 99.4% of processing time is between 0.1 and 0.3 seconds
func TestProcessTime(t *testing.T) {
	var count int
	for i := 0; i < 10000; i++ {
		sec := ProcessingTime().Seconds()
		assert.True(t, sec >= 0.1 && sec <= 5) // processing time should be between 0.1 and 5 seconds

		if sec >= 0.1 && sec <= 0.3 {
			count++
		}
	}
	assert.True(t, count >= 9940) // 99.4% (0.5% error margin) of processing time should be between 0.1 and 0.3 seconds
}
