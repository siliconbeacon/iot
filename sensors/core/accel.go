package core

import (
	"time"
)

// AccelRange represents the dynamic range you wish to use with the
// accelerometer sensor
type AccelRange int

const (
	AccelRange8g AccelRange = 8
	AccelRange4g AccelRange = 4
	AccelRange2g AccelRange = 2
)

// AccelReading represents 3 axis accelerometer readings in milli-gravities
// 1mg = 1/1000 earth standard gravity
type AccelReading struct {
	Timestamp time.Time
	Xmg       float32
	Ymg       float32
	Zmg       float32
}
