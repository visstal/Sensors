package datasource

import (
	"math/rand"
	sensor "sensor-node/internal/models"
	"strings"
	"time"
)

// SensorDataSource handles sensor data generation
type SensorDataSource struct {
	sensorName string
}

// NewSensorDataSource creates a new sensor data source
func NewSensorDataSource(sensorName string) *SensorDataSource {
	return &SensorDataSource{
		sensorName: sensorName,
	}
}

// GetReading generates and returns a new sensor reading with realistic mock values
func (ds *SensorDataSource) GetReading() *sensor.SensorReading {
	value := ds.generateMockValue()

	return &sensor.SensorReading{
		SensorName: ds.sensorName,
		Value:      value,
		Timestamp:  time.Now(),
	}
}

// generateMockValue creates realistic sensor values based on sensor type patterns
func (ds *SensorDataSource) generateMockValue() int {
	// Generate different realistic ranges based on common sensor types
	switch {
	case containsAny(ds.sensorName, []string{"temp", "temperature"}):
		// Temperature: -10 to 50 Celsius
		return rand.Intn(61) - 10
	case containsAny(ds.sensorName, []string{"humid", "moisture"}):
		// Humidity: 0 to 100 percent
		return rand.Intn(101)
	case containsAny(ds.sensorName, []string{"press", "pressure"}):
		// Pressure: 980 to 1050 hPa
		return rand.Intn(71) + 980
	case containsAny(ds.sensorName, []string{"cpu", "processor"}):
		// CPU usage: 0 to 100 percent
		return rand.Intn(101)
	case containsAny(ds.sensorName, []string{"volt", "voltage"}):
		// Voltage: 0 to 12 volts (scaled to int)
		return rand.Intn(121) // 0.0 to 12.0 volts (as 0 to 120)
	default:
		// Generic sensor: 0 to 100
		return rand.Intn(101)
	}
}

// containsAny checks if the sensor name contains any of the given keywords
func containsAny(text string, keywords []string) bool {
	lowerText := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lowerText, keyword) {
			return true
		}
	}
	return false
}
