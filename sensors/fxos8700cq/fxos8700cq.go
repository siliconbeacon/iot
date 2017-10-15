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

	fxos8700cqRegisterAccelSensitivity = 0x0E

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

type fxos8700cqAccelRange struct {
	ctrlSensitivity byte
	sensitivity     float64 // in mg per LSB
	ctrl1           byte    // enable low noise mode or not
}

var (
	fxos8700cqAccelRanges = map[core.AccelRange]*fxos8700cqAccelRange{
		2: {ctrlSensitivity: 0x00, sensitivity: 0.244, ctrl1: 0x04},
		4: {ctrlSensitivity: 0x01, sensitivity: 0.122, ctrl1: 0x04},
		8: {ctrlSensitivity: 0x02, sensitivity: 0.061, ctrl1: 0x00},
	}
)

type fxos8700cqRegisterValues struct {
	ctrl1    byte
	magctrl1 byte
	magctrl2 byte
}

type Configuration struct {
	AccelerometerEnabled bool
	AccelerometerRange   core.AccelRange
	accelRangeInfo       *fxos8700cqAccelRange
	MagnetometerEnabled  bool
	Rate                 core.DataRate
	registerConfigInfo   *fxos8700cqRegisterValues
}

type Fxos8700cq struct {
	Bus embd.I2CBus
	// Addr of the sensor.
	address byte

	active bool
	mu     sync.RWMutex

	conf *Configuration

	accelReadings chan *core.AccelReading
	magReadings   chan *core.MagReading
	closing       chan chan struct{}
}

// New creates a new fxos8700cq sensor.
func New(bus embd.I2CBus) *Fxos8700cq {
	return &Fxos8700cq{
		Bus:     bus,
		address: fxos8700cqAddrDefault,
	}
}

func (d *Fxos8700cq) Configure(conf *Configuration) error {
	if !(conf.AccelerometerEnabled || conf.MagnetometerEnabled) {
		return errors.New("Must enable at least one sensor")
	}

	// get configuration data for the data rate
	if conf.AccelerometerEnabled && conf.MagnetometerEnabled {
		// hybrid mode, both sensors enabled
		switch conf.Rate.ID {
		case core.DataRate400Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x01,
			}
		case core.DataRate200Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x09,
			}
		case core.DataRate100Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x11,
			}
		case core.DataRate50Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x19,
			}
		case core.DataRate25Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x21,
			}
		case core.DataRate6_25Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x29,
			}
		case core.DataRate3_125Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x31,
			}
		case core.DataRate0_78125Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x39,
			}
		default:
			return fmt.Errorf("Data rate %v is not supported in hybrid mode.", conf.Rate.Name)
		}

		conf.registerConfigInfo.magctrl1 = 0x1F
		conf.registerConfigInfo.magctrl2 = 0x20

	} else {
		// Single sensor
		switch conf.Rate.ID {
		case core.DataRate800Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x01,
			}
		case core.DataRate400Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x09,
			}
		case core.DataRate200Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x11,
			}
		case core.DataRate100Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x19,
			}
		case core.DataRate50Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x21,
			}
		case core.DataRate12_5Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x29,
			}
		case core.DataRate6_25Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x31,
			}
		case core.DataRate1_5625Hz.ID:
			conf.registerConfigInfo = &fxos8700cqRegisterValues{
				ctrl1: 0x39,
			}
		default:
			return fmt.Errorf("Data rate %v is not supported in single sensor mode.", conf.Rate.Name)
		}

		if conf.AccelerometerEnabled {
			conf.registerConfigInfo.magctrl1 = 0x1C
		} else {
			// Mag only
			conf.registerConfigInfo.magctrl1 = 0x1D
		}

		conf.registerConfigInfo.magctrl2 = 0x00
	}

	// Set a valid range info for the accelerometer
	if conf.AccelerometerEnabled {
		conf.accelRangeInfo = fxos8700cqAccelRanges[conf.AccelerometerRange]
	} else {
		conf.accelRangeInfo = fxos8700cqAccelRanges[core.AccelRange2g]
	}

	d.conf = conf
	return nil
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
func (d *Fxos8700cq) Start() error {
	if !d.active {
		if err := d.Activate(); err != nil {
			return err
		}
	}

	if d.accelReadings == nil && d.magReadings == nil {
		// buffer is based on data rate
		if d.conf.AccelerometerEnabled {
			d.accelReadings = make(chan *core.AccelReading, d.conf.Rate.OneSecondBuffer)
		}
		if d.conf.MagnetometerEnabled {
			d.magReadings = make(chan *core.MagReading, d.conf.Rate.OneSecondBuffer)
		}

		go func() {
			time.Sleep(d.conf.Rate.Period + 2*time.Millisecond)
			ticker := time.NewTicker(d.conf.Rate.Period)
			if d.conf.AccelerometerEnabled && d.conf.MagnetometerEnabled {
				for {
					select {
					case <-ticker.C:
						var (
							accelReading *core.AccelReading
							magReading   *core.MagReading
							err          error
						)

						accelReading, magReading, err = d.ReadAccelAndMag()
						if err != nil {
							glog.Errorf("fxos8700c: %v", err)
							continue
						}
						accelReading.Timestamp = time.Now().UTC()
						magReading.Timestamp = accelReading.Timestamp
						d.accelReadings <- accelReading
						d.magReadings <- magReading
					case waitc := <-d.closing:
						waitc <- struct{}{}
						ticker.Stop()
						close(d.accelReadings)
						close(d.magReadings)
						d.accelReadings = nil
						d.magReadings = nil
						return
					}
				}
			}
			if d.conf.AccelerometerEnabled {
				for {
					select {
					case <-ticker.C:
						var (
							accelReading *core.AccelReading
							err          error
						)

						accelReading, err = d.ReadAccel()
						if err != nil {
							glog.Errorf("fxos8700c: %v", err)
							continue
						}
						accelReading.Timestamp = time.Now().UTC()
						d.accelReadings <- accelReading
					case waitc := <-d.closing:
						waitc <- struct{}{}
						ticker.Stop()
						close(d.accelReadings)
						d.accelReadings = nil
						return
					}
				}
			}
			for {
				select {
				case <-ticker.C:
					var (
						magReading *core.MagReading
						err        error
					)

					magReading, err = d.ReadMag()
					if err != nil {
						glog.Errorf("fxos8700c: %v", err)
						continue
					}
					magReading.Timestamp = time.Now().UTC()
					d.magReadings <- magReading
				case waitc := <-d.closing:
					waitc <- struct{}{}
					ticker.Stop()
					close(d.magReadings)
					d.magReadings = nil
					return
				}
			}
		}()
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

	d.active = false
	return nil
}

func (d *Fxos8700cq) Activate() error {
	if err := d.setup(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Bus.WriteByteToReg(d.address, fxos8700cqRegisterMCtrl1, d.conf.registerConfigInfo.magctrl1); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxos8700cqRegisterMCtrl2, d.conf.registerConfigInfo.magctrl2); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxos8700cqRegisterAccelSensitivity, d.conf.accelRangeInfo.ctrlSensitivity); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.address, fxos8700cqRegisterCtrl1, d.conf.accelRangeInfo.ctrl1|d.conf.registerConfigInfo.ctrl1); err != nil {
		return err
	}

	d.active = true
	return nil
}

