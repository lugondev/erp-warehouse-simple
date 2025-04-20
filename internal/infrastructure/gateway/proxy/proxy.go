package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/config"
)

// ServiceProxy handles proxying requests to backend services
type ServiceProxy struct {
	services map[string]config.ServiceConfig
	client   *http.Client
}

// NewServiceProxy creates a new service proxy
func NewServiceProxy(services map[string]config.ServiceConfig) *ServiceProxy {
	return &ServiceProxy{
		services: services,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     60 * time.Second,
			},
		},
	}
}

// ProxyRequest returns a handler that proxies requests to a backend service
func (p *ServiceProxy) ProxyRequest(serviceName, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, exists := p.services[serviceName]
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Service %s not configured", serviceName),
			})
			return
		}

		// Build the target URL
		targetPath := path
		for _, param := range c.Params {
			targetPath = strings.Replace(targetPath, ":"+param.Key, param.Value, -1)
		}

		// Add query parameters
		targetQuery := c.Request.URL.RawQuery
		var targetURL string
		if targetQuery != "" {
			targetURL = fmt.Sprintf("%s%s?%s", service.URL, targetPath, targetQuery)
		} else {
			targetURL = fmt.Sprintf("%s%s", service.URL, targetPath)
		}

		// Create the request
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(reqBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to create request: %v", err),
			})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Add X-Forwarded headers
		req.Header.Set("X-Forwarded-For", c.ClientIP())
		req.Header.Set("X-Forwarded-Proto", c.Request.Proto)
		req.Header.Set("X-Forwarded-Host", c.Request.Host)

		// Add tracing headers if available
		if requestID, exists := c.Get("request_id"); exists {
			req.Header.Set("X-Request-ID", requestID.(string))
		}

		// Add user context if available
		if userID, exists := c.Get("user_id"); exists {
			req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
		}
		if username, exists := c.Get("username"); exists {
			req.Header.Set("X-Username", username.(string))
		}
		if role, exists := c.Get("role"); exists {
			req.Header.Set("X-Role", role.(string))
		}

		// Execute the request
		resp, err := p.client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("Service %s unavailable: %v", serviceName, err),
			})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Copy response status code
		c.Writer.WriteHeader(resp.StatusCode)

		// Copy response body
		io.Copy(c.Writer, resp.Body)
	}
}

// LoadBalancer represents a simple load balancer
type LoadBalancer struct {
	targets []string
	current int
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(targets []string) *LoadBalancer {
	return &LoadBalancer{
		targets: targets,
		current: 0,
	}
}

// NextTarget returns the next target in a round-robin fashion
func (lb *LoadBalancer) NextTarget() string {
	if len(lb.targets) == 0 {
		return ""
	}

	target := lb.targets[lb.current]
	lb.current = (lb.current + 1) % len(lb.targets)
	return target
}

// ServiceDiscovery represents a service discovery mechanism
type ServiceDiscovery struct {
	services map[string][]string
}

// NewServiceDiscovery creates a new service discovery
func NewServiceDiscovery() *ServiceDiscovery {
	return &ServiceDiscovery{
		services: make(map[string][]string),
	}
}

// Register registers a service instance
func (sd *ServiceDiscovery) Register(name, url string) {
	if _, exists := sd.services[name]; !exists {
		sd.services[name] = []string{}
	}
	sd.services[name] = append(sd.services[name], url)
}

// Unregister unregisters a service instance
func (sd *ServiceDiscovery) Unregister(name, url string) {
	if instances, exists := sd.services[name]; exists {
		for i, instance := range instances {
			if instance == url {
				sd.services[name] = append(instances[:i], instances[i+1:]...)
				break
			}
		}
	}
}

// GetInstances returns all instances of a service
func (sd *ServiceDiscovery) GetInstances(name string) []string {
	if instances, exists := sd.services[name]; exists {
		return instances
	}
	return []string{}
}

// GetInstance returns a single instance of a service
func (sd *ServiceDiscovery) GetInstance(name string) (string, error) {
	instances := sd.GetInstances(name)
	if len(instances) == 0 {
		return "", fmt.Errorf("no instances available for service %s", name)
	}
	return instances[0], nil
}
