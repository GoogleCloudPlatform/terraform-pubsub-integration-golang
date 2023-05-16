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

// Package main is the entry point of event generator.
package main

import (
	"google/jss/up12/eventgen/api"
	"google/jss/up12/eventgen/config"
	"google/jss/up12/eventgen/generator"
	"log"
)

func main() {
	if err := generator.Start(generator.NewEvent, config.Config.Threads, config.Config.Timeout, 0, config.Config.Sleep); err != nil {
		log.Fatalf("fail to start generator, err: %v", err)
	}
	api.StartRESTServer()
}
