// internal/models/telemetry.go
package models

import "time"

// Telemetry holds a single data point for the vehicle's state.
type Telemetry struct {
	ID          int       `json:"id"`
	Timestamp   time.Time `json:"ts"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Speed       float64   `json:"speed"`   // in km/h
	Heading     float64   `json:"heading"` // in degrees
	IsAnomaly   bool      `json:"is_anomaly"`
	AnomalyType *string   `json:"anomaly_type,omitempty"`
}
