// package fxas21001c is a driver for the Freescale 21002c 3-Axis gyroscope.
package fxas21002c

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
	"github.com/siliconbeacon/iot/sensors/core"
)

const (
	fxas21002cAddrDefault   = 0x21
	fxas21002cAddrAlternate = 0x20

	fxas21002cIdentifier = 0xD7

	fxas21002cRegisterStatus = 0x00
	fxas21002cRegisterOutX1  = 0x01
	fxas21002cRegisterOutX2  = 0x02
	fxas21002cRegisterOutY1  = 0x03
	fxas21002cRegisterOutY2  = 0x04
	fxas21002cRegisterOutZ1  = 0x05
	fxas21002cRegisterOutZ2  = 0x06
	fxas21002cRegisterWhoami = 0x0C
	fxas21002cRegisterCtrl0  = 0x0D
	fxas21002cRegisterCtrl1  = 0x13
	fxas21002cRegisterCtrl2  = 0x14
	fxas21002cRegisterCtrl3  = 0x15

	fxas21002cCtrlReset0 = 0x00
	fxas21002cCtrlReset1 = 0x40
)

const (
	fxas21002cActiveTransitionTime = 60 * time.Millisecond
)

type fxas21002cDataRate struct {
	ctrl1      byte
	bufferSize int
}

var (
	fxas21002cDataRates = map[core.DataRate]*fxas21002cDataRate{
		core.DataRate800Hz:  {ctrl1: 0x02, bufferSize: 800},
		core.DataRate400Hz:  {ctrl1: 0x06, bufferSize: 400},
		core.DataRate200Hz:  {ctrl1: 0x0A, bufferSize: 200},
		core.DataRate100Hz:  {ctrl1: 0x0E, bufferSize: 100},
		core.DataRate50Hz:   {ctrl1: 0x12, bufferSize: 50},
		core.DataRate25Hz:   {ctrl1: 0x16, bufferSize: 25},
		core.DataRate12_5Hz: {ctrl1: 0x1A, bufferSize: 13},
	}
)

type fxas21002cGyroRange struct {
	ctrl0       byte
	ctrl3       byte
	sensitivity float64
}

var (
	fxas21002cGyroRanges = map[core.GyroRange]*fxas21002cGyroRange{
		250:  {ctrl0: 0x03, ctrl3: 0x00, sensitivity: 0.0078125},
		500:  {ctrl0: 0x02, ctrl3: 0x00, sensitivity: 0.015625},
		1000: {ctrl0: 0x01, ctrl3: 0x00, sensitivity: 0.03125},
		2000: {ctrl0: 0x00, ctrl3: 0x00, sensitivity: 0.0625},
		4000: {ctrl0: 0x00, ctrl3: 0x01, sensitivity: 0.125},
	}
)

type Fxas21002c struct {
	Bus embd.I2CBus
	// Addr of the sensor.
	address     byte
	initialized bool
	mu          sync.RWMutex

	rge      core.GyroRange
	rgeInfo  *fxas21002cGyroRange
	rate     core.DataRate
	rateInfo *fxas21002cDataRate

	readings chan *core.GyroReading
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
func (d *Fxas21002c) Readings() <-chan *core.GyroReading {
	return d.readings
}

// Start produces a stream of gyroscope readings in the Readings() channel
func (d *Fxas21002c) Start(rge core.GyroRange, rate core.DataRate) error {

	if err := d.Activate(rge, rate); err != nil {
		return err
	}

	// buffer is based on data rate
	d.readings = make(chan *core.GyroReading, d.rateInfo.bufferSize)

	go func() {
		time.Sleep(fxas21002cActiveTransitionTime)
		ticker := time.NewTicker(d.rate.Period)
		for {
			select {
			case <-ticker.C:
				var (
					reading *core.GyroReading
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

func (d *Fxas21002c) Activate(rge core.GyroRange, rate core.DataRate) error {
	if rate.ID > core.DataRate12_5Hz.ID {
		return errors.New("Lowest supported data rate is 12.5Hz")
	}

	if err := d.setup(); err != nil {
		return err
	}

	// grab configuration from tables
	d.rate = rate
	d.rge = rge
	d.rgeInfo = fxas21002cGyroRanges[rge]
	d.rateInfo = fxas21002cDataRates[rate]

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl0, d.rgeInfo.ctrl0); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl3, d.rgeInfo.ctrl3); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl1, d.rateInfo.ctrl1); err != nil {
		return err
	}
	return nil
}

func (d *Fxas21002c) ReadGyro() (*core.GyroReading, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var result []byte
	var err error
	result, err = d.Bus.ReadBytes(d.address, 7)
	if err != nil {
		return nil, err
	}
	return &core.GyroReading{
		Xdps: float64(int16(binary.BigEndian.Uint16(result[1:3]))) * d.rgeInfo.sensitivity,
		Ydps: float64(int16(binary.BigEndian.Uint16(result[3:5]))) * d.rgeInfo.sensitivity,
		Zdps: float64(int16(binary.BigEndian.Uint16(result[5:7]))) * d.rgeInfo.sensitivity,
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
