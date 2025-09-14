// internal/services/telemetry_service.go
package services

import (
	"math"
	"math/rand"
	"time"

	"github.com/lucasdamasceno96/defender-vehicle/internal/models"
)

// TelemetryService is the interface for our telemetry operations.
type TelemetryService interface {
	GetTelemetry(page, limit int) []models.Telemetry
	GetAllTelemetry() []models.Telemetry
}

type telemetryServiceImpl struct {
	telemetryData []models.Telemetry
}

// NewTelemetryService creates a new instance of the telemetry service and generates data.
func NewTelemetryService() TelemetryService {
	service := &telemetryServiceImpl{}
	service.generateTelemetry(100) // Generate 100 data points on startup
	return service
}

func (s *telemetryServiceImpl) GetTelemetry(page, limit int) []models.Telemetry {
	start := (page - 1) * limit
	end := start + limit

	if start > len(s.telemetryData) {
		return []models.Telemetry{}
	}
	if end > len(s.telemetryData) {
		end = len(s.telemetryData)
	}
	return s.telemetryData[start:end]
}

func (s *telemetryServiceImpl) GetAllTelemetry() []models.Telemetry {
	return s.telemetryData
}

// generateTelemetry is now a private method of the service.
func (s *telemetryServiceImpl) generateTelemetry(count int) {
	// (O código da função GenerateTelemetry que criamos anteriormente vai aqui)
	// A única mudança é que ele agora popula s.telemetryData em vez de retornar um slice.
	// E usamos models.Telemetry em vez de Telemetry.

	s.telemetryData = make([]models.Telemetry, count)
	now := time.Now()

	startLat, startLon := 34.0522, -118.2437
	speed := 60.0
	heading := 45.0

	anomalyIndices := map[int]string{
		25: "gps_spoof",
		50: "speed_spike",
		75: "sensor_dropout",
	}

	for i := 0; i < count; i++ {
		// (Lógica de simulação de movimento idêntica à anterior)
		speed += (rand.Float64() - 0.5) * 2
		if speed < 40 {
			speed = 40
		}
		if speed > 80 {
			speed = 80
		}
		heading += (rand.Float64() - 0.5) * 1
		distance := (speed * 1000 / 3600) * 1
		earthRadius := 6371000.0
		radHeading := heading * (math.Pi / 180.0)
		latRad := startLat * (math.Pi / 180.0)
		lonRad := startLon * (math.Pi / 180.0)
		newLatRad := math.Asin(math.Sin(latRad)*math.Cos(distance/earthRadius) + math.Cos(latRad)*math.Sin(distance/earthRadius)*math.Cos(radHeading))
		newLonRad := lonRad + math.Atan2(math.Sin(radHeading)*math.Sin(distance/earthRadius)*math.Cos(latRad), math.Cos(distance/earthRadius)-math.Sin(latRad)*math.Sin(newLatRad))
		startLat = newLatRad * (180.0 / math.Pi)
		startLon = newLonRad * (180.0 / math.Pi)

		point := models.Telemetry{
			ID:        i, // Assign a simple index as ID
			Timestamp: now.Add(time.Duration(i) * time.Second),
			Lat:       startLat,
			Lon:       startLon,
			Speed:     speed,
			Heading:   heading,
			IsAnomaly: false,
		}

		if anomalyType, ok := anomalyIndices[i]; ok {
			point.IsAnomaly = true
			anomalyStr := anomalyType
			point.AnomalyType = &anomalyStr

			switch anomalyType {
			case "gps_spoof":
				point.Lat += 0.05
				point.Lon -= 0.05
			case "speed_spike":
				point.Speed *= 2.5
			case "sensor_dropout":
				point.Speed = 0
				point.Lat = 0
				point.Lon = 0
			}
		}
		s.telemetryData[i] = point
	}
}
