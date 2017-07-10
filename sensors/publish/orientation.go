package publish

import (
	"fmt"
	//"time"

	//"github.com/golang/protobuf/proto"
	//"github.com/golang/protobuf/ptypes"

	"github.com/kidoman/embd"
	//"github.com/siliconbeacon/iot/messages"
	//"github.com/siliconbeacon/iot/mqtt"
	"github.com/siliconbeacon/iot/sensors/fxas21002c"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Orientation(station string, i2cbus embd.I2CBus, mq MQTT.Client, shutdown chan bool) {
	gyro := fxas21002c.New(i2cbus)

	if present := gyro.IsPresent(); !present {
		fmt.Println("Unable to initialize FXAS21002C sensor.")
		return
	}

	gyroRange := fxas21002c.GyroRange500dps
	gyroRate := fxas21002c.DataRate200Hz

	fmt.Println("FXAS21002C gyroscope found.  Initializing for a range of %v dps, reading at 200Hz.", gyroRange)
	if err := gyro.Start(gyroRange, gyroRate); err != nil {
		fmt.Println("Unable to commence sensor reads from FXAS21002C sensor.")
		return
	}

	readings := gyro.Readings()
	for {
		select {
		case reading := <-readings:
			fmt.Printf("x: %v, y: %v, z: %v\n", reading.Xdps, reading.Ydps, reading.Zdps)
		case <-shutdown:
			gyro.Close()
			return
		}
	}
}
