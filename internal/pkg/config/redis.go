package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	URL          string        `mapstructure:"url"`
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	PoolSize     int           `mapstructure:"pool_size"`
}

func (r *RedisConfig) SetDefaults() {
	// Redis defaults
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")
	viper.SetDefault("redis.pool_size", 10)
}

func (r *RedisConfig) BindEnvironmentVars() {
	// Redis
	viper.BindEnv("redis.url", "REDIS_URL")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
}

// Existing helper methods remain the same
func (c *RedisConfig) ParseRedisURL() error {
	if c.URL != "" {
		urlToParse := c.URL
		if after, ok := strings.CutPrefix(urlToParse, "redis://"); ok {
			urlToParse = after
		}

		if strings.Contains(urlToParse, "@") {
			parts := strings.Split(urlToParse, "@")
			if len(parts) == 2 {
				authPart := parts[0]
				if strings.Contains(authPart, ":") {
					authParts := strings.Split(authPart, ":")
					if len(authParts) == 2 {
						c.Password = authParts[1]
					}
				}
				urlToParse = parts[1]
			}
		}

		if strings.Contains(urlToParse, ":") {
			parts := strings.Split(urlToParse, ":")
			if len(parts) == 2 {
				c.Host = parts[0]
				c.Port = parts[1]
			}
		} else {
			c.Host = urlToParse
			c.Port = "6379"
		}
	}
	return nil
}

func (c *RedisConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
