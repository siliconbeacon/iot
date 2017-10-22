package core

import (
	"time"
)

// MagReading represents 3 axis magnetometer readings in micro-Teslas
type MagReading struct {
	Timestamp time.Time
	XuT       float32
	YuT       float32
	ZuT       float32
}
