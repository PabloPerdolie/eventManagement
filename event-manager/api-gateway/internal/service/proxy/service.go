package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	coreServiceURL          string
	notificationServiceURL  string
	communicationServiceURL string
	logger                  *zap.SugaredLogger
}

func New(coreURL, notificationURL, communicationURL string, logger *zap.SugaredLogger) *Service {
	return &Service{
		coreServiceURL:          coreURL,
		notificationServiceURL:  notificationURL,
		communicationServiceURL: communicationURL,
		logger:                  logger,
	}
}

func (s *Service) NewCoreServiceProxy() (*httputil.ReverseProxy, error) {
	coreURL, err := url.Parse(s.coreServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(coreURL)
	s.updateProxyDirector(proxy, coreURL)
	s.setupProxyErrorHandler(proxy, "core-service")

	return proxy, nil
}

func (s *Service) NewNotificationServiceProxy() (*httputil.ReverseProxy, error) {
	notificationURL, err := url.Parse(s.notificationServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(notificationURL)
	s.updateProxyDirector(proxy, notificationURL)
	s.setupProxyErrorHandler(proxy, "notification-service")

	return proxy, nil
}

func (s *Service) NewCommunicationServiceProxy() (*httputil.ReverseProxy, error) {
	communicationURL, err := url.Parse(s.communicationServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(communicationURL)
	s.updateProxyDirector(proxy, communicationURL)
	s.setupProxyErrorHandler(proxy, "communication-service")

	return proxy, nil
}

func (s *Service) updateProxyDirector(proxy *httputil.ReverseProxy, target *url.URL) {
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		start := time.Now()
		originalDirector(req)

		req.Host = target.Host

		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-Proto", "http")

		req.Header.Set("X-Proxy-By", "api-gateway")

		s.logger.Infow("Proxying request",
			"method", req.Method,
			"path", req.URL.Path,
			"target_host", target.Host,
			"target_path", req.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}

func (s *Service) setupProxyErrorHandler(proxy *httputil.ReverseProxy, serviceName string) {
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		s.logger.Errorw("Proxy error",
			"service", serviceName,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		)

		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("Content-Type", "application/json")

		errorResponse := `{"error":"Service Unavailable","message":"The requested service is temporarily unavailable. Please try again later."}`
		w.Write([]byte(errorResponse))
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		s.logger.Infow("Proxy response received",
			"service", serviceName,
			"method", resp.Request.Method,
			"path", resp.Request.URL.Path,
			"status", resp.StatusCode,
		)
		return nil
	}
}
