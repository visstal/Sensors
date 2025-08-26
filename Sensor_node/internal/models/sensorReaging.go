package sensor

import "time"

// SensorReading represents a single sensor measurement
type SensorReading struct {
	SensorName string    `json:"sensor_name"`
	Value      int       `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}
