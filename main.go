package main

import (
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/avro/schema"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/chs.go/kafka/producer"
	"github.com/companieshouse/chs.go/kafka/resilience"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/company-search-consumer/config"
	"github.com/companieshouse/company-search-consumer/service"
	"github.com/companieshouse/company-search-consumer/upsert"
	gologger "log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.Get()

	log.Namespace = "company-search-consumer"
	log.Trace("company-search-consumer starting...", log.Data{"config": cfg})

	// Push the Sarama logs into our custom writer
	sarama.Logger = gologger.New(&log.Writer{}, "[Sarama] ", gologger.LstdFlags)

	var resetOffset bool
	if cfg.InitialOffset != -1 {
		resetOffset = true
	}

	resourceChangedDataSchema, err := schema.Get(cfg.SchemaRegistryURL, "resource-changed-data")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	p, err := producer.New(&producer.Config{Acks: &producer.WaitForAll, BrokerAddrs: cfg.StreamingBrokerAddr})

	retry := &resilience.ServiceRetry{
		time.Duration(cfg.RetryThrottleRate),
		cfg.MaxRetryAttempts,
	}

	ppSchema, err := schema.Get(cfg.SchemaRegistryURL, resourceChangedDataSchema)

	rh := resilience.NewHandler(cfg.StreamCompanyProfileTopic, "consumer", retry, p, &avro.Schema{Definition: ppSchema})

	consumerConfig := &consumer.Config{
		Topics:       []string{cfg.StreamCompanyProfileTopic},
		ZookeeperURL: cfg.StreamingZookeeperURL,
		BrokerAddr:   cfg.StreamingBrokerAddr,
	}
	groupConfig := &consumer.GroupConfig{
		GroupName:   cfg.ConsumerGroupName,
		ResetOffset: resetOffset,
		Chroot:      cfg.ZookeeperChroot,
	}

	log.Debug("Joining consumer group", log.Data{"group": cfg.ConsumerGroupName})

	consumer := consumer.NewConsumerGroup(consumerConfig)
	if err := consumer.JoinGroup(groupConfig); err != nil {
		log.ErrorC("Error joining consumer group", err, nil)
		return
	}

	log.Debug("Consumer has been successfully initialised")

	avro := &avro.Schema{
		Definition: resourceChangedDataSchema,
	}

	upsert := &upsert.APIUpsert{
		HTTPClient:          http.DefaultClient,
		UpsertCompanyAPIUrl: cfg.UpsertURL,
	}

	svc := &service.Service{
		Schema:        resourceChangedDataSchema,
		Consumer:      consumer,
		Marshaller:    avro,
		HandleError:   rh.HandleError,
		Upsert:        upsert,
		InitialOffset: cfg.InitialOffset,
	}

	mainChannel := make(chan os.Signal, 1)
	svc.Start(mainChannel)
}
