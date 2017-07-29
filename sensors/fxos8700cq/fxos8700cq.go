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

	"github.com/siliconbeacon/iot/sensors/core"
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

	fxos8700cqCtrlStandby = 0x00
)

type Fxos8700cq struct {
	Bus embd.I2CBus
	// Addr of the sensor.
	address     byte
	initialized bool
	mu          sync.RWMutex

	accelReadings chan *core.AccelReading
	magReadings   chan *core.MagReading
	closing       chan chan struct{}
}

// New creates a new fxas21002c sensor.
func New(bus embd.I2CBus) *Fxos8700cq {
	return &Fxos8700cq{
		Bus:     bus,
		address: fxos8700cqAddrDefault,
	}
}

// IsPresent returns true if it looks like we were able to see the sensor
func (d *Fxos8700cq) IsPresent() bool {
	if err := d.setup(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// AccelReadings is a channel that will contain sensor readings after calling Start()
func (d *Fxos8700cq) AccelReadings() <-chan *core.AccelReading {
	return d.accelReadings
}

// MagReadings is a channel that will contain sensor readings after calling Start()
func (d *Fxos8700cq) MagReadings() <-chan *core.MagReading {
	return d.magReadings
}

// Start produces a stream of gyroscope readings in the Readings() channel
func (d *Fxos8700cq) Start(rate core.DataRate) error {
	return nil
}

// Close stops any period reads in progress.  Call this to stop the readings that
// Start() begins
func (d *Fxos8700cq) Close() error {
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

func (d *Fxos8700cq) setup() error {

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.validate(); err != nil {
		return err
	}

	// check whoami register first
	whoami, err := d.Bus.ReadByteFromReg(d.address, fxos8700cqRegisterWhoami)
	if err != nil {
		return err
	}
	if whoami != fxos8700cqIdentifier {
		return fmt.Errorf("fxos8700cq whoami check failed.  Expected %#x, got %#x", fxos8700cqIdentifier, whoami)
	}

	// put sensor in standby mode
	err = d.Bus.WriteByteToReg(d.address, fxos8700cqRegisterCtrl1, fxos8700cqCtrlStandby)
	if err != nil {
		return errors.New("Unable to reset fxos8700cq")
	}

	return nil
}

func (d *Fxos8700cq) Activate() error {

}

func (d *Fxos8700cq) ReadAccel() (*core.AccelReading, error) {

}

func (d *Fxos8700cq) ReadMag() (*core.MagReading, error) {

}

func (d *Fxos8700cq) ReadAccelAndMag() (*core.AccelReading, *core.MagReading, error) {

}

func (d *Fxos8700cq) validate() error {
	if d.Bus == nil {
		return errors.New("fxos8700cq: no i2c bus")
	}
	if d.address == 0x00 {
		return fmt.Errorf("fxos8700cq: invalid address %#x", d.address)
	}
	return nil
}
