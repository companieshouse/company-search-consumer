// Package config defines the environment variable and command-line flags
package config

import (
	"os"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/gofigure"
)

// Config is the company-search-consumer config
type Config struct {
	gofigure                  interface{} `order:"env,flag"`
	ConsumerGroupName         string      `env:"CONSUMER_GROUP_NAME"                 flag:"consumer-group-name"                 flagDesc:"Consumer group name"`
	StreamCompanyProfileTopic string      `env:"STREAM_COMPANY_PROFILE_TOPIC"        flag:"stream-company-profile-topic"        flagDesc:"Stream Company Profile Topic"`
	SchemaRegistryURL         string      `env:"SCHEMA_REGISTRY_URL"                 flag:"schema-registry-url"                 flagDesc:"Schema registry url"`
	StreamingBrokerAddr       []string    `env:"KAFKA_STREAMING_BROKER_ADDR"         flag:"streaming-broker-addr"               flagDesc:"Streaming CH Kafka broker cluster address"`
	StreamingZookeeperURL     string      `env:"KAFKA_ZOOKEEPER_STREAMING_ADDR"      flag:"zookeeper-addr"                      flagDesc:"Main CH Zookeeper address"`
	ZookeeperChroot           string      `env:"KAFKA_ZOOKEEPER_CHROOT"              flag:"zookeeper-chroot"                    flagDesc:"Zookeeper chroot"`
	InitialOffset             int64       `env:"INITIAL_OFFSET"                      flag:"initial-offset"                      flagDesc:"Initial offset for consumer group"`
	LogTopic                  string      `env:"LOG_TOPIC"                           flag:"log-topic"                           flagDesc:"Log topic"`
	UpsertURL                 string      `env:"UPSERT_URL"                          flag:"upsert-url"                          flagDesc:"Url for upserting a company"`
	RetryThrottleRate         int         `env:"RETRY_THROTTLE_RATE_SECONDS"         flag:"retry-throttle-rate-seconds"         flagDesc:"Retry throttle rate seconds"`
	MaxRetryAttempts          int         `env:"MAXIMUM_RETRY_ATTEMPTS"              flag:"max-retry-attempts"                  flagDesc:"Maximum retry attempts"`
	CHSAPIKey                 string      `env:"CHS_INTERNAL_API_KEY"                flag:"chs-api-key"                         flagDesc:"API key to call api's through eric"`
}

var cfg *Config

// Get configures the application and returns the configuration
func Get() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = &Config{
		ConsumerGroupName:         "stream-company-profile-consumer-group",
		StreamCompanyProfileTopic: "stream-company-profile",
		ZookeeperChroot:           "",
		InitialOffset:             int64(-1),
	}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	return cfg
}
