package config

import (
	"log"
	"reflect"
	"testing"
)

func TestMetricsAppType(t *testing.T) {
	ma := string(MetricsAck)
	mytype := MetricsAppType(ma + "a")
	log.Println(reflect.TypeOf(mytype))
	log.Println("mytype=", mytype)
}
