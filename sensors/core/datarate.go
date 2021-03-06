// core contains data structures and configuration elements
// that are agnostic to any specific sensor
package core

import (
	"time"
)

// DataRate allows selection from 0.78125Hz to 800Hz
type DataRate struct {
	ID              int
	Name            string
	Period          time.Duration
	OneSecondBuffer int
}

var DataRate800Hz = DataRate{ID: 0, Name: "800Hz", Period: 1250 * time.Microsecond, OneSecondBuffer: 800}
var DataRate400Hz = DataRate{ID: 1, Name: "400Hz", Period: 2500 * time.Microsecond, OneSecondBuffer: 400}
var DataRate200Hz = DataRate{ID: 2, Name: "200Hz", Period: 5 * time.Millisecond, OneSecondBuffer: 200}
var DataRate100Hz = DataRate{ID: 3, Name: "100Hz", Period: 10 * time.Millisecond, OneSecondBuffer: 100}
var DataRate50Hz = DataRate{ID: 4, Name: "50Hz", Period: 20 * time.Millisecond, OneSecondBuffer: 50}
var DataRate25Hz = DataRate{ID: 5, Name: "25Hz", Period: 40 * time.Millisecond, OneSecondBuffer: 25}
var DataRate12_5Hz = DataRate{ID: 6, Name: "12.5Hz", Period: 80 * time.Millisecond, OneSecondBuffer: 13}
var DataRate6_25Hz = DataRate{ID: 7, Name: "6.25Hz", Period: 160 * time.Millisecond, OneSecondBuffer: 7}
var DataRate3_125Hz = DataRate{ID: 8, Name: "3.125Hz", Period: 320 * time.Millisecond, OneSecondBuffer: 4}
var DataRate1_5625Hz = DataRate{ID: 9, Name: "1.5625Hz", Period: 640 * time.Millisecond, OneSecondBuffer: 2}
var DataRate0_78125Hz = DataRate{ID: 10, Name: "0.78125Hz", Period: 1280 * time.Millisecond, OneSecondBuffer: 1}
