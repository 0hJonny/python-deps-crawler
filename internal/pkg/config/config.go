package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Logger   LoggerConfig   `mapstructure:"logging"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Database DatabaseConfig `mapstructure:"database"`
	PyPI     PyPIConfig     `mapstructure:"pypi"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Автоматическое связывание с переменными окружения
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Установка значений по умолчанию
	setDefaults()

	// Привязка переменных окружения
	bindEnvironmentVars()

	// Чтение конфигурационного файла (опционально)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Постобработка конфигурации
	if err := postProcessConfig(&config); err != nil {
		return nil, fmt.Errorf("error post-processing config: %w", err)
	}

	// Валидация конфигурации
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	var (
		s ServerConfig
		c CORSConfig
		d DatabaseConfig
		k KafkaConfig
		l LoggerConfig
		r RedisConfig
		p PyPIConfig
	)

	// Init defaults ServerConfig
	s.SetDefaults()

	// Init defaults CORSConfig
	c.SetDefaults()

	// Init defaults DatabaseConfig
	d.SetDefaults()

	// Init defaults KafkaConfig
	k.SetDefaults()

	// Init defaults LoggerConfig
	l.SetDefaults()

	// Init defaults RedisConfig
	r.SetDefaults()

	// Init defaults PyPIConfig
	p.SetDefaults()
}

func bindEnvironmentVars() {
	var (
		s ServerConfig
		d DatabaseConfig
		k KafkaConfig
		l LoggerConfig
		r RedisConfig
		p PyPIConfig
	)

	// Bind ServerConfig vars
	s.BindEnvironmentVars()

	// Bind DatabaseConfig vars
	d.BindEnvironmentVars()

	// Bind KafkaConfig vars
	k.BindEnvironmentVars()

	// Bind LoggerConfig vars
	l.BindEnvironmentVars()

	// Bind RedisConfig vars
	r.BindEnvironmentVars()

	// Bind PyPIConfig vars
	p.BindEnvironmentVars()
}

func postProcessConfig(config *Config) error {
	if err := config.Redis.ParseRedisURL(); err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	config.Kafka.ParseBrokers()

	return nil
}

func validateConfig(config *Config) error {
	if len(config.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka brokers are required")
	}

	if config.Kafka.Topic == "" {
		return fmt.Errorf("kafka topic is required")
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if config.PyPI.APIURL == "" {
		return fmt.Errorf("PyPI API URL is required")
	}

	return nil
}

func (c *Config) LogConfig() {
	fmt.Printf("Configuration loaded:\n")
	fmt.Printf("  Server: %s:%s (mode: %s)\n", c.Server.Host, c.Server.Port, c.Server.Mode)
	fmt.Printf("  Kafka Brokers: %v\n", c.Kafka.Brokers)
	fmt.Printf("  Kafka Topic: %s\n", c.Kafka.Topic)
	fmt.Printf("  Database: %s:%s/%s\n", c.Database.Host, c.Database.Port, c.Database.DBName)
	fmt.Printf("  Redis: %s:%s\n", c.Redis.Host, c.Redis.Port)
	fmt.Printf("  PyPI API: %s\n", c.PyPI.APIURL)
	fmt.Printf("  Log Level: %s (%s)\n", c.Logger.Level, c.Logger.Encoding)
}