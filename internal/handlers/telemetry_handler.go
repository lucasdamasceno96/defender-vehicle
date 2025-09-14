// internal/handlers/telemetry_handler.go
package handlers

import (
	"net/http"
	"strconv"

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
