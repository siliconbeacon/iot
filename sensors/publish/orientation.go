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
	"github.com/siliconbeacon/iot/sensors/fxos8700cq"
)

func Orientation(station string, i2cbus embd.I2CBus, mq MQTT.Client, shutdown chan bool) {
	accelmag := fxos8700cq.New(i2cbus)

	accelmag.Configure(&fxos8700cq.Configuration{
		AccelerometerEnabled: true,
		AccelerometerRange:   core.AccelRange4g,
		MagnetometerEnabled:  true,
		Rate:                 core.DataRate200Hz,
	})

	if present := accelmag.IsPresent(); !present {
		fmt.Println("Unable to initialize FXOS8700CQ sensor.")
		return
	}

	fmt.Println("FXOS8700CQ accelerometer and magnetometer found. Initialized for an acceleromter range of 4g, reading at 200Hz.  Commencing Sensor Reads.")

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
