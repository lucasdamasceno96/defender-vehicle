// internal/handlers/telemetry_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/lucasdamasceno96/defender-vehicle/internal/models"
	"github.com/lucasdamasceno96/defender-vehicle/internal/services"

	"github.com/gin-gonic/gin"
)

// TelemetryHandler holds the dependencies for the telemetry handlers.
type TelemetryHandler struct {
	service services.TelemetryService
}

// NewTelemetryHandler creates a new TelemetryHandler with the given service.
func NewTelemetryHandler(service services.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{service: service}
}

// GetTelemetry is the handler for the /api/telemetry endpoint.
func (h *TelemetryHandler) GetTelemetry(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	data := h.service.GetTelemetry(page, limit)
	c.JSON(http.StatusOK, data)
}

func (h *TelemetryHandler) DetectAnomaly(c *gin.Context) {
	var req models.DetectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	correct, err := h.service.ProcessDetection(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if correct {
		c.JSON(http.StatusOK, gin.H{"status": "Correct! Anomaly detected.", "result": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Incorrect. This point is not an anomaly.", "result": false})
	}
}

func (h *TelemetryHandler) TriggerMitigation(c *gin.Context) {
	// Let's get the telemetry_id from the URL, e.g., /api/mitigate/25
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telemetry ID"})
		return
	}

	response, err := h.service.TriggerMitigation(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
