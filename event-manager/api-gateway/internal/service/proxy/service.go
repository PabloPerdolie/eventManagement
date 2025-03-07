package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

// Service handles forwarding requests to appropriate services
type Service struct {
	coreServiceURL          string
	notificationServiceURL  string
	communicationServiceURL string
	logger                  *zap.SugaredLogger
}

// New creates a new proxy service
func New(coreURL, notificationURL, communicationURL string, logger *zap.SugaredLogger) *Service {
	return &Service{
		coreServiceURL:          coreURL,
		notificationServiceURL:  notificationURL,
		communicationServiceURL: communicationURL,
		logger:                  logger,
	}
}

// NewCoreServiceProxy creates a new reverse proxy for the core service
func (s *Service) NewCoreServiceProxy() (*httputil.ReverseProxy, error) {
	coreURL, err := url.Parse(s.coreServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(coreURL)
	s.updateProxyDirector(proxy, coreURL)

	return proxy, nil
}

// NewNotificationServiceProxy creates a new reverse proxy for the notification service
func (s *Service) NewNotificationServiceProxy() (*httputil.ReverseProxy, error) {
	notificationURL, err := url.Parse(s.notificationServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(notificationURL)
	s.updateProxyDirector(proxy, notificationURL)

	return proxy, nil
}

// NewCommunicationServiceProxy creates a new reverse proxy for the communication service
func (s *Service) NewCommunicationServiceProxy() (*httputil.ReverseProxy, error) {
	communicationURL, err := url.Parse(s.communicationServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(communicationURL)
	s.updateProxyDirector(proxy, communicationURL)

	return proxy, nil
}

// updateProxyDirector modifies the proxy director to update request headers and path
func (s *Service) updateProxyDirector(proxy *httputil.ReverseProxy, target *url.URL) {
	originalDirector := proxy.Director
	
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		
		// Update the Host header to the target host
		req.Host = target.Host
		
		// Set X-Forwarded-* headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http") // Use "https" if necessary
		
		// Add any additional headers or modifications
		req.Header.Set("X-Proxy-By", "api-gateway")
		
		s.logger.Debugw("Proxying request",
			"path", req.URL.Path,
			"target", target.String(),
		)
	}
}
