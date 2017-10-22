package publish

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/kidoman/embd"
	"github.com/siliconbeacon/iot/messages"
	"github.com/siliconbeacon/iot/mqtt"
	"github.com/siliconbeacon/iot/sensors/core"
	"github.com/siliconbeacon/iot/sensors/si70xx"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Weather(station string, i2cbus embd.I2CBus, mq MQTT.Client, shutdown chan bool) {

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

	// let's sample at 12.5Hz
	if err = sensor.Start(core.DataRate12_5Hz); err != nil {
		fmt.Println("Unable to commence sensor reads from Si70xx.")
		return
	}
	readings := sensor.Readings()
	var buffer [25]*core.TemperatureAndHumidityReading
	sampleCount := 0
	for {
		select {
		case reading := <-readings:
			buffer[sampleCount] = reading
			sampleCount++
			if sampleCount == 25 {
				weatherBatch(station, mq, buffer[0:25])
				sampleCount = 0
			}
		case <-shutdown:
			sensor.Close()
			if sampleCount > 0 {
				weatherBatch(station, mq, buffer[0:sampleCount])
			}
			return
		}
	}
}

func weatherBatch(station string, mq MQTT.Client, readings []*core.TemperatureAndHumidityReading) error {
	var msg []byte
	var err error
	if msg, err = serializeWeather(station, readings); err != nil {
		return err
	}
	topic := mqtt.CreateWeatherTopic(station)
	if token := mq.Publish(topic, 0, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func serializeWeather(station string, readings []*core.TemperatureAndHumidityReading) ([]byte, error) {
	baseTime := readings[0].Timestamp
	basePTime, _ := ptypes.TimestampProto(baseTime)
	msg := &messages.WeatherReadings{
		Device:   station,
		BaseTime: basePTime,
	}
	for _, reading := range readings {
		msg.Readings = append(msg.Readings, &messages.WeatherReading{
			RelativeTimeUs:             uint32(reading.Timestamp.Sub(baseTime) / time.Microsecond),
			TemperatureDegreesC:        reading.TemperatureDegreesCelsius,
			RelativeHumidityPercentage: reading.RelativeHumidityPercentage,
		})
	}
	return proto.Marshal(msg)
}
