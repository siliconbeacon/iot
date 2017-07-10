package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	riak "github.com/basho/riak-go-client"

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
	cluster := connectToRiak("localhost")
	defer cluster.Stop()
	for {
		select {
		case reading := <-readings:
			fmt.Println(reading[0].Station, reading[0].TemperatureDegreesCelsius)
			rows := make([][]riak.TsCell, len(reading))
			for i, r := range reading {
				row := []riak.TsCell{
					riak.NewStringTsCell(r.Station),
					riak.NewTimestampTsCell(r.Timestamp),
					riak.NewDoubleTsCell(r.TemperatureDegreesCelsius),
					riak.NewDoubleTsCell(r.RelativeHumidityPercentage),
				}
				rows[i] = row
			}
			cmd, err := riak.NewTsStoreRowsCommandBuilder().WithTable("Weather").WithRows(rows).Build()

			if err != nil {
				fmt.Println("Error writing to Riak...", err)
				continue
			}
			err = cluster.Execute(cmd)
			if err != nil {
				fmt.Println("Error writing to Riak...", err)
			}

		case <-shutdown:
			return
		}
	}
}

func connectToRiak(server string) *riak.Cluster {
	addr := fmt.Sprintf("%s:8087", server)
	nodeOpts := &riak.NodeOptions{
		MinConnections: 10,
		MaxConnections: 256,
		RemoteAddress:  addr,
	}
	if node, err := riak.NewNode(nodeOpts); err != nil {
		fmt.Println("Riak config issue...", err)
	} else {
		clusterOpts := &riak.ClusterOptions{
			Nodes:             []*riak.Node{node},
			ExecutionAttempts: 3,
		}
		cluster, err := riak.NewCluster(clusterOpts)
		if err != nil {
			fmt.Println("Riak issue...", err)
			return nil
		}
		err = cluster.Start()
		if err != nil {
			fmt.Println("Riak issue...", err)
			return nil
		}
		return cluster
	}
	return nil
}
