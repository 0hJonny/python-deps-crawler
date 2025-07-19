package config

import (
	"strings"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers       []string       `mapstructure:"brokers"`
	Topic         string         `mapstructure:"topic"`
	ConsumerGroup string         `mapstructure:"consumer_group"`
	Producer      ProducerConfig `mapstructure:"producer"`
	Consumer      ConsumerConfig `mapstructure:"consumer"`
}

func (k *KafkaConfig) SetDefaults() {
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
}

func (k *KafkaConfig) BindEnvironmentVars() {
	// Kafka
	viper.BindEnv("kafka.brokers", "API_GATEWAY_KAFKA_BROKERS")
	viper.BindEnv("kafka.topic", "API_GATEWAY_KAFKA_TOPIC")
	viper.BindEnv("kafka.consumer_group", "API_GATEWAY_KAFKA_CONSUMER_GROUP")
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
