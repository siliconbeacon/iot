package core

import (
	"time"
)

// MagReading represents 3 axis magnetometer readings in micro-Teslas
type MagReading struct {
	Timestamp time.Time
	XuT       float64
	YuT       float64
	ZuT       float64
}
