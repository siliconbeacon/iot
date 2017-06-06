package listener

import (
	"time"
)

type WeatherReading struct {
	Station                    string
	Timestamp                  time.Time
	TemperatureDegreesCelsius  float64
	RelativeHumidityPercentage float64
}

type WeatherReadings []*WeatherReading
