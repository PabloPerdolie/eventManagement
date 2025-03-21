package handler

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// Health represents the health status of the service
type Health struct {
    Status    string `json:"status"`
    Timestamp string `json:"timestamp"`
    Version   string `json:"version"`
    Service   string `json:"service"`
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
    health := Health{
        Status:    "UP",
        Timestamp: time.Now().Format(time.RFC3339),
        Version:   "1.0.0",
        Service:   "notification-service",
    }

    c.JSON(http.StatusOK, health)
}

// ServiceInfo returns information about the notification service
func (h *Handler) ServiceInfo(c *gin.Context) {
    stats := h.service.Notification.GetStats()
    stats["service"] = "notification-service"
    stats["version"] = "1.0.0"
    stats["timestamp"] = time.Now().Format(time.RFC3339)

    c.JSON(http.StatusOK, stats)
}
