package config

import "time"

type ProducerConfig struct {
	RetryMax     int           `mapstructure:"retry_max"`
	Compression  string        `mapstructure:"compression"`
	Idempotent   bool          `mapstructure:"idempotent"`
	Timeout      time.Duration `mapstructure:"timeout"`
	RequiredAcks int           `mapstructure:"required_acks"`
}
