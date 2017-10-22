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

	fxas21002cCtrlStandby = 0x00
)

const (
	fxas21002cActiveTransitionTime = 60 * time.Millisecond
)

type fxas21002cDataRate struct {
	ctrl1 byte
}

var (
	fxas21002cDataRates = map[core.DataRate]*fxas21002cDataRate{
		core.DataRate800Hz:  {ctrl1: 0x02},
		core.DataRate400Hz:  {ctrl1: 0x06},
		core.DataRate200Hz:  {ctrl1: 0x0A},
		core.DataRate100Hz:  {ctrl1: 0x0E},
		core.DataRate50Hz:   {ctrl1: 0x12},
		core.DataRate25Hz:   {ctrl1: 0x16},
		core.DataRate12_5Hz: {ctrl1: 0x1A},
	}
)

type fxas21002cGyroRange struct {
	ctrl0       byte
	ctrl3       byte
	sensitivity float32
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

type Configuration struct {
	Range     core.GyroRange
	rangeInfo *fxas21002cGyroRange
	Rate      core.DataRate
	rateInfo  *fxas21002cDataRate
}

type Fxas21002c struct {
	Bus embd.I2CBus
	// Addr of the sensor.
	address byte
	active  bool
	mu      sync.RWMutex

	conf *Configuration

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

func (d *Fxas21002c) Configure(conf *Configuration) error {
	if conf.Rate.ID > core.DataRate12_5Hz.ID {
		return errors.New("Lowest supported data rate is 12.5Hz")
	}
	conf.rangeInfo = fxas21002cGyroRanges[conf.Range]
	conf.rateInfo = fxas21002cDataRates[conf.Rate]
	d.conf = conf
	return nil
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
func (d *Fxas21002c) Start() error {
	if !d.active {
		if err := d.Activate(); err != nil {
			return err
		}
	}

	if d.readings == nil {
		// buffer is based on data rate
		d.readings = make(chan *core.GyroReading, d.conf.Rate.OneSecondBuffer)

		go func() {
			time.Sleep(fxas21002cActiveTransitionTime)
			ticker := time.NewTicker(d.conf.Rate.Period)
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
					d.readings = nil
					return
				}
			}
		}()
	}
	return nil
}

// Close stops any period reads in progress.  Call this to stop the readings that
// Start() begins
func (d *Fxas21002c) Close() error {
	if !d.active {
		return nil
	}
	if d.closing != nil {
		waitc := make(chan struct{})
		d.closing <- waitc
		<-waitc
	}
	if err := d.setup(); err != nil {
		return err
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

	// put sensor in standby mode
	err = d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl1, fxas21002cCtrlStandby)
	if err != nil {
		return errors.New("Unable to put fxas21002c in standby")
	}

	d.active = false
	return nil
}

func (d *Fxas21002c) Activate() error {

	if err := d.setup(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl0, d.conf.rangeInfo.ctrl0); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl3, d.conf.rangeInfo.ctrl3); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxas21002cRegisterCtrl1, d.conf.rateInfo.ctrl1); err != nil {
		return err
	}
	d.active = true
	return nil
}

func (d *Fxas21002c) ReadGyro() (*core.GyroReading, error) {
	if !d.active {
		return nil, errors.New("fxas21002c: sensor not active")
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	result := make([]byte, 7, 7)
	var err error
	err = d.Bus.ReadFromReg(d.address, fxas21002cRegisterStatus, result)
	if err != nil {
		return nil, err
	}
	sensitivity := d.conf.rangeInfo.sensitivity
	return &core.GyroReading{
		Xdps: float32(int16(binary.BigEndian.Uint16(result[1:3]))) * sensitivity,
		Ydps: float32(int16(binary.BigEndian.Uint16(result[3:5]))) * sensitivity,
		Zdps: float32(int16(binary.BigEndian.Uint16(result[5:7]))) * sensitivity,
	}, nil
}

func (d *Fxas21002c) validate() error {
	if d.Bus == nil {
		return errors.New("fxas21002c: no i2c bus")
	}
	if d.address == 0x00 {
		return fmt.Errorf("fxas21002c: invalid address %#x", d.address)
	}
	if d.conf == nil {
		return errors.New("fxas21001x: not configured")
	}
	return nil
}
