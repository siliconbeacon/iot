package mqtt

import (
	"strings"
)

const (
	// WeatherTopicPattern is the topic pattern to listen to for readings
	// from all weather stations
	WeatherTopicPattern = "weather/#/readings"
)

// CreateWeatherTopic creates a specific topic for weather readings from
// the named weather station
func CreateWeatherTopic(station string) string {
	return strings.Replace(WeatherTopicPattern, "#", station, 1)
}
