package config

import (
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/gofigure"
	"os"
)

// Config is the company-search-consumer config
type Config struct {
	gofigure                  interface{} `order:"env,flag"`
	SchemaRegistryURL         string      `env:"SCHEMA_REGISTRY_URL"                 flag:"schema-registry-url"                 flagDesc:"Schema registry url"`
	StreamingBrokerAddr       []string    `env:"KAFKA_STREAMING_BROKER_ADDR"         flag:"streaming-broker-addr"               flagDesc:"Streaming CH Kafka broker cluster address"`
	LogTopic                  string      `env:"LOG_TOPIC"                           flag:"log-topic"                           flagDesc:"Log topic"`
	InitialNotifyOffset       int64       `env:"INITIAL_OFFSET"                      flag:"initial-offset"                      flagDesc:"Initial offset for consumer group"`
	StreamCompanyProfileTopic string      `env:"STREAM_COMPANY_PROFILE_TOPIC"        flag:"stream-company-profile-topic"        flagDesc:"Stream Company Profile Topic"`
	ZookeeperURL              string      `env:"KAFKA_ZOOKEEPER_ADDR"                flag:"zookeeper-addr"                      flagDesc:"Zookeeper address"`
	ConsumerGroupName         string      `env:"CONSUMER_GROUP_NAME"                 flag:"consumer-group-name"                 flagDesc:"Consumer group name"`
	ZookeeperChroot           string      `env:"KAFKA_ZOOKEEPER_CHROOT"              flag:"zookeeper-chroot"                    flagDesc:"Zookeeper chroot"`

}

var cfg *Config

// Get configures the application and returns the configuration
func Get() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	return cfg
}