func (d *Fxos8700cq) ReadAccel() (*core.AccelReading, error) {
	if !d.active || !d.conf.AccelerometerEnabled {
		return nil, errors.New("fxos8700cq: accelerometer not active")
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	result := make([]byte, 7, 7)
	var err error
	err = d.Bus.ReadFromReg(d.address, fxos8700cqRegisterAccelStatus, result)
	if err != nil {
		return nil, err
	}

	return newAccel(result[1:7], d.conf.accelRangeInfo.sensitivity), nil
}

func newAccel(accelBytes []byte, sensitivity float64) *core.AccelReading {
	return &core.AccelReading{
		Xmg: float64(int16(binary.BigEndian.Uint16(accelBytes[0:2]))) * sensitivity,
		Ymg: float64(int16(binary.BigEndian.Uint16(accelBytes[2:4]))) * sensitivity,
		Zmg: float64(int16(binary.BigEndian.Uint16(accelBytes[4:6]))) * sensitivity,
	}
}

func (d *Fxos8700cq) ReadMag() (*core.MagReading, error) {
	if !d.active || !d.conf.MagnetometerEnabled {
		return nil, errors.New("fxos8700cq: magnetometer not active")
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	result := make([]byte, 7, 7)
	var err error
	err = d.Bus.ReadFromReg(d.address, fxos8700cqRegisterMagStatus, result)
	if err != nil {
		return nil, err
	}

	return newMag(result[1:7]), nil
}

func newMag(magBytes []byte) *core.MagReading {
	return &core.MagReading{
		XuT: float64(int16(binary.BigEndian.Uint16(magBytes[0:2]))) * 0.1,
		YuT: float64(int16(binary.BigEndian.Uint16(magBytes[2:4]))) * 0.1,
		ZuT: float64(int16(binary.BigEndian.Uint16(magBytes[4:6]))) * 0.1,
	}
}

func (d *Fxos8700cq) ReadAccelAndMag() (*core.AccelReading, *core.MagReading, error) {
	if !d.active || !d.conf.AccelerometerEnabled || !d.conf.MagnetometerEnabled {
		return nil, nil, errors.New("fxos8700cq: not active in dual sensor mode")
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	result := make([]byte, 13, 13)
	var err error
	err = d.Bus.ReadFromReg(d.address, fxos8700cqRegisterAccelStatus, result)
	if err != nil {
		return nil, nil, err
	}

	return newAccel(result[1:7], d.conf.accelRangeInfo.sensitivity), newMag(result[7:13]), nil
}

// Close stops any period reads in progress.  Call this to stop the readings that
// Start() begins
func (d *Fxos8700cq) Close() error {
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

func (d *Fxos8700cq) validate() error {
	if d.Bus == nil {
		return errors.New("fxos8700cq: no i2c bus")
	}
	if d.address == 0x00 {
		return fmt.Errorf("fxos8700cq: invalid address %#x", d.address)
	}
	return nil
}
