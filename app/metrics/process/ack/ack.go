package ack

import (
	"context"
	"log"
	"time"
	"up12/metrics/config"
	"up12/metrics/process"
	"up12/pubsub/pubsub"
)

type EventProcessor struct {
	MsgChnn chan map[string]interface{}
}

func (processor *EventProcessor) AvscFile() string {
	return config.Config.MetricsAckAvsc
}

func (processor *EventProcessor) EventHandler() pubsub.MessageHandler {
	return func(ctx context.Context, message *pubsub.Message) {
		log.Printf("Processing event=%v", message.Data)

		processingTime := process.ProcessingTime()
		time.Sleep(processingTime) // Simulate processing time

		message.Ack()
		ackTime := time.Now()
		metricsAck := toMetricsAck(message, ackTime, processingTime)
		log.Printf("metricsAck=%v", metricsAck)
		processor.MsgChnn <- metricsAck
	}
}

func toMetricsAck(message *pubsub.Message, ackTime time.Time, processingTime time.Duration) map[string]interface{} {
	event := message.Data
	return map[string]interface{}{
		"session_id":           event["session_id"],
		"station_id":           event["station_id"],
		"location":             event["location"],
		"event_timestamp":      event["session_end_time"],
		"publish_timestamp":    message.PublishTime,
		"processing_time_sec":  processingTime.Seconds(),
		"ack_timestamp":        ackTime,
		"session_duration_hr":  event["session_end_time"].(time.Time).Sub(event["session_start_time"].(time.Time)).Hours(),
		"avg_charge_rate_kw":   event["avg_charge_rate_kw"],
		"battery_capacity_kwh": event["battery_capacity_kwh"],
		"battery_level_start":  event["battery_level_start"],
	}

}

type Member interface {
	GetName() string
	GetAge() int
}

type Robot struct {
	name  string
	age   int
	power int
}

func (r *Robot) Work() {
}

func (r *Robot) GetName() string {
	return r.name
}
func (r *Robot) GetAge() int {
	return r.age
}

var m Member = &Robot{"r", 1, 1}
