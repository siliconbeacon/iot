package publish

import (
	"fmt"
	//"time"

	//"github.com/golang/protobuf/proto"
	//"github.com/golang/protobuf/ptypes"

	"github.com/kidoman/embd"
	//"github.com/siliconbeacon/iot/messages"
	//"github.com/siliconbeacon/iot/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/siliconbeacon/iot/sensors/core"
	"github.com/siliconbeacon/iot/sensors/fxas21002c"
)

func Orientation(station string, i2cbus embd.I2CBus, mq MQTT.Client, shutdown chan bool) {
	gyro := fxas21002c.New(i2cbus)

	gyro.Configure(&fxas21002c.Configuration{
		Range: core.GyroRange500dps,
		Rate:  core.DataRate200Hz,
	})

	if present := gyro.IsPresent(); !present {
		fmt.Println("Unable to initialize FXAS21002C sensor.")
		return
	}

	fmt.Println("FXAS21002C gyroscope found. Initialized for a range of 500dps, reading at 200Hz.  Commencing Sensor Reads.")
	if err := gyro.Start(); err != nil {
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
