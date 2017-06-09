package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/siliconbeacon/iot/mqtt"
	"github.com/siliconbeacon/iot/stores/listener"
	"github.com/siliconbeacon/iot/stores/listener/weather"
)

func main() {
	osExit := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	signal.Notify(osExit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Connecting to MQTT Server...")
	mqClient := mqtt.NewTLSClient(fmt.Sprintf("store.riakts.%v", os.Getpid()), "ssl://localhost:8883/")
	l := listener.New(mqClient)
	if err := l.Start(mqtt.WeatherTopicPattern); err != nil {
		fmt.Println(err)
		return
	}
	go writeToRiak(l.Weather(), shutdown)
	fmt.Println("Running...  Press Ctrl-C to exit.")
	<-osExit
	fmt.Println("Exiting.")
	l.Stop()
	close(shutdown)
}

func writeToRiak(readings <-chan weather.Readings, shutdown chan bool) {
	for {
		select {
		case reading := <-readings:
			fmt.Println(reading[0].Station, reading[0].TemperatureDegreesCelsius)
		case <-shutdown:
			return
		}
	}
}
