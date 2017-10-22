package weather

import (
	"time"
)

type Reading struct {
	Station                    string
	Timestamp                  time.Time
	TemperatureDegreesCelsius  float32
	RelativeHumidityPercentage float32
}

type Readings []*Reading
