package config

import (
	"fmt"
	"strings"
	"time"

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

// TODO: Split the config

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

type KafkaConfig struct {
	Brokers       []string       `mapstructure:"brokers"`
	Topic         string         `mapstructure:"topic"`
	ConsumerGroup string         `mapstructure:"consumer_group"`
	Producer      ProducerConfig `mapstructure:"producer"`
	Consumer      ConsumerConfig `mapstructure:"consumer"`
}

type ProducerConfig struct {
	RetryMax     int           `mapstructure:"retry_max"`
	Compression  string        `mapstructure:"compression"`
	Idempotent   bool          `mapstructure:"idempotent"`
	Timeout      time.Duration `mapstructure:"timeout"`
	RequiredAcks int           `mapstructure:"required_acks"`
}

type ConsumerConfig struct {
	SessionTimeout   time.Duration `mapstructure:"session_timeout"`
	HeartbeatTimeout time.Duration `mapstructure:"heartbeat_timeout"`
	InitialOffset    string        `mapstructure:"initial_offset"`
	MaxPollRecords   int           `mapstructure:"max_poll_records"`
}

type LoggerConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
	Development      bool     `mapstructure:"development"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

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

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"db_name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type PyPIConfig struct {
	APIURL         string        `mapstructure:"api_url"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	MaxRetries     int           `mapstructure:"max_retries"`
	RateLimit      int           `mapstructure:"rate_limit"`
	CacheTimeout   time.Duration `mapstructure:"cache_timeout"`
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
	// Server defaults with graceful shutdown
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.max_request_size", 4*1024*1024)
	viper.SetDefault("server.enable_tls", false)

	// Base Kafka architecture defaults
	viper.SetDefault("kafka.topic", "dependency.analysis.request")
	viper.SetDefault("kafka.consumer_group", "api-gateway-consumer")
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})

	// Kafka Producer defaults (simplified for base architecture)
	viper.SetDefault("kafka.producer.retry_max", 3)
	viper.SetDefault("kafka.producer.compression", "lz4")
	viper.SetDefault("kafka.producer.idempotent", true)
	viper.SetDefault("kafka.producer.timeout", "10s")
	viper.SetDefault("kafka.producer.required_acks", 1)

	// Kafka Consumer defaults
	viper.SetDefault("kafka.consumer.session_timeout", "10s")
	viper.SetDefault("kafka.consumer.heartbeat_timeout", "3s")
	viper.SetDefault("kafka.consumer.initial_offset", "latest")
	viper.SetDefault("kafka.consumer.max_poll_records", 500)

	// Zap logger defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")
	viper.SetDefault("logging.output_paths", []string{"stdout"})
	viper.SetDefault("logging.error_output_paths", []string{"stderr"})
	viper.SetDefault("logging.development", false)

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"})

	// Redis defaults
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")
	viper.SetDefault("redis.pool_size", 10)

	// Database defaults
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// PyPI defaults
	viper.SetDefault("pypi.api_url", "https://pypi.org/pypi")
	viper.SetDefault("pypi.request_timeout", "30s")
	viper.SetDefault("pypi.max_retries", 3)
	viper.SetDefault("pypi.rate_limit", 100)
	viper.SetDefault("pypi.cache_timeout", "1h")
}

func bindEnvironmentVars() {
	// Server
	viper.BindEnv("server.port", "API_GATEWAY_SERVER_PORT")
	viper.BindEnv("server.host", "API_GATEWAY_SERVER_HOST")
	viper.BindEnv("server.mode", "API_GATEWAY_SERVER_MODE")

	// Kafka
	viper.BindEnv("kafka.brokers", "API_GATEWAY_KAFKA_BROKERS")
	viper.BindEnv("kafka.topic", "API_GATEWAY_KAFKA_TOPIC")
	viper.BindEnv("kafka.consumer_group", "API_GATEWAY_KAFKA_CONSUMER_GROUP")

	// Logger
	viper.BindEnv("logging.level", "API_GATEWAY_LOGGING_LEVEL")
	viper.BindEnv("logging.encoding", "API_GATEWAY_LOGGING_ENCODING")
	viper.BindEnv("logging.development", "API_GATEWAY_LOGGING_DEVELOPMENT")

	// Database
	viper.BindEnv("database.host", "POSTGRES_HOST")
	viper.BindEnv("database.db_name", "POSTGRES_DB")
	viper.BindEnv("database.user", "POSTGRES_USER")
	viper.BindEnv("database.password", "POSTGRES_PASSWORD")
	viper.BindEnv("database.port", "POSTGRES_PORT")

	// Redis
	viper.BindEnv("redis.url", "REDIS_URL")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")

	// PyPI
	viper.BindEnv("pypi.api_url", "PYPI_API_URL")
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

func (c *KafkaConfig) ParseBrokers() {
	if len(c.Brokers) == 1 {
		brokersList := strings.Split(c.Brokers[0], ",")
		var cleanBrokers []string
		for _, broker := range brokersList {
			cleanBroker := strings.TrimSpace(broker)
			if cleanBroker != "" {
				cleanBrokers = append(cleanBrokers, cleanBroker)
			}
		}
		c.Brokers = cleanBrokers
	}
}

func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *RedisConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
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

func (s ServerConfig) GetConfig() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
