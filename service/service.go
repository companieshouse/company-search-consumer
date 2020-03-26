// Package service contains the creation and start of the listener for the Kafka topic
package service

import (
	"fmt"

	event "github.com/companieshouse/chs-go-avro-schemas/common"
	"github.com/companieshouse/chs.go/avro"
	consumer "github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/company-search-consumer/upsert"

	"os"
	"os/signal"
	"syscall"
)

// messageMetadata represents resource change data json unmarshal
type messageMetadata struct {
	ResourceKind string            `avro:"resource_kind"`
	ResourceURI  string            `avro:"resource_uri"`
	ContextID    string            `avro:"context_id"`
	ResourceID   string            `avro:"resource_id"`
	Data         string            `avro:"data"`
	Event        event.EventRecord `avro:"event"`
}

// Service contains the necessary config for the company-search-consumer service
type Service struct {
	Consumer      *consumer.GroupConsumer
	Marshaller    avro.Marshaller
	Upsert        upsert.Upsert
	HandleError   func(err error, offset int64, str interface{}) error
	InitialOffset int64
	Schema        string
}

// Start is called to run the service
func (svc *Service) Start(c chan os.Signal) {
	log.Debug("svc", log.Data{"service": svc})

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case <-c:
			log.Info("Closing consumer")
			err := svc.Consumer.Close()
			if err != nil {
				log.Error(fmt.Errorf("error closing consumer: %s", err))
			}
			log.Info("Consumer successfully closed")
			return

		case err := <-svc.Consumer.Errors():
			log.Error(err)

		case message := <-svc.Consumer.Messages():

			if svc.InitialOffset == -1 {
				svc.InitialOffset = message.Offset
			}

			if message.Offset >= svc.InitialOffset {

				var mm messageMetadata
				if err := svc.Marshaller.Unmarshal(message.Value, &mm); err != nil {
					log.ErrorC("Error unmarshalling avro", err, nil)
					if err = svc.HandleError(err, message.Offset, &message.Value); err != nil {
						log.ErrorC("Error handling unmarshalling error, kafka message will not be retried", err, nil)
					}
					continue
				}

				if err := svc.Upsert.SendViaAPI(mm.Data); err != nil {
					log.ErrorC("Error calling upsert on the search api", err, nil)
					if err = svc.HandleError(err, message.Offset, &mm); err != nil {
						log.ErrorC("Error handling upserting error, kafka message will not be retried", err, nil)
					}
					continue
				}

				log.Trace("Committing message", log.Data{"offset": message.Offset})
				svc.Consumer.MarkOffset(message, "")
				if err := svc.Consumer.CommitOffsets(); err != nil {
					log.ErrorC("Error committing message offset", err, log.Data{"offset": message.Offset})
					if err = svc.HandleError(err, message.Offset, &mm); err != nil {
						log.ErrorC("Error handling committing error, kafka message will not be retired", err, nil)
					}
					continue
				}
			}
		}
	}
}
