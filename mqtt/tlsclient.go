package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// NewTLSClient returns a new MQTT client, with TLS configured
// already.
func NewTLSClient(name string, server string) MQTT.Client {
	options := MQTT.NewClientOptions().AddBroker(server)
	options.SetClientID(name).SetTLSConfig(newTLSConfig())
	return MQTT.NewClient(options)
}

func newTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("/etc/ssl/cert.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	} else {
		fmt.Println(err)
		return nil
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}
}
