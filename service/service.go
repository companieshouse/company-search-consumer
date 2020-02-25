package service

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/companieshouse/chs.go/avro"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/chs.go/log"
)

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
	HTTPClient         *http.Client
	NotificationAPIURL string
	Consumer           *consumer.GroupConsumer
	InitialOffset      int64
	Schema             string
}

//Start is called to run the service
func (svc *Service) Start() {
	log.Debug("svc", log.Data{"service": svc})

	avro := &avro.Schema{
		Definition: svc.Schema,
	}

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
			log.Info("Begin consume")
			if svc.InitialOffset == -1 {
				svc.InitialOffset = event.Offset
			}

			var err error

			if event.Offset >= svc.InitialOffset {

				var mm messageMetadata
				err = avro.Unmarshal(event.Value, &mm)
				if err != nil {
					log.ErrorC("Error unmarshalling avro", err, nil)
					continue
				}

				log.Info("Unmarsheled the avro message")

				// currently we only email - later switch on message type
				// email := &email.Template{
				// 	HTTPClient:         svc.HTTPClient,
				// 	NotificationAPIURL: svc.NotificationAPIURL,
				// }

				// if err != nil {
				// 	log.ErrorC("Error formatting message template data", err, nil)
				// 	continue
				// }

				// log.Event("trace", "", log.Data{"messagetype": mm.MessageType, "appid": mm.AppID, "data": mm.Data})
				// if err = email.SendViaAPI(mm.MessageType, mm.AppID, mm.Data); err != nil {
				// 	log.ErrorC("Error calling Notification API", err, nil)
				// 	continue
				// }

				svc.Consumer.MarkOffset(event, "")
				if err := svc.Consumer.CommitOffsets(); err != nil {
					log.ErrorC("Error commiting message offset", err, nil)
					continue
				}
			}
		}
	}
}
