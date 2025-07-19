package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	MaxRequestSize  int64         `mapstructure:"max_request_size"`
	EnableTLS       bool          `mapstructure:"enable_tls"`
	CertFile        string        `mapstructure:"cert_file"`
	KeyFile         string        `mapstructure:"key_file"`
}

func (s *ServerConfig) SetDefaults() {
	// Server defaults with graceful shutdown
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.max_request_size", 4*1024*1024)
	viper.SetDefault("server.enable_tls", false)
}

func (s *ServerConfig) BindEnvironmentVars() {
	// Server
	viper.BindEnv("server.port", "API_GATEWAY_SERVER_PORT")
	viper.BindEnv("server.host", "API_GATEWAY_SERVER_HOST")
	viper.BindEnv("server.mode", "API_GATEWAY_SERVER_MODE")
}

func (s ServerConfig) GetConfig() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
