package listener

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/siliconbeacon/iot/messages"
	"github.com/siliconbeacon/iot/stores/listener/weather"
)

type Listener struct {
	client  MQTT.Client
	closing chan chan struct{}
	weather chan weather.Readings
}

func New(client MQTT.Client) *Listener {
	return &Listener{
		client: client,
	}
}

func (l *Listener) Weather() <-chan weather.Readings {
	return l.weather
}

func (l *Listener) Start(topic string) error {
	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	var subscriber MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		readings := deserializeMessage(msg.Payload())
		if readings != nil {
			l.weather <- readings
		}
	}
	if token := l.client.Subscribe(topic, 0, subscriber); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	l.weather = make(chan weather.Readings, 5)
	go func() {
		for {
			select {
			case waitc := <-l.closing:
				l.client.Unsubscribe(topic)
				l.client.Disconnect(250)
				close(l.weather)
				waitc <- struct{}{}
				return
			}
		}
	}()
	return nil
}

func (l *Listener) Stop() {
	if l.closing != nil {
		waitc := make(chan struct{})
		l.closing <- waitc
		<-waitc
	}
}

func deserializeMessage(body []byte) weather.Readings {
	pb := &messages.WeatherReadings{}
	if err := proto.Unmarshal(body, pb); err != nil {
		return nil
	}
	var readings weather.Readings

	baseTime, _ := ptypes.Timestamp(pb.BaseTime)
	for _, reading := range pb.Readings {
		delta := time.Duration(reading.RelativeTimeUs) * time.Microsecond
		readings = append(readings, &weather.Reading{
			Station:                    pb.Device,
			Timestamp:                  baseTime.Add(delta),
			TemperatureDegreesCelsius:  reading.TemperatureDegreesC,
			RelativeHumidityPercentage: reading.RelativeHumidityPercentage,
		})
	}
	return readings
}
