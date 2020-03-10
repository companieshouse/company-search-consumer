package service

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/companieshouse/chs.go/kafka/consumer/cluster"
	"reflect"
)

const resourceKind = "resourceKind"
const resourceUri = "resourceUri"
const contextID = "contextID"
const resourceID = "resourceID"

type MockConsumer struct{}

type MockMarshaller struct{}

type MockUpsert struct{}

type MockUpsertFail struct{}

type MockGroup struct{}

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

func (m MockGroup) MarkOffset(msg *sarama.ConsumerMessage, metadata string) {}

func (m MockGroup) CommitOffsets() error {
	return nil
}

func createMockConsumerWithMessage() *consumer.GroupConsumer {

	return &consumer.GroupConsumer{
		GConsumer: MockConsumer{},
		Group:     MockGroup{},
	}
}

func (m MockMarshaller) Unmarshal(message []byte, s interface{}) error {

	reflect.ValueOf(s).Elem().Field(0).SetString("resourceKind")
	reflect.ValueOf(s).Elem().Field(1).SetString("resourceUri")
	reflect.ValueOf(s).Elem().Field(2).SetString("contextId")
	reflect.ValueOf(s).Elem().Field(3).SetString("resourceId")
	reflect.ValueOf(s).Elem().Field(4).SetString("data")
	return nil
}

func (m MockMarshaller) Marshal(s interface{}) ([]byte, error) {
	return nil, nil
}

func (m MockUpsert) SendViaAPI(data string) (int, error) {
	return 200, nil
}

func (m MockUpsertFail) SendViaAPI(data string) (int, error) {
	return 400, errors.New("failed to send to api")
}

func getDefaultSchema() string {

	return "{ \"type\": \"record\", \"name\": \"resource_changed_data\", \"namespace\": \"stream\", \"fields\": [ { \"name\": \"resource_kind\", \"type\": \"string\" }, { \"name\": \"resource_uri\", \"type\": \"string\" }, { \"name\": \"context_id\", \"type\": \"string\" }, { \"name\": \"resource_id\", \"type\": \"string\" }, { \"name\": \"data\", \"type\": \"string\" }, { \"name\": \"event\", \"type\": { \"type\": \"record\", \"name\": \"event_record\", \"fields\": [ { \"name\": \"published_at\", \"type\": \"string\", \"default\":\"\" }, { \"name\": \"type\", \"type\": \"string\", \"default\":\"\" }, { \"name\": \"fields_changed\", \"type\": [ \"null\", { \"type\": \"array\", \"items\": \"string\" } ] } ] } } ] }"
}
