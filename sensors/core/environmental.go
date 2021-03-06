package core

import (
	"time"
)

// TemperatureAndHumdityReading is a single reading of both Temperature
// and Humdity values
type TemperatureAndHumidityReading struct {
	Timestamp                  time.Time
	TemperatureDegreesCelsius  float32
	RelativeHumidityPercentage float32
}
