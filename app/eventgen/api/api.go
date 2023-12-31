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

// Package api provides RESTful API for controlling event generator
package api

import (
	"google/jss/pubsub-integration/eventgen/config"
	"google/jss/pubsub-integration/eventgen/generator"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func response(c *gin.Context, code int, body interface{}) {
	if body == nil {
		c.String(code, "")
	} else {
		c.JSON(code, body)
	}
}

func responseError(c *gin.Context, statusCode int, err error) {
	log.Printf("encounter server error: %v", err)
	response(c, statusCode, "")
}

// GeneratorReq holds the request parameter for generating event
type GeneratorReq struct {
	Threads int     `form:"threads"`
	Runtime float64 `form:"runtime"` // in minutes
}

func random(c *gin.Context) {
	log.Printf("start to generate event")
	req := GeneratorReq{
		Threads: config.Config.Threads,
		Runtime: config.Config.Timeout.Minutes(),
	}
	if err := c.Bind(&req); err != nil {
		log.Printf("bad request parameters, err: %v", err)
		response(c, http.StatusBadRequest, nil)
		return
	}
	log.Printf("request parameters: %+v", req)
	timeout := time.Duration(req.Runtime * float64(time.Minute))
	if err := generator.Start(generator.NewEvent, req.Threads, timeout); err != nil {
		responseError(c, http.StatusBadRequest, err)
	}
}

func shutdown(c *gin.Context) {
	generator.Stop()
}

// StartRESTServer starts the REST server
func StartRESTServer() {
	log.Printf("start REST server")

	router := gin.Default()
	msgRouter := router.Group("/api/msg")

	msgRouter.POST("/random", random)
	msgRouter.POST("/shutdown", shutdown)

	server := &http.Server{
		Addr:    ":" + config.Config.RESTPort,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
