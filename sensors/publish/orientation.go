package publish

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/kidoman/embd"
	"github.com/siliconbeacon/iot/messages"
	"github.com/siliconbeacon/iot/mqtt"
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

	gyroReadings := gyro.Readings()
	accelReadings := accelmag.AccelReadings()
	magReadings := accelmag.MagReadings()

	var buffer [50]*messages.OrientationReading
	sampleCount := 0
	baseTime := time.Now().UTC()

	for {
		select {
		case gyroReading := <-gyroReadings:
			magReading := <-magReadings
			accelReading := <-accelReadings

			if sampleCount == 0 {
				baseTime = gyroReading.Timestamp
			}

			buffer[sampleCount] = &messages.OrientationReading{
				RelativeTimeUs: uint32(gyroReading.Timestamp.Sub(baseTime) / time.Microsecond),
				Gyroscope: &messages.GyroscopeReading{
					XDps: gyroReading.Xdps,
					YDps: gyroReading.Ydps,
					ZDps: gyroReading.Zdps,
				},
				Magnetometer: &messages.MagnetometerReading{
					XUt: magReading.XuT,
					YUt: magReading.YuT,
					ZUt: magReading.ZuT,
				},
				Accelerometer: &messages.AccelerometerReading{
					XMg: accelReading.Xmg,
					YMg: accelReading.Ymg,
					ZMg: accelReading.Zmg,
				},
			}
			sampleCount++
			if sampleCount == 50 {
				orientationBatch(station, mq, buffer[0:50], baseTime)
				sampleCount = 0
			}
		case <-shutdown:
			gyro.Close()
			accelmag.Close()
			return
		}
	}
}

func orientationBatch(station string, mq MQTT.Client, readings []*messages.OrientationReading, baseTime time.Time) error {
	basePTime, _ := ptypes.TimestampProto(baseTime)
	protobuf := &messages.OrientationReadings{
		Device:   station,
		BaseTime: basePTime,
		Readings: readings,
	}
	msg, _ := proto.Marshal(protobuf)
	topic := mqtt.CreateWeatherTopic(station)
	if token := mq.Publish(topic, 0, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
