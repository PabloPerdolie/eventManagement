package handler

import (
	"net/http"
	"net/http/httputil"

	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
)

type ProxyService interface {
	NewCoreServiceProxy() (*httputil.ReverseProxy, error)
	NewNotificationServiceProxy() (*httputil.ReverseProxy, error)
	NewCommunicationServiceProxy() (*httputil.ReverseProxy, error)
}

type Proxy struct {
	proxyService ProxyService
}

func NewProxy(proxyService ProxyService) Proxy {
	return Proxy{
		proxyService: proxyService,
	}
}

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
func (h *Proxy) ProxyToEventService(c *gin.Context) {
	proxy, err := h.proxyService.NewCoreServiceProxy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to proxy request to event service",
		})
		return
	}

	if userId, err := getUserIdFromContext(c); err == nil {
		c.Request.Header.Set("X-User-Id", string(rune(userId)))
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
func (h *Proxy) ProxyToCommunicationService(c *gin.Context) {
	proxy, err := h.proxyService.NewCommunicationServiceProxy()
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to proxy request to communication service",
		})
		return
	}

	if userId, err := getUserIdFromContext(c); err == nil {
		c.Request.Header.Set("X-User-Id", string(rune(userId)))
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
