// Package upsert contains the logic to send the upsert request to the API
package upsert

import (
	"errors"
	"net/http"
	"strings"
)

var (
	// ErrInvalidResponse is returned when the Search API returns a non 200 status
	ErrInvalidResponse = errors.New("invalid status returned by search.api.ch.gov.uk")
)

// Upsert is an interface for calling api to upsert data to an elastic search index
type Upsert interface {
	// SendViaAPI provides ability to post to search api
	SendViaAPI(data string) error
}

// Template config contains recipient address, http client and send email endpoint
type APIUpsert struct {
	HTTPClient          *http.Client
	UpsertCompanyAPIUrl string
}

// SendViaAPI makes a call to the search.api.ch.gov.uk and passes it
// the required data for upserting to alpha_search index
func (upsert *APIUpsert) SendViaAPI(data string) error {

	req, err := http.NewRequest("PUT", upsert.UpsertCompanyAPIUrl, strings.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := upsert.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrInvalidResponse
	}
	return nil
}
