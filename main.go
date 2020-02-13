package main

import (
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/company-search-consumer/config"
	"github.com/companieshouse/company-search-consumer/service"
	gologger "log"
	"net/http"
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

	log.Debug("consumer has been successfully initialised")

	svc := &service.Service{
		Schema:             "resource-changed-data",
		NotificationAPIURL: "http://api.chs-dev.internal:4089/upsert",
		HTTPClient:         http.DefaultClient,
		Consumer:           consumer,
		InitialOffset:      cfg.InitialOffset,
	}

	svc.Start()
}
