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

// Package main is the entry point of MetricsNack.
package main

import (
	"context"
	"errors"
	"google/jss/pubsub-integration/metrics"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	if err := metrics.Start(ctx, newNackMetrics); err != nil {
		log.Fatalf("fail to start metircs nack, err: %v", err)
	}
}

func newNackMetrics(event map[string]interface{}, publishTime time.Time, ackTime time.Time, processingTime time.Duration) (map[string]interface{}, error) {
	// simulate a bug 🐞 and nack the message
	return nil, errors.New("simulate a bug 🐞 and nack the message")
}
