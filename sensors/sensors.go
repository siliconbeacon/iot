package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"

	"github.com/siliconbeacon/iot/mqtt"
	"github.com/siliconbeacon/iot/sensors/publish"
)

const (
	rpiI2cBus  = 1
	deviceName = "sensor.pi3650"
)

func main() {
	osExit := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	signal.Notify(osExit, syscall.SIGINT, syscall.SIGTERM)

	mqClient := mqtt.NewTLSClient(deviceName, "ssl://upsilon:8883")
	fmt.Println("Connecting to Server...")

	if token := mqClient.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		fmt.Println("Unable to connect to MQTT server.")
		return
	}

	defer mqClient.Disconnect(250)

	fmt.Println("Initializing Sensors...")
	embd.InitI2C()
	defer embd.CloseI2C()
	i2cbus := embd.NewI2CBus(rpiI2cBus)

	// start all sensors
	go publish.Weather(deviceName, i2cbus, mqClient, shutdown)

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
