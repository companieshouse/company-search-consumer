package service

import (
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"os"
	"testing"
)

const upsertCompanyAPIUrl = "upsertCompanyAPIUrl"
const resourceKind = "resourceKind"
const resourceUri = "resourceUri"
const contextID = "contextID"
const resourceID = "resourceID"

func createMockService() *Service {

	return &Service{
		HTTPClient:          &http.Client{},
		UpsertCompanyAPIUrl: upsertCompanyAPIUrl,
		InitialOffset:       int64(-1),
		Schema:              getDefaultSchema(),
	}
}

func createMockConsumerWithMessage() *consumer.GroupConsumer {

	return &consumer.GroupConsumer{
		GConsumer: MockConsumer{},
	}
}

type MockConsumer struct{}

func (m MockConsumer) Close() error {
	return nil
}

func (m MockConsumer) Messages() <-chan *sarama.ConsumerMessage {
	out := make(chan *sarama.ConsumerMessage)

	go func() {
		out <- &sarama.ConsumerMessage{
			Value: []byte("\"" + resourceKind + "\"" + resourceUri + "\"" + contextID + "\"" + resourceID + "\""),
		}
		close(out)
	}()

	return out
}

func (m MockConsumer) Errors() <-chan error {
	return nil
}

func getDefaultSchema() string {

	return "{ \"type\": \"record\", \"name\": \"resource_changed_data\", \"namespace\": \"stream\", \"fields\": [ { \"name\": \"resource_kind\", \"type\": \"string\" }, { \"name\": \"resource_uri\", \"type\": \"string\" }, { \"name\": \"context_id\", \"type\": \"string\" }, { \"name\": \"resource_id\", \"type\": \"string\" }, { \"name\": \"data\", \"type\": \"string\" }, { \"name\": \"event\", \"type\": { \"type\": \"record\", \"name\": \"event_record\", \"fields\": [ { \"name\": \"published_at\", \"type\": \"string\", \"default\":\"\" }, { \"name\": \"type\", \"type\": \"string\", \"default\":\"\" }, { \"name\": \"fields_changed\", \"type\": [ \"null\", { \"type\": \"array\", \"items\": \"string\" } ] } ] } } ] }"
}

type MockAvro struct{}

func (m MockAvro) Unmarshal() error {
	return nil
}

func TestUnitStart(t *testing.T) {
	Convey("Successful process of a single Kafka message", t, func() {

		c := make(chan os.Signal)
		svc := createMockService()

		Convey("Given a message is readily available for the service to consume", func() {
			svc.Consumer = createMockConsumerWithMessage()

			// Convey("When the avro message is unmarshalled", func() {
			// when(mockAvro.Unmarshal(event.Value, &mm)).thenReturn(nil);

			  // Convey("And the data is sent to search api", func() {
			  // when(upsert.SendViaAPI(mm.data).thenReturn(nil)

				// Convey("Then a successful response is returned", func() {
				// assert(Http Response is successful)

					// Consumer needs to be shut down otherwise will run forever
						//Convey("And the consumer is shutdown", func() {

						// Not sure what the condition should be
						// if (err == nil) {
						//   endConsumerProcess(svc, c);
						// }
		})

		svc.Start(c)
	})
}

// endConsumerProcess facilitates service termination
func endConsumerProcess(svc *Service, c chan os.Signal) {

	// Decrement the offset to escape an endless loop in the service
	svc.InitialOffset = int64(-2)

	// Send a kill command to the input channel to terminate program execution
	go func() {
		c <- os.Kill
		close(c)
	}()
}
