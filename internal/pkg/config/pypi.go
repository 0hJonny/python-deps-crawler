package config

import (
	"time"

	"github.com/spf13/viper"
)

type PyPIConfig struct {
	APIURL         string        `mapstructure:"api_url"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	MaxRetries     int           `mapstructure:"max_retries"`
	RateLimit      int           `mapstructure:"rate_limit"`
	CacheTimeout   time.Duration `mapstructure:"cache_timeout"`
}

func (p *PyPIConfig) SetDefaults() {
	// PyPI defaults
	viper.SetDefault("pypi.api_url", "https://pypi.org/pypi")
	viper.SetDefault("pypi.request_timeout", "30s")
	viper.SetDefault("pypi.max_retries", 3)
	viper.SetDefault("pypi.rate_limit", 100)
	viper.SetDefault("pypi.cache_timeout", "1h")
}

func (p *PyPIConfig) BindEnvironmentVars() {
	// PyPI
	viper.BindEnv("pypi.api_url", "PYPI_API_URL")
}
