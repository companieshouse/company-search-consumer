package main

import (
	"github.com/companieshouse/chs.go/kafka/klog"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/company-search-consumer/config"
	"os"
)

func main() {
	cfg := config.Get()

	log.Namespace = "company-search-consumer"
	if len(os.Getenv("HUMAN_LOG")) < 1 {
		klg := &klog.Klog{OutTopic: cfg.LogTopic, BrokerAddr: cfg.StreamingBrokerAddr, SchemaRegistryURL: cfg.SchemaRegistryURL}
		go klg.Capture()
	}
	log.Trace("company-search-consumer starting...", log.Data{"config": cfg})
}
