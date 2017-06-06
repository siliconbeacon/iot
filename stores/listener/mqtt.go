package listener

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttListener struct {
	client  MQTT.Client
	closing chan chan struct{}
	weather chan WeatherReadings
}

func New(name string, server string) *MqttListener {
	options := MQTT.NewClientOptions().AddBroker(server)
	options.SetClientID(name).SetTLSConfig(newTlsConfig())

	mqttClient := MQTT.NewClient(options)
	return &MqttListener{
		client: mqttClient,
	}
}

func (l *MqttListener) Weather() <-chan WeatherReadings {
	return l.weather
}

func (l *MqttListener) Start(topic string) error {

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

func (l *MqttListener) Stop() {
	if l.closing != nil {
		waitc := make(chan struct{})
		l.closing <- waitc
		<-waitc
	}
}

func weatherListener(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("%s\n", msg.Payload())
}

func newTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("certs/tls_cert.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	} else {
		fmt.Println(err)
		return nil
	}
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
	}
}
