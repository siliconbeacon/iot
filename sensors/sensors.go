package main

import "fmt"

import "os"
import "os/signal"
import "syscall"
import "time"

import "github.com/kidoman/embd"
import _ "github.com/kidoman/embd/host/rpi"

import "github.com/siliconbeacon/iot/sensors/si70xx"

const (
	rpiI2cBus = 1
)

func runTemperatureSensor(i2cbus embd.I2CBus, shutdown chan bool) {
	sensor := si70xx.New(i2cbus)
	var (
		err      error
		serial   string
		model    string
		firmware string
	)
	if serial, err = sensor.SerialNumber(); err != nil {
		fmt.Println("Unable to initialize Si70xx sensor.")
		return
	}
	if model, err = sensor.ModelName(); err != nil {
		fmt.Println("Unable to initialize Si70xx sensor.")
		return
	}
	if firmware, err = sensor.FirmwareVersion(); err != nil {
		fmt.Println("Unable to initialize Si70xx sensor.")
		return
	}
	fmt.Println(model, "found. Serial:", serial, "running firmware version", firmware)

	// let's sample at 5Hz
	if err = sensor.Start(200 * time.Millisecond); err != nil {
		fmt.Println("Unable to commence sensor reads from Si70xx.")
		return
	}
	readings := sensor.Readings()
	var buffer [10]*si70xx.TemperatureAndHumidity
	sampleCount := 0
	for {
		select {
		case reading := <-readings:
			buffer[sampleCount] = reading
			sampleCount++
			if sampleCount == 10 {
				fmt.Println("Got 10.")
				// do mqtt stuff with buffer(0:sampleCount)
				for i := range buffer {
					fmt.Println(buffer[i].TemperatureDegreesCelsius)
					buffer[i] = nil
				}
				sampleCount = 0
			}
		case <-shutdown:
			sensor.Close()
			// do mqtt stuff with buffer(0:sampleCount)
			return
		}
	}
}

func main() {
	osExit := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	signal.Notify(osExit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Initializing Sensors...")
	embd.InitI2C()
	defer embd.CloseI2C()
	i2cbus := embd.NewI2CBus(rpiI2cBus)

	// start all sensors
	go runTemperatureSensor(i2cbus, shutdown)

	fmt.Println("Running...  Press Ctrl-C to exit.")

	<-osExit
	fmt.Println("Exiting.")
	close(shutdown)
}

// func blink() {
// 	osExit := make(chan os.Signal, 1)
// 	shutdown := make(chan bool, 1)
// 	signal.Notify(osExit, syscall.SIGINT, syscall.SIGTERM)

// 	fmt.Println("Initializing Sensors...")

// 	embd.InitGPIO()
// 	defer embd.CloseGPIO()

// 	green, _ := embd.NewDigitalPin(5)
// 	green.SetDirection(embd.Out)

// 	yellow, _ := embd.NewDigitalPin(6)
// 	yellow.SetDirection(embd.Out)

// 	blue, _ := embd.NewDigitalPin(19)
// 	blue.SetDirection(embd.Out)

// 	red, _ := embd.NewDigitalPin(26)
// 	red.SetDirection(embd.Out)

// 	fmt.Println("Blinking.  Press Ctrl-C to exit.")

// 	go func() {
// 		ticker := time.NewTicker(time.Millisecond * 200)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				red.Write(embd.High)
// 				time.Sleep(time.Millisecond * 100)
// 				red.Write(embd.Low)
// 			case <-shutdown:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()

// 	go func() {
// 		ticker := time.NewTicker(time.Millisecond * 350)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				blue.Write(embd.High)
// 				time.Sleep(time.Millisecond * 175)
// 				blue.Write(embd.Low)
// 			case <-shutdown:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()

// 	go func() {
// 		ticker := time.NewTicker(time.Millisecond * 500)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				green.Write(embd.High)
// 				time.Sleep(time.Millisecond * 250)
// 				green.Write(embd.Low)
// 			case <-shutdown:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()

// 	go func() {
// 		ticker := time.NewTicker(time.Millisecond * 1200)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				yellow.Write(embd.High)
// 				time.Sleep(time.Millisecond * 600)
// 				yellow.Write(embd.Low)
// 			case <-shutdown:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()

// 	<-osExit
// 	fmt.Println("Exiting.")
// 	close(shutdown)

// }
