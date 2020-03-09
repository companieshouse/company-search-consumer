package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/company-search-consumer/upsert"
)

// Value required to unit test the consumer
var testHttpResponseValue int

// messageMetadata represents resource change data json unmarshal
type messageMetadata struct {
	ResourceKind string `avro:"resource_kind"`
	ResourceURI  string `avro:"resource_uri"`
	ContextID    string `avro:"context_id"`
	ResourceID   string `avro:"resource_id"`
	Data         string `avro:"data"`
	Event        Event  `avro:"event"`
}

// Event - resource changed data event
type Event struct {
	FieldsChanged []string `avro:"fields_changed,omitempty"`
	PublishedAt   string   `avro:"published_at"`
	Type          string   `avro:"type"`
}

// Service contains the necessary config for the company-search-consumer service
type Service struct {
	Consumer      *consumer.GroupConsumer
	Marshaller    avro.Marshaller
	Upsert        upsert.Upsert
	InitialOffset int64
	Schema        string
}

//Start is called to run the service
func (svc *Service) Start(c chan os.Signal) {
	log.Debug("svc", log.Data{"service": svc})

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

			var err error

			if event.Offset >= svc.InitialOffset {

				var mm messageMetadata
				err = svc.Marshaller.Unmarshal(event.Value, &mm)
				if err != nil {
					log.ErrorC("Error unmarshalling avro", err, nil)
					continue
				}

				if err != nil {
					log.ErrorC("Error formatting message template data", err, nil)
					continue
				}

				log.Event("trace", "", log.Data{"data": mm.Data})
				statusCode, err := svc.Upsert.SendViaAPI(mm.Data)
				testHttpResponseValue = statusCode
				if err != nil {
					log.ErrorC("Error calling upsert on the search api", err, nil)
					continue
				}

				svc.Consumer.MarkOffset(event, "")
				if err := svc.Consumer.CommitOffsets(); err != nil {
					log.ErrorC("Error committing message offset", err, nil)
					continue
				}
			}
		}
	}
}
