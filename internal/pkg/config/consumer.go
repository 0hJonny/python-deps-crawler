package config

import "time"

type ConsumerConfig struct {
	SessionTimeout   time.Duration `mapstructure:"session_timeout"`
	HeartbeatTimeout time.Duration `mapstructure:"heartbeat_timeout"`
	InitialOffset    string        `mapstructure:"initial_offset"`
	MaxPollRecords   int           `mapstructure:"max_poll_records"`
}
