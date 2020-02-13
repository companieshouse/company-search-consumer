package service

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// "github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/chs.go/log"
)

// messageMetadata represents resource change data json unmarshal
type messageMetadata struct {
	ResourceKind        string `avro:"resource_kind"`
	ResourceURI    		string `avro:"resource_uri"`
	ContextID  		    string `avro:"context_id"`
	ResourceID 		    string `avro:"resource_id"`
	Data         		string `avro:"data"`
	Event               Event  `json:"event"`
}

// Event - resource changed data event
type Event struct {
	FieldsChanged []string `json:"fields_changed,omitempty"`
	Timepoint     int64    `json:"timepoint"`
	PublishedAt   string   `json:"published_at"`
	Type          string   `json:"type"`
}

// Service contains the necessary config for the company-search-consumer service
type Service struct {
	HTTPClient         *http.Client
	NotificationAPIURL string
	Consumer           *consumer.GroupConsumer
	InitialOffset      int64
	Schema             string
}

//Start is called to run the service
func (svc *Service) Start() {
	log.Debug("svc", log.Data{"service": svc})

	// avro := &avro.Schema{
	// 	Definition: svc.Schema,
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case <-c:
			log.Info("Closing consumer")
			err := svc.Consumer.Close()
			if err != nil {
				log.Error(fmt.Errorf("Error closing consumer: %s", err))
			}
			log.Info("Consumer successfully closed")
			return

		case err := <-svc.Consumer.Errors():
			log.Error(err)

		case event := <-svc.Consumer.Messages():
			if svc.InitialOffset == -1 {
				svc.InitialOffset = event.Offset
			}
		}
	}
}
