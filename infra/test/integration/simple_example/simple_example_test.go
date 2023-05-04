// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simple_example

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestSimpleExample(t *testing.T) {
	example := tft.NewTFBlueprintTest(t)

	example.DefineVerify(func(assert *assert.Assertions) {
		projectID := example.GetTFSetupStringOutput("project_id")
		gcloudArgs := gcloud.WithCommonArgs([]string{"--project", projectID})

		// Check subscriptions state
		subscriptions := gcloud.Run(t, ("pubsub subscriptions list --format=json"), gcloudArgs).Array()
		for _, subscription := range subscriptions {
			state := subscription.Get("state").String()
			assert.Equal("ACTIVE", state, "expected subscriptions to be active")
		}

		// Check if the ErrorTopic exists
		errorTopicName := example.GetStringOutput("errors_topic_name")
		errorTopic := gcloud.Run(t, fmt.Sprintf("pubsub topics describe %s --format=json", errorTopicName), gcloudArgs)
		assert.NotEmpty(errorTopic)

		// Check if the MetricsTopic exists
		metricsTopicName := example.GetStringOutput("metrics_topic_name")
		metricsTopic := gcloud.Run(t, fmt.Sprintf("pubsub topics describe %s --format=json", metricsTopicName), gcloudArgs)
		assert.NotEmpty(metricsTopic)

		// Check if the EventTopic exists
		eventTopicName := example.GetStringOutput("event_topic_name")
		eventTopic := gcloud.Run(t, fmt.Sprintf("pubsub topics describe %s --format=json", eventTopicName), gcloudArgs)
		assert.NotEmpty(eventTopic)

		// Check container clusters state
		clusters := gcloud.Run(t, ("container clusters list --format=json"), gcloudArgs).Array()
		for _, cluster := range clusters {
			state := cluster.Get("state").String()
			assert.Equal("RUNNING", state, "expected container clusters to be active")
		}

	})

	example.Test()
}
