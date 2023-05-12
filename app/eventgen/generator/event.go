package generator

import (
	"google/jss/up12/eventgen/config"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var avgChargeRateKWValues = []int{20, 72, 100, 120, 250}
var batteryCapacityKWH = []int{40, 50, 58, 62, 75, 77, 82, 100, 129, 131}

func newSessionStartTime(now time.Time) time.Time {
	return now.Add(time.Duration(-1*(random.Intn(86)+5)) * time.Minute)
}

func newAvgChargeRateKW() int {
	return avgChargeRateKWValues[random.Intn(len(avgChargeRateKWValues))] + random.Intn(3) - 1
}

func newBatteryCapacityKWH() int {
	return batteryCapacityKWH[random.Intn(len(avgChargeRateKWValues))]
}

func newBatteryLevelStart() float32 {
	return (float32(random.Intn(76)) + 5) / 100 // 0.05 ~ 0.8
}

func NewEvent() map[string]interface{} {
	now := time.Now()
	return map[string]interface{}{
		"session_id":           uuid.New().String(),
		"station_id":           random.Intn(101),
		"location":             config.Config.Location,
		"session_start_time":   newSessionStartTime(now),
		"session_end_time":     now,
		"avg_charge_rate_kw":   newAvgChargeRateKW(),
		"battery_capacity_kwh": newBatteryCapacityKWH(),
		"battery_level_start":  newBatteryLevelStart(),
	}
}
