// internal/services/telemetry_service.go
package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/lucasdamasceno96/defender-vehicle/internal/models"
)

// TelemetryService is the interface for our telemetry operations.
type TelemetryService interface {
	GetTelemetry(page, limit int) []models.Telemetry
	GetAllTelemetry() []models.Telemetry
	ProcessDetection(req models.DetectionRequest) (bool, error)
	TriggerMitigation(telemetryID int) (models.MitigationResponse, error)
	ExportLogsToCSV(writer io.Writer) error
	GetGameState() models.GameState
}

type telemetryServiceImpl struct {
	telemetryData []models.Telemetry
	detectionLogs []models.DetectionLog
	playerScore   int
	playerBadges  map[string]bool // Usamos um map para evitar badges duplicados
	correctHits   int
}

// NewTelemetryService creates a new instance of the telemetry service and generates data.
func NewTelemetryService() TelemetryService {
	service := &telemetryServiceImpl{
		detectionLogs: make([]models.DetectionLog, 0),
		playerScore:   0,
		playerBadges:  make(map[string]bool),
		correctHits:   0,
	}
	service.generateTelemetry(100)
	return service
}

func (s *telemetryServiceImpl) GetTelemetry(page, limit int) []models.Telemetry {
	// For the dashboard, we'll send all data at once.
	// In a real-world scenario with millions of points, we would stream or paginate.
	if limit > len(s.telemetryData) {
		limit = len(s.telemetryData)
	}
	return s.telemetryData[0:limit]
}

func (s *telemetryServiceImpl) GetAllTelemetry() []models.Telemetry {
	return s.telemetryData
}

// generateTelemetry is now a private method of the service.
func (s *telemetryServiceImpl) generateTelemetry(count int) {
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
			ID:        i,
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
} // <<< CORREÇÃO AQUI: Esta chave estava faltando!

var mitigationPlaybook = map[string][]string{
	"gps_spoof": {
		"Cross-validate with other sensors (IMU, wheel speed)",
		"Increase reliance on inertial navigation for short-term",
		"Alert SOC and flag data as untrusted",
		"Enter fail-operational mode (limp mode)",
	},
	"speed_spike": {
		"Validate against wheel encoders and accelerometer data",
		"Filter sensor readings using a Kalman filter",
		"Isolate the faulty sensor",
		"Trigger alert for sensor malfunction",
	},
	"sensor_dropout": {
		"Switch to redundant sensor unit",
		"Activate safe pull-over procedure",
		"Alert SOC about critical sensor failure",
		"Log system health status for maintenance",
	},
}

func (s *telemetryServiceImpl) ProcessDetection(req models.DetectionRequest) (bool, error) {
	if req.TelemetryID < 0 || req.TelemetryID >= len(s.telemetryData) {
		return false, fmt.Errorf("telemetry ID %d out of bounds", req.TelemetryID)
	}

	point := s.telemetryData[req.TelemetryID]
	isCorrect := point.IsAnomaly

	if isCorrect {
		s.playerScore += 10
		s.correctHits++
		s.checkAndAwardBadges() // <<< Lógica de Badges
	} else {
		s.playerScore -= 5
	}

	log := models.DetectionLog{
		Timestamp:      time.Now(),
		TelemetryID:    req.TelemetryID,
		UserGuess:      true,
		Correct:        isCorrect,
		ResponseTimeMs: 250,
	}
	s.detectionLogs = append(s.detectionLogs, log)

	return isCorrect, nil
}

func (s *telemetryServiceImpl) TriggerMitigation(telemetryID int) (models.MitigationResponse, error) {
	if telemetryID < 0 || telemetryID >= len(s.telemetryData) {
		return models.MitigationResponse{}, fmt.Errorf("telemetry ID %d out of bounds", telemetryID)
	}

	point := s.telemetryData[telemetryID]
	if !point.IsAnomaly || point.AnomalyType == nil {
		return models.MitigationResponse{Status: "No threat detected at this point."}, nil
	}

	actions, found := mitigationPlaybook[*point.AnomalyType]
	if !found {
		return models.MitigationResponse{Status: "Unknown threat type, no playbook available."}, nil
	}

	log.Printf("MITIGATION TRIGGERED for threat '%s' at ID %d. Actions: %v", *point.AnomalyType, telemetryID, actions)

	return models.MitigationResponse{
		ThreatType:       *point.AnomalyType,
		SuggestedActions: actions,
		Status:           "Mitigation playbook initiated.",
	}, nil
}

func (s *telemetryServiceImpl) ExportLogsToCSV(writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	header := []string{"timestamp", "telemetry_id", "user_guess_is_anomaly", "was_correct", "response_time_ms"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	for _, log := range s.detectionLogs {
		record := []string{
			log.Timestamp.Format(time.RFC3339),
			strconv.Itoa(log.TelemetryID),
			strconv.FormatBool(log.UserGuess),
			strconv.FormatBool(log.Correct),
			strconv.FormatInt(log.ResponseTimeMs, 10),
		}
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// NOVO MÉTODO: checkAndAwardBadges
func (s *telemetryServiceImpl) checkAndAwardBadges() {
	if s.correctHits == 1 {
		s.playerBadges["First Catch"] = true
	}
	if s.correctHits == 3 {
		s.playerBadges["Threat Hunter"] = true
	}
	// A anomalia do tipo 'sensor_dropout' é a última e mais difícil
	if len(s.detectionLogs) > 0 && s.detectionLogs[len(s.detectionLogs)-1].TelemetryID == 75 {
		s.playerBadges["Critical Failure Averted"] = true
	}
}

// NOVO MÉTODO: GetGameState
func (s *telemetryServiceImpl) GetGameState() models.GameState {
	badges := make([]string, 0, len(s.playerBadges))
	for badge := range s.playerBadges {
		badges = append(badges, badge)
	}
	return models.GameState{
		PlayerScore:  s.playerScore,
		PlayerBadges: badges,
	}
}
