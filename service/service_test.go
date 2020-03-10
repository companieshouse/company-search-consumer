package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

func createMockService() *Service {

	return &Service{
		Marshaller:    MockMarshaller{},
		InitialOffset: int64(-1),
		Schema:        getDefaultSchema(),
	}
}

func TestUnitStart(t *testing.T) {
	Convey("Successful process of a single Kafka message", t, func() {

		c := make(chan os.Signal)
		svc := createMockService()

		Convey("Given a message is readily available for the service to consume", func() {
			svc.Consumer = createMockConsumerWithMessage()

			Convey("When the upsert request is successful", func() {

				svc.Upsert = MockUpsert{}

				Convey("Then the consumer will close succesfully", func() {
					time.AfterFunc(1*time.Millisecond, func() {
						endConsumerProcess(svc, c)
					})
				})
				svc.Start(c)
			})
		})
		So(testHttpResponseValue, ShouldEqual, 200)
	})

	Convey("Error when attempting to upsert a document", t, func() {

		c := make(chan os.Signal)
		svc := createMockService()

		Convey("Given a message is readily available for the service to consume", func() {
			svc.Consumer = createMockConsumerWithMessage()

			Convey("When the upsert request fails", func() {

				svc.Upsert = MockUpsertFail{}

				Convey("Then the consumer will close succesfully", func() {
					time.AfterFunc(1*time.Millisecond, func() {
						endConsumerProcess(svc, c)
					})
				})
				svc.Start(c)
			})
		})
		So(testHttpResponseValue, ShouldEqual, 400)
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
