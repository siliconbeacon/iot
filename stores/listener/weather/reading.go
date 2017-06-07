package weather

import (
	"time"
)

type Reading struct {
	Station                    string
	Timestamp                  time.Time
	TemperatureDegreesCelsius  float64
	RelativeHumidityPercentage float64
}

type Readings []*Reading
