package upsert

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/smartystreets/goconvey/convey"

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

		upsert := &APIUpsert{
			HTTPClient:                      httpClient,
			AlphabeticalUpsertCompanyAPIUrl: "http://api.chs-dev.internal:4089/alphabetical-search/companies/",
			AdvancedUpsertCompanyAPIUrl:     "http://api.chs-dev.internal:4089/advanced-search/companies/",
		}

		err := upsert.SendViaAPI("{'data' : '1' }", "testAPIKey")
		So(err, ShouldBeNil)
	})

	Convey("Test call to search.api.ch.gov.uk returns error when invalid url is passed in alphabetical", t, func() {
		httpClient := createMockClient(500)

		upsert := &APIUpsert{
			HTTPClient:                      httpClient,
			AlphabeticalUpsertCompanyAPIUrl: "http://invalid-url",
			AdvancedUpsertCompanyAPIUrl:     "http://api.chs-dev.internal:4089/advanced-search/companies/",
		}
		err := upsert.SendViaAPI("{'data' : '1' }", "testAPIKey")
		So(err, ShouldEqual, ErrInvalidResponse)
	})

	Convey("Test call to search.api.ch.gov.uk returns error when invalid url is passed in advanced", t, func() {
		httpClient := createMockClient(500)

		upsert := &APIUpsert{
			HTTPClient:                      httpClient,
			AlphabeticalUpsertCompanyAPIUrl: "http://api.chs-dev.internal:4089/alphabetical-search/companies/",
			AdvancedUpsertCompanyAPIUrl:     "http://invalid-url",
		}
		err := upsert.SendViaAPI("{'data' : '1' }", "testAPIKey")
		So(err, ShouldEqual, ErrInvalidResponse)
	})

	Convey("Test search.api.ch.gov.uk returns error when no protocol in front of url alphabetical", t, func() {
		httpClient := createMockClient(500)

		upsert := &APIUpsert{
			HTTPClient:                      httpClient,
			AlphabeticalUpsertCompanyAPIUrl: "invalid-url",
			AdvancedUpsertCompanyAPIUrl:     "http://api.chs-dev.internal:4089/advanced-search/companies/",
		}
		err := upsert.SendViaAPI("{ 'data' : '1' }", "testAPIKey")
		So(err, ShouldNotBeNil)
	})

	Convey("Test search.api.ch.gov.uk returns error when no protocol in front of url advanced", t, func() {
		httpClient := createMockClient(500)

		upsert := &APIUpsert{
			HTTPClient:                      httpClient,
			AlphabeticalUpsertCompanyAPIUrl: "http://api.chs-dev.internal:4089/alphabetical-search/companies/",
			AdvancedUpsertCompanyAPIUrl:     "invalid-url",
		}
		err := upsert.SendViaAPI("{ 'data' : '1' }", "testAPIKey")
		So(err, ShouldNotBeNil)
	})
}
