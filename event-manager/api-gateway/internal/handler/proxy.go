package handler

import (
	"net/http"

	"github.com/event-management/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
)

// ProxyToEventService forwards requests to the Event Service
// @Summary Proxy to Event Service
// @Description Forward requests to the Event Service
// @Tags proxy
// @Accept json
// @Produce json
// @Security Bearer
// @Param path path string true "Request path"
// @Router /events/{path} [get]
// @Router /events/{path} [post]
// @Router /events/{path} [put]
// @Router /events/{path} [delete]
func (h *Handler) ProxyToEventService(c *gin.Context) {
	proxy, err := h.service.Proxy.NewCoreServiceProxy()
	if err != nil {
		h.logger.Errorf("Failed to create event service proxy: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to proxy request to event service",
		})
		return
	}

	// Add user Id to request context if available
	if userId, err := getUserIdFromContext(c); err == nil {
		c.Request.Header.Set("X-User-Id", userId.String())
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

// ProxyToNotificationService forwards requests to the Notification Service
// @Summary Proxy to Notification Service
// @Description Forward requests to the Notification Service
// @Tags proxy
// @Accept json
// @Produce json
// @Security Bearer
// @Param path path string true "Request path"
// @Router /notifications/{path} [get]
// @Router /notifications/{path} [post]
// @Router /notifications/{path} [put]
// @Router /notifications/{path} [delete]
func (h *Handler) ProxyToNotificationService(c *gin.Context) {
	proxy, err := h.service.Proxy.NewNotificationServiceProxy()
	if err != nil {
		h.logger.Errorf("Failed to create notification service proxy: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to proxy request to notification service",
		})
		return
	}

	// Add user Id to request context if available
	if userId, err := getUserIdFromContext(c); err == nil {
		c.Request.Header.Set("X-User-Id", userId.String())
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

// ProxyToCommunicationService forwards requests to the Communication Service
// @Summary Proxy to Communication Service
// @Description Forward requests to the Communication Service
// @Tags proxy
// @Accept json
// @Produce json
// @Security Bearer
// @Param path path string true "Request path"
// @Router /comments/{path} [get]
// @Router /comments/{path} [post]
// @Router /comments/{path} [put]
// @Router /comments/{path} [delete]
func (h *Handler) ProxyToCommunicationService(c *gin.Context) {
	proxy, err := h.service.Proxy.NewCommunicationServiceProxy()
	if err != nil {
		h.logger.Errorf("Failed to create communication service proxy: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to proxy request to communication service",
		})
		return
	}

	// Add user Id to request context if available
	if userId, err := getUserIdFromContext(c); err == nil {
		c.Request.Header.Set("X-User-Id", userId.String())
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
