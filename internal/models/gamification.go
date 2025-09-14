// internal/models/gamification.go
package models

import "time"

// DetectionRequest is the payload sent by the user to flag an anomaly.
type DetectionRequest struct {
	TelemetryID int `json:"telemetry_id" binding:"required"`
}

// DetectionLog stores the result of a user's action.
type DetectionLog struct {
	Timestamp      time.Time `json:"ts"`
	TelemetryID    int       `json:"telemetry_id"`
	UserGuess      bool      `json:"user_guess"` // true for anomaly, false for normal
	Correct        bool      `json:"correct"`
	ResponseTimeMs int64     `json:"response_time_ms"`
}

// MitigationResponse represents the outcome of a mitigation action.
type MitigationResponse struct {
	ThreatType       string   `json:"threat_type"`
	SuggestedActions []string `json:"suggested_actions"`
	Status           string   `json:"status"`
}

// GameState represents the player's current progress.
type GameState struct {
	PlayerScore  int      `json:"player_score"`
	PlayerBadges []string `json:"player_badges"`
}
