package config

import "github.com/spf13/viper"

type LoggerConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
	Development      bool     `mapstructure:"development"`
}

func (l *LoggerConfig) SetDefaults() {
	// Zap logger defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")
	viper.SetDefault("logging.output_paths", []string{"stdout"})
	viper.SetDefault("logging.error_output_paths", []string{"stderr"})
	viper.SetDefault("logging.development", false)
}

func (l *LoggerConfig) BindEnvironmentVars() {
	// Logger
	viper.BindEnv("logging.level", "API_GATEWAY_LOGGING_LEVEL")
	viper.BindEnv("logging.encoding", "API_GATEWAY_LOGGING_ENCODING")
	viper.BindEnv("logging.development", "API_GATEWAY_LOGGING_DEVELOPMENT")
}
