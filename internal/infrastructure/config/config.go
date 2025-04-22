package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	APIGateway APIGatewayConfig
}

type ServerConfig struct {
	Port string
	Mode string // "debug" or "release"
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
}

type APIGatewayConfig struct {
	Enabled      bool
	Port         string
	Services     map[string]ServiceConfig
	RateLimit    RateLimitConfig
	CircuitBreak CircuitBreakConfig
	Tracing      bool
	Logging      bool
}

type ServiceConfig struct {
	URL         string
	Timeout     int
	RetryCount  int
	HealthCheck string
}

type RateLimitConfig struct {
	RequestsPerSecond int
	Burst             int
}

type CircuitBreakConfig struct {
	MaxRequests      uint32
	Interval         int
	Timeout          int
	ConsecutiveError int
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "erp_db")

	viper.SetDefault("jwt.access_secret", "your-access-secret-key")
	viper.SetDefault("jwt.refresh_secret", "your-refresh-secret-key")

	// API Gateway defaults
	viper.SetDefault("apigateway.enabled", true)
	viper.SetDefault("apigateway.port", "8000")
	viper.SetDefault("apigateway.tracing", true)
	viper.SetDefault("apigateway.logging", true)

	// Rate limiting defaults
	viper.SetDefault("apigateway.ratelimit.requests_per_second", 100)
	viper.SetDefault("apigateway.ratelimit.burst", 50)

	// Circuit breaking defaults
	viper.SetDefault("apigateway.circuitbreak.max_requests", 100)
	viper.SetDefault("apigateway.circuitbreak.interval", 60)
	viper.SetDefault("apigateway.circuitbreak.timeout", 30)
	viper.SetDefault("apigateway.circuitbreak.consecutive_error", 5)

	// Default services
	viper.SetDefault("apigateway.services.auth.url", "http://localhost:8080")
	viper.SetDefault("apigateway.services.auth.timeout", 30)
	viper.SetDefault("apigateway.services.auth.retry_count", 3)
	viper.SetDefault("apigateway.services.auth.health_check", "/health")

	viper.SetDefault("apigateway.services.user.url", "http://localhost:8080")
	viper.SetDefault("apigateway.services.user.timeout", 30)
	viper.SetDefault("apigateway.services.user.retry_count", 3)
	viper.SetDefault("apigateway.services.user.health_check", "/health")

	viper.SetDefault("apigateway.services.audit.url", "http://localhost:8080")
	viper.SetDefault("apigateway.services.audit.timeout", 30)
	viper.SetDefault("apigateway.services.audit.retry_count", 3)
	viper.SetDefault("apigateway.services.audit.health_check", "/health")

	viper.SetDefault("apigateway.services.stock.url", "http://localhost:8080")
	viper.SetDefault("apigateway.services.stock.timeout", 30)
	viper.SetDefault("apigateway.services.stock.retry_count", 3)
	viper.SetDefault("apigateway.services.stock.health_check", "/health")

	viper.SetDefault("apigateway.services.vendor.url", "http://localhost:8080")
	viper.SetDefault("apigateway.services.vendor.timeout", 30)
	viper.SetDefault("apigateway.services.vendor.retry_count", 3)
	viper.SetDefault("apigateway.services.vendor.health_check", "/health")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ERP")

	// Load services from environment
	services := make(map[string]ServiceConfig)
	for _, svc := range []string{"audit", "user", "auth", "stock", "sku", "vendor", "manufacturing", "purchase", "order", "client", "finance", "report"} {
		services[svc] = ServiceConfig{
			URL:         viper.GetString(fmt.Sprintf("apigateway.services.%s.url", svc)),
			Timeout:     viper.GetInt(fmt.Sprintf("apigateway.services.%s.timeout", svc)),
			RetryCount:  viper.GetInt(fmt.Sprintf("apigateway.services.%s.retry_count", svc)),
			HealthCheck: viper.GetString(fmt.Sprintf("apigateway.services.%s.health_check", svc)),
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
			Mode: viper.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
		},
		JWT: JWTConfig{
			AccessSecret:  viper.GetString("jwt.access_secret"),
			RefreshSecret: viper.GetString("jwt.refresh_secret"),
		},
		APIGateway: APIGatewayConfig{
			Enabled:  viper.GetBool("apigateway.enabled"),
			Port:     viper.GetString("apigateway.port"),
			Services: services,
			RateLimit: RateLimitConfig{
				RequestsPerSecond: viper.GetInt("apigateway.ratelimit.requests_per_second"),
				Burst:             viper.GetInt("apigateway.ratelimit.burst"),
			},
			CircuitBreak: CircuitBreakConfig{
				MaxRequests:      viper.GetUint32("apigateway.circuitbreak.max_requests"),
				Interval:         viper.GetInt("apigateway.circuitbreak.interval"),
				Timeout:          viper.GetInt("apigateway.circuitbreak.timeout"),
				ConsecutiveError: viper.GetInt("apigateway.circuitbreak.consecutive_error"),
			},
			Tracing: viper.GetBool("apigateway.tracing"),
			Logging: viper.GetBool("apigateway.logging"),
		},
	}

	return cfg, nil
}
