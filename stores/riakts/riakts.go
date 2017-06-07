package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/siliconbeacon/iot/mqtt"
	"github.com/siliconbeacon/iot/stores/listener"
)

func main() {
	osExit := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	signal.Notify(osExit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Connecting to MQTT Server...")
	mqClient := mqtt.NewTLSClient(fmt.Sprintf("store.riakts.%v", os.Getpid()), "ssl://localhost:8883/")
	listener := listener.New(mqClient)
	if err := listener.Start(mqtt.WeatherTopicPattern); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Running...  Press Ctrl-C to exit.")
	<-osExit
	fmt.Println("Exiting.")
	listener.Stop()
	close(shutdown)
}
