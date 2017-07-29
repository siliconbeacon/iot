// package fxos8700cq is a driver for the Freescale 8700 6-Axis
// integrated linear accelerometer and magnetometer
package fxos8700cq

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

const (
	fxos8700cqAddrDefault    = 0x1F
	fxos8700cqAddrAlternate1 = 0x1E
	fxos8700cqAddrAlternate2 = 0x1D
	fxos8700cqAddrAlternate3 = 0x1C

	fxos8700cqIdentifier = 0xC7

	fxos8700cqRegisterAccelStatus = 0x00
	fxos8700cqRegisterAccelOutX1  = 0x01
	fxos8700cqRegisterAccelOutX2  = 0x02
	fxos8700cqRegisterAccelOutY1  = 0x03
	fxos8700cqRegisterAccelOutY2  = 0x04
	fxos8700cqRegisterAccelOutZ1  = 0x05
	fxos8700cqRegisterAccelOutZ2  = 0x06

	fxos8700cqRegisterWhoami = 0x0D
	fxos8700cqRegisterCtrl1  = 0x2A
	fxos8700cqRegisterCtrl2  = 0x2B
	fxos8700cqRegisterCtrl3  = 0x2C
	fxos8700cqRegisterCtrl4  = 0x2D
	fxos8700cqRegisterCtrl5  = 0x2E

	fxos8700cqRegisterMagStatus = 0x32
	fxos8700cqRegisterMagOutX1  = 0x33
	fxos8700cqRegisterMagOutX2  = 0x34
	fxos8700cqRegisterMagOutY1  = 0x35
	fxos8700cqRegisterMagOutY2  = 0x36
	fxos8700cqRegisterMagOutZ1  = 0x37
	fxos8700cqRegisterMagOutZ2  = 0x38

	fxos8700cqRegisterMCtrl1 = 0x5B
	fxos8700cqRegisterMCtrl2 = 0x5C
	fxos8700cqRegisterMCtrl3 = 0x5D

	fxos8700cqCtrlReset0 = 0x00
	fxos8700cqCtrlReset1 = 0x40
)

const (
	fxas21002cActiveTransitionTime = 60 * time.Millisecond
)

type fxas21002cDataRate struct {
	ctrl1        byte
	bufferSize   int
	readInterval time.Duration
}

var (
	fxas21002cDataRates = map[DataRate]*fxas21002cDataRate{
		DataRate800Hz:  {ctrl1: 0x02, bufferSize: 800, readInterval: 1250 * time.Microsecond},
		DataRate400Hz:  {ctrl1: 0x06, bufferSize: 400, readInterval: 2500 * time.Microsecond},
		DataRate200Hz:  {ctrl1: 0x0A, bufferSize: 200, readInterval: 5 * time.Millisecond},
		DataRate100Hz:  {ctrl1: 0x0E, bufferSize: 100, readInterval: 10 * time.Millisecond},
		DataRate50Hz:   {ctrl1: 0x12, bufferSize: 50, readInterval: 20 * time.Millisecond},
		DataRate25Hz:   {ctrl1: 0x16, bufferSize: 25, readInterval: 40 * time.Millisecond},
		DataRate12_5Hz: {ctrl1: 0x1A, bufferSize: 13, readInterval: 80 * time.Millisecond},
	}
)

type GyroRange int

const (
	GyroRange250dps  GyroRange = 250
	GyroRange500dps  GyroRange = 500
	GyroRange1000dps GyroRange = 1000
	GyroRange2000dps GyroRange = 2000
	GyroRange4000dps GyroRange = 4000
)

type fxas21002cGyroRange struct {
	ctrl0       byte
	ctrl3       byte
	sensitivity float64
}

var (
	fxas21002cGyroRanges = map[GyroRange]*fxas21002cGyroRange{
		250:  {ctrl0: 0x03, ctrl3: 0x00, sensitivity: 0.0078125},
		500:  {ctrl0: 0x02, ctrl3: 0x00, sensitivity: 0.015625},
		1000: {ctrl0: 0x01, ctrl3: 0x00, sensitivity: 0.03125},
		2000: {ctrl0: 0x00, ctrl3: 0x00, sensitivity: 0.0625},
		4000: {ctrl0: 0x00, ctrl3: 0x01, sensitivity: 0.125},
	}
)

// Acceleromter Reading on 3 axes, in std Earth gravities
type AccelerometerReading struct {
	Timestamp time.Time
	Xg        float64
	Yg        float64
	Zg        float64
}

