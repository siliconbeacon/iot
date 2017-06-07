package listener

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/siliconbeacon/iot/stores/listener/weather"
)

type Listener struct {
	client  MQTT.Client
	closing chan chan struct{}
	weather chan weather.Readings
}

func New(client MQTT.Client) {
	return &Listener{
		client: client,
	}
}

func (l *Listener) Weather() <-chan WeatherReadings {
	return l.weather
}

func (l *Listener) Start(topic string) error {

	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := l.client.Subscribe(topic, 0, weatherListener); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	l.weather = make(chan WeatherReadings, 5)
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

func weatherListener(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("%s\n", msg.Payload())
}
