// Package si70xx is a driver for the Si70xx temperature / humidity sensor
package si70xx

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
	si70xxAddrDefault                = 0x40
	si70xxCmdReadRhtRegister         = 0xE7
	si70xxCmdMeasureRelativeHumidity = 0xF5
	si70xxCmdMeasureTemperature      = 0xF3
	si70xxCmdReadPreviousTemperature = 0xE0
	si70xxCmdReset                   = 0xFE
	si70xxRhtRegisterValue           = 0x3A
)

var (
	si70xxCmdIdentify1       = []byte{0xFA, 0x0F}
	si70xxCmdIdentify2       = []byte{0xFC, 0xC9}
	si70xxCmdFirmwareVersion = []byte{0x84, 0xB8}
	si70xxFirmwareVersions   = map[byte]string{
		0x20: "2.0",
		0xFF: "1.0",
	}
	si70xxChipModels = map[byte]string{
		0x00: "Engineering Sample",
		0xFF: "Engineering Sample",
		0x0D: "Si7013",
		0x14: "Si7020",
		0x15: "Si7021",
	}
)

// Si70xx represents The Si70xx sensor (Si7013, Si7020, or Si7021)
type Si70xx struct {
	// Bus to communicate over.
	Bus embd.I2CBus
	// Addr of the sensor.
	address byte

	initialized bool
	mu          sync.RWMutex

	readings chan *core.TemperatureAndHumidityReading
	closing  chan chan struct{}
	serial   string
	firmware string
	model    string
}

// New creates a new si70xx sensor.
func New(bus embd.I2CBus) *Si70xx {
	return &Si70xx{
		Bus:     bus,
		address: si70xxAddrDefault,
	}
}

// IsPresent returns true if it looks like we were able to see the sensor
func (d *Si70xx) IsPresent() bool {
	if err := d.setup(); err != nil {
		return false
	}
	return true
}

// SerialNumber returns the sensor's hardware serial number
func (d *Si70xx) SerialNumber() (string, error) {
	if err := d.setup(); err != nil {
		return "", err
	}
	return d.serial, nil
}

// FirmwareVersion returns the sensor's firmware version
func (d *Si70xx) FirmwareVersion() (string, error) {
	if err := d.setup(); err != nil {
		return "", err
	}
	return d.firmware, nil
}

// ModelName returns the model name of the sensor (Si7013, Si7020, or Si7021)
func (d *Si70xx) ModelName() (string, error) {
	if err := d.setup(); err != nil {
		return "", err
	}
	return d.model, nil
}

