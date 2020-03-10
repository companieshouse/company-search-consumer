package upsert

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func createMockClient(status int) *http.Client {
	mockStreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(status)
	}))
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(mockStreamServer.URL)
		},
	}
	httpClient := &http.Client{Transport: transport}
	return httpClient
}

func TestUnitUpsert(t *testing.T) {
	Convey("Test call to search.api.ch.gov.uk is successful when valid fields passed in", t, func() {
		httpClient := createMockClient(200)

		upsert := &Template{
			HTTPClient:          httpClient,
			UpsertCompanyAPIUrl: "http://api.chs-dev.internal:4089/upsert-company",
		}

		statusCode, err := upsert.SendViaAPI("{'data' : '1' }")
		So(err, ShouldBeNil)
		So(statusCode, ShouldEqual, 200)
	})

	Convey("Test call to search.api.ch.gov.uk returns error when invalid url is passed in", t, func() {
		httpClient := createMockClient(500)

		upsert := &Template{
			HTTPClient:          httpClient,
			UpsertCompanyAPIUrl: "http://invalid-url",
		}
		statusCode, err := upsert.SendViaAPI("{'data' : '1' }")
		So(err, ShouldEqual, ErrInvalidResponse)
		So(statusCode, ShouldEqual, 500)
	})

	Convey("Test search.api.ch.gov.uk returns error when no protocol in front of url", t, func() {
		httpClient := createMockClient(500)

		upsert := &Template{
			HTTPClient:          httpClient,
			UpsertCompanyAPIUrl: "invalid-url",
		}
		statusCode, err := upsert.SendViaAPI("{ 'data' : '1' }")
		So(err, ShouldNotBeNil)
		So(statusCode, ShouldEqual, 400)
	})
}
