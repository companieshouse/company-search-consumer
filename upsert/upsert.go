package upsert

import (
	"errors"
	"net/http"
	"strings"
)

var (
	//ErrInvalidResponse is returned when the Search API returns a non 200 status
	ErrInvalidResponse = errors.New("invalid status returned by search.api.ch.gov.uk")
)

// Interface for calling api to upsert data to an elastic search index
type Upsert interface {
	// SendViaApi provides ability to post to search api
	SendViaAPI(data string) (int, error)
}

// Template config contains recipient address, http client and send email endpoint
type APIUpsert struct {
	HTTPClient          *http.Client
	UpsertCompanyAPIUrl string
}

// SendViaAPI makes a call to the search.api.ch.gov.uk and passes it
// the required data for upserting to alpha_search index
func (upsert *APIUpsert) SendViaAPI(data string) (int, error) {

	req, err := http.NewRequest("POST", upsert.UpsertCompanyAPIUrl, strings.NewReader(data))
	if err != nil {
		return 400, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := upsert.HTTPClient.Do(req)
	if err != nil {
		return 400, err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return resp.StatusCode, ErrInvalidResponse
	}
	return resp.StatusCode, nil
}