// Humdity performs a Relative Humidity reading
func (d *Si70xx) Humidity() (float64, error) {
	if err := d.setup(); err != nil {
		return -1, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.Bus.WriteByte(d.address, si70xxCmdMeasureRelativeHumidity); err != nil {
		return -1, errors.New("Error asking sensor for humidity reading")
	}
	time.Sleep(25 * time.Millisecond)
	reading, err := d.Bus.ReadBytes(d.address, 3)
	if err != nil {
		return -1, errors.New("Error reading humidity value from sensor")
	}
	rawValue := binary.BigEndian.Uint16(reading[0:2])
	return float64(rawValue)*125.0/65536.0 - 6.0, nil
}

// Temperature performs a Temperature reading
func (d *Si70xx) Temperature() (float64, error) {
	if err := d.setup(); err != nil {
		return -1, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.Bus.WriteByte(d.address, si70xxCmdMeasureTemperature); err != nil {
		return -1, errors.New("Error asking sensor for temperature reading")
	}
	time.Sleep(15 * time.Millisecond)
	reading, err := d.Bus.ReadBytes(d.address, 3)
	if err != nil {
		return -1, errors.New("Error reading temperature value from sensor")
	}

	return calcTemperature(reading), nil
}

// LastTemperature returns a temperature value that was measured during the last Humidity reading,
// without waiting for the sensor to perform a new reading.
func (d *Si70xx) LastTemperature() (float64, error) {
	if err := d.setup(); err != nil {
		return -1, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.Bus.WriteByte(d.address, si70xxCmdReadPreviousTemperature); err != nil {
		return -1, errors.New("Error asking sensor for temperature reading")
	}
	reading, err := d.Bus.ReadBytes(d.address, 2)
	if err != nil {
		return -1, errors.New("Error reading temperature value from sensor")
	}

	return calcTemperature(reading), nil
}

// Readings is a channel that will contain sensor readings after calling Start()
func (d *Si70xx) Readings() <-chan *core.TemperatureAndHumidityReading {
	return d.readings
}

// Start produces a stream of humidity and temperature readings
// in the Readings() channel
func (d *Si70xx) Start(rate core.DataRate) error {
	if err := d.setup(); err != nil {
		return err
	}
	if rate.ID < core.DataRate25Hz.ID {
		return errors.New("Cannot sample at a higher frequency than 25Hz")
	}
	// determine buffer size based on frequency.
	bufferSize := 2
	if rate.Period < time.Second {
		bufferSize = int(time.Second/rate.Period) + 1
	}

	d.readings = make(chan *core.TemperatureAndHumidityReading, bufferSize)

	go func() {
		ticker := time.NewTicker(rate.Period)
		for {
			select {
			case <-ticker.C:
				var err error
				reading := &core.TemperatureAndHumidityReading{}
				if reading.RelativeHumidityPercentage, err = d.Humidity(); err != nil {
					glog.Errorf("si70xx: %v", err)
					continue
				}
				reading.Timestamp = time.Now().UTC()
				if reading.TemperatureDegreesCelsius, err = d.LastTemperature(); err != nil {
					glog.Errorf("si70xx: %v", err)
					continue
				}
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
func (d *Si70xx) Close() error {
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

func calcTemperature(reading []byte) float64 {
	rawValue := binary.BigEndian.Uint16(reading[0:2])
	return float64(rawValue)*175.72/65536.0 - 46.85
}

func (d *Si70xx) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.validate(); err != nil {
		return err
	}

	// issue reset
	if err := d.Bus.WriteByte(d.address, si70xxCmdReset); err != nil {
		return errors.New("Unable to reset si70xx")
	}
	time.Sleep(50 * time.Millisecond)

	// make sure this looks like the right sensor
	reg, err := d.Bus.ReadByteFromReg(d.address, si70xxCmdReadRhtRegister)
	if err != nil {
		return err
	}
	if reg != si70xxRhtRegisterValue {
		return fmt.Errorf("si70xx: Unexpected RHT Register Value: %#x", reg)
	}
	fmt.Printf("Reg: %#x\n", reg)

	// read firmware version
	if err := d.Bus.WriteByteToReg(d.address, si70xxCmdFirmwareVersion[0], si70xxCmdFirmwareVersion[1]); err != nil {
		return errors.New("Unable to read Firmware Version")
	}

	firmware, err := d.Bus.ReadByte(d.address)
	if err != nil {
		return errors.New("Unable to read Firmware Version")
	}
	d.firmware = si70xxFirmwareVersions[firmware]

	// read serial number and firmware version since we're connected
	if err := d.Bus.WriteByteToReg(d.address, si70xxCmdIdentify1[0], si70xxCmdIdentify1[1]); err != nil {
		return errors.New("Unable to read Serial Number (Part 1)")
	}
	serial1, err := d.Bus.ReadBytes(d.address, 8)
	if err != nil {
		return errors.New("Unable to read Serial Number (Part 1)")
	}
	if err := d.Bus.WriteByteToReg(d.address, si70xxCmdIdentify2[0], si70xxCmdIdentify2[1]); err != nil {
		return errors.New("Unable to read Serial Number (Part 2)")
	}
	serial2, err := d.Bus.ReadBytes(d.address, 8)
	if err != nil {
		return errors.New("Unable to read Serial Number (Part 2)")
	}
	d.serial = fmt.Sprintf("%x%x%x%x%x%x%x%x", serial1[0], serial1[2], serial1[4], serial1[6], serial2[0], serial2[2], serial2[4], serial2[6])
	d.model = si70xxChipModels[serial2[0]] // derived from serial number
	d.initialized = true

	return nil
}

func (d *Si70xx) validate() error {
	if d.Bus == nil {
		return errors.New("si70xx: no i2c bus")
	}
	if d.address == 0x00 {
		return fmt.Errorf("si70xx: invalid address %#x", d.address)
	}
	return nil
}
