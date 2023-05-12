// Package api provides RESTful API for controlling event generator
package api

import (
	"google/jss/up12/eventgen/config"
	"google/jss/up12/eventgen/generator"
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
	Thread        int     `form:"thread"`
	ExecutionTime float64 `form:"executionTime"` // in minutes
	Times         int     `form:"times"`
	Sleep         float64 `form:"sleep"` // in seconds
}

func random(c *gin.Context) {
	log.Printf("start to generate event")
	req := GeneratorReq{
		Thread:        config.Config.Threads,
		ExecutionTime: config.Config.Timeout.Minutes(),
		Sleep:         config.Config.Sleep.Seconds(),
	}
	if err := c.BindQuery(&req); err != nil {
		log.Printf("bad request parameters, err: %v", err)
		response(c, http.StatusBadRequest, nil)
		return
	}
	log.Printf("request parameters: %+v", req)
	timeout := time.Duration(req.ExecutionTime * float64(time.Minute))
	sleep := time.Duration(req.Sleep * float64(time.Second))
	if err := generator.Start(generator.NewEvent, req.Thread, timeout, req.Times, sleep); err != nil {
		responseError(c, http.StatusBadRequest, err)
	}
}

func shutdown(c *gin.Context) {
	generator.Stop()
}

// StartRESTServer starts the REST server
func StartRESTServer() {
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
