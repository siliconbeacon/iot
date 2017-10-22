package core

import (
	"time"
)

// GyroRange represents the dynamic range you wish to use with the
// gryoscope sensor
type GyroRange int

const (
	GyroRange250dps  GyroRange = 250
	GyroRange500dps  GyroRange = 500
	GyroRange1000dps GyroRange = 1000
	GyroRange2000dps GyroRange = 2000
	GyroRange4000dps GyroRange = 4000
)

// GyroReading represents 3 axis gyroscope readings in degrees per second
type GyroReading struct {
	Timestamp time.Time
	Xdps      float32
	Ydps      float32
	Zdps      float32
}