// Magnetometer Reading on 3 axes, in std

type Fxas21002c struct {
	Bus embd.I2CBus
	// Addr of the sensor.
	address     byte
	rangeInfo   *fxas21002cGyroRange
	rateInfo    *fxas21002cDataRate
	initialized bool
	mu          sync.RWMutex

	readings chan *GyroReading
	closing  chan chan struct{}
}

// New creates a new fxas21002c sensor.
func New(bus embd.I2CBus) *Fxas21002c {
	return &Fxas21002c{
		Bus:     bus,
		address: fxas21002cAddrDefault,
	}
}

// IsPresent returns true if it looks like we were able to see the sensor
func (d *Fxas21002c) IsPresent() bool {
	if err := d.setup(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Readings is a channel that will contain sensor readings after calling Start()
func (d *Fxas21002c) Readings() <-chan *GyroReading {
	return d.readings
}

// Start produces a stream of gyroscope readings in the Readings() channel
func (d *Fxas21002c) Start(rge GyroRange, rate DataRate) error {

	go func() {
		// we are in standby mode. Configure Sensor
		if err := d.Activate(rge, rate); err != nil {
			glog.Errorf("fxas21002c: %v", err)
			return
		}

		// buffer is based on data rate
		d.readings = make(chan *GyroReading, d.rateInfo.bufferSize)
		time.Sleep(fxas21002cActiveTransitionTime)
		ticker := time.NewTicker(d.rateInfo.readInterval)
		for {
			select {
			case <-ticker.C:
				var (
					reading *GyroReading
					err     error
				)
				reading, err = d.ReadGyro()
				if err != nil {
					glog.Errorf("fxas21002c: %v", err)
					continue
				}
				reading.Timestamp = time.Now().UTC()

				d.readings <- reading
			case waitc := <-d.closing:
				waitc <- struct{}{}
				ticker.Stop()
				close(d.readings)
				return
			}
		}
	}()
	return nil
}

// Close stops any period reads in progress.  Call this to stop the readings that
// Start() begins
func (d *Fxas21002c) Close() error {
	if err := d.setup(); err != nil {
		return err
	}
	if d.closing != nil {
		waitc := make(chan struct{})
		d.closing <- waitc
		<-waitc
	}

	return nil
}

func (d *Fxas21002c) setup() error {

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.validate(); err != nil {
		return err
	}

	// check whoami register first
	whoami, err := d.Bus.ReadByteFromReg(d.address, fxas21002cRegisterWhoami)
	if err != nil {
		return err
	}
	if whoami != fxas21002cIdentifier {
		return fmt.Errorf("fxas21002c whoami check failed.  Expected %#x, got %#x", fxas21002cIdentifier, whoami)
	}

	// reset sensor, leave in sleep mode
	err = d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl1, fxas21002cCtrlReset0)
	if err != nil {
		return errors.New("Unable to reset fxas21002c")
	}

	return nil
}

func (d *Fxas21002c) Activate(rge GyroRange, rate DataRate) error {
	if err := d.setup(); err != nil {
		return err
	}

	// grab configuration from tables
	d.rangeInfo = fxas21002cGyroRanges[rge]
	d.rateInfo = fxas21002cDataRates[rate]

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl0, d.rangeInfo.ctrl0); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl3, d.rangeInfo.ctrl3); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl1, d.rateInfo.ctrl1); err != nil {
		return err
	}
	return nil
}

func (d *Fxas21002c) ReadGyro() (*GyroReading, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []byte
	var err error
	result, err = d.Bus.ReadBytes(d.address, 7)
	if err != nil {
		return nil, err
	}
	return &GyroReading{
		Xdps: float64(int16(binary.BigEndian.Uint16(result[1:3]))) * d.rangeInfo.sensitivity,
		Ydps: float64(int16(binary.BigEndian.Uint16(result[3:5]))) * d.rangeInfo.sensitivity,
		Zdps: float64(int16(binary.BigEndian.Uint16(result[5:7]))) * d.rangeInfo.sensitivity,
	}, nil
}

func (d *Fxas21002c) validate() error {
	if d.Bus == nil {
		return errors.New("fxas21002c: no i2c bus")
	}
	if d.address == 0x00 {
		return fmt.Errorf("fxas21002c: invalid address %#x", d.address)
	}
	return nil
}
