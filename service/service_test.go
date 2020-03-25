package service

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	consumer "github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"github.com/companieshouse/company-search-consumer/upsert"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"reflect"
	"testing"
	"time"
)

func createMockService(mockUpsert *upsert.MockUpsert) *Service {

	return &Service{
		InitialOffset: int64(-1),
		Schema:        getDefaultData(),
		Upsert:        mockUpsert,
		Marshaller:    MockMarshaller{},
		HandleError:   MockHandleError,
	}
}

func createMockServiceErrorUnmarshalling(mockUpsert *upsert.MockUpsert) *Service {
	return &Service{
		InitialOffset: int64(-1),
		Schema:        getDefaultData(),
		Upsert:        mockUpsert,
		Marshaller:    MockErrorMarshaller{},
		HandleError:   MockHandleError,
	}
}

func getDefaultData() string {
	return "data"
}

var errorMessage error

func MockHandleError(err error, offset int64, str interface{}) error {
	errorMessage = err
	return nil
}

const resourceKind = "resourceKind"
const resourceUri = "resourceURI"
const contextID = "contextID"
const resourceID = "resourceID"

// Marshaller that works
type MockMarshaller struct{}

func (m MockMarshaller) Unmarshal(message []byte, s interface{}) error {

	reflect.ValueOf(s).Elem().Field(0).SetString("resourceKind")
	reflect.ValueOf(s).Elem().Field(1).SetString("resourceUri")
	reflect.ValueOf(s).Elem().Field(2).SetString("contextId")
	reflect.ValueOf(s).Elem().Field(3).SetString("resourceId")
	reflect.ValueOf(s).Elem().Field(4).SetString(getDefaultData())
	return nil
}

func (m MockMarshaller) Marshal(s interface{}) ([]byte, error) {
	return nil, nil
}

// Marshaller that throws error
type MockErrorMarshaller struct{}

func (m MockErrorMarshaller) Unmarshal(message []byte, s interface{}) error {
	return fmt.Errorf("error unmarshalling data")
}

func (m MockErrorMarshaller) Marshal(s interface{}) ([]byte, error) {
	return nil, nil
}

// MockGroup that works
type MockGroup struct{}

func (m MockGroup) MarkOffset(msg *sarama.ConsumerMessage, metadata string) {}

func (m MockGroup) CommitOffsets() error {
	return nil
}

// MockGroup that throws error
type MockErrorGroup struct{}

func (m MockErrorGroup) MarkOffset(msg *sarama.ConsumerMessage, metadata string) {}

func (m MockErrorGroup) CommitOffsets() error {
	return fmt.Errorf("error committing message offset")
}

// MockConsumer that works
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

func createMockConsumerWithMessage() *consumer.GroupConsumer {

	return &consumer.GroupConsumer{
		GConsumer: MockConsumer{},
		Group:     MockGroup{},
	}
}

func createMockConsumerWithErrorCommitting() *consumer.GroupConsumer {

	return &consumer.GroupConsumer{
		GConsumer: MockConsumer{},
		Group:     MockErrorGroup{},
	}
}

func TestUnitStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Successful process of a single Kafka message", t, func() {
		c := make(chan os.Signal)

		Convey("Given a message is readily available for the service to consume", func() {
			mockUpsert := upsert.NewMockUpsert(ctrl)
			svc := createMockService(mockUpsert)

			svc.Consumer = createMockConsumerWithMessage()

			Convey("Then the upsert request is successfully sent once", func() {
				mockUpsert.EXPECT().SendViaAPI(getDefaultData()).DoAndReturn(func(data string) error {
					// Since this is the last thing the service does, we send a signal to kill the consumer process gracefully
					endConsumerProcess(svc, c)
					return nil
				}).Times(1)

				svc.Start(c)
			})
		})
	})

	Convey("Error unmarshalling document", t, func() {
		c := make(chan os.Signal)

		Convey("Given a message is readily available but there is an error unmarshalling its contents", func() {
			mockUpsert := upsert.NewMockUpsert(ctrl)
			svc := createMockServiceErrorUnmarshalling(mockUpsert)

			svc.Consumer = createMockConsumerWithMessage()

			Convey("Then the upsert request is not sent", func() {
				mockUpsert.EXPECT().SendViaAPI(getDefaultData()).Times(0)

				time.AfterFunc(1*time.Millisecond, func() {
					endConsumerProcess(svc, c)
				})

				Convey("And the correct error is handled", func() {
					svc.Start(c)
					So(errorMessage.Error(), ShouldEqual, "error unmarshalling data")
				})
			})
		})
	})

	Convey("Error calling upsert on the search api", t, func() {
		c := make(chan os.Signal)

		Convey("Given a message is readily available for the service to consume", func() {
			mockUpsert := upsert.NewMockUpsert(ctrl)
			svc := createMockService(mockUpsert)

			svc.Consumer = createMockConsumerWithMessage()

			Convey("When the upsert request fails", func() {
				mockUpsert.EXPECT().SendViaAPI(getDefaultData()).DoAndReturn(func(data string) error {
					// Since this is the last thing the service does, we send a signal to kill the consumer process gracefully
					endConsumerProcess(svc, c)
					return fmt.Errorf("error calling upsert")
				}).Times(1)

				Convey("Then the correct error is handled", func() {
					svc.Start(c)
					So(errorMessage.Error(), ShouldEqual, "error calling upsert")
				})
			})
		})
	})

	Convey("Error committing offset", t, func() {
		c := make(chan os.Signal)

		Convey("Given a message is readily available for the service to consume", func() {
			mockUpsert := upsert.NewMockUpsert(ctrl)
			svc := createMockService(mockUpsert)

			svc.Consumer = createMockConsumerWithErrorCommitting()

			Convey("When the upsert request is successfully sent once", func() {
				mockUpsert.EXPECT().SendViaAPI(getDefaultData()).DoAndReturn(func(data string) error {
					// Since this is the last thing the service does, we send a signal to kill the consumer process gracefully
					endConsumerProcess(svc, c)
					return nil
				}).Times(1)

				Convey("Then the correct error is handled", func() {
					svc.Start(c)
					So(errorMessage.Error(), ShouldEqual, "error committing message offset")
				})
			})
		})
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
	time.Sleep(1 * time.Millisecond)
}
