// Package upsert contains the logic to send the upsert request to the API
package upsert

import (
	"encoding/json"
	"errors"
	"github.com/companieshouse/chs.go/log"
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
	SendViaAPI(data string, apiKey string) error
}

// Template config contains recipient address, http client and send email endpoint
type APIUpsert struct {
	HTTPClient                      *http.Client
	AlphabeticalUpsertCompanyAPIUrl string
	AdvancedUpsertCompanyAPIUrl     string
}

// Company profile delta data json contains fields such as company_number
type CompanyProfileDelta struct {
	CompanyNumber string `json:"company_number"`
}

// SendViaAPI makes a call to the search.api.ch.gov.uk and passes it
// the required data for upserting to alpha_search index
func (upsert *APIUpsert) SendViaAPI(data string, apiKey string) error {

	companyProfileDelta := CompanyProfileDelta{}
	json.Unmarshal([]byte(data), &companyProfileDelta)

	alphabeticalRequest, err := http.NewRequest("PUT", upsert.AlphabeticalUpsertCompanyAPIUrl+companyProfileDelta.CompanyNumber, strings.NewReader(data))
	if err != nil {
		return err
	}
	advancedRequest, err := http.NewRequest("PUT", upsert.AdvancedUpsertCompanyAPIUrl+companyProfileDelta.CompanyNumber, strings.NewReader(data))
	if err != nil {
		return err
	}

	alphabeticalRequest.Header.Add("Content-Type", "application/json")
	advancedRequest.Header.Add("Content-Type", "application/json")
	alphabeticalRequest.Header.Add("Authorization", apiKey)
	advancedRequest.Header.Add("Authorization", apiKey)

	log.Info("Attemting to upsert company " + companyProfileDelta.CompanyNumber + " to alphabetical index")
	alphabeticalResponse, err := upsert.HTTPClient.Do(alphabeticalRequest)
	if err != nil {
		return err
	}
	if err := alphabeticalResponse.Body.Close(); err != nil {
		return err
	}
	if alphabeticalResponse.StatusCode != http.StatusOK {
		return ErrInvalidResponse
	}
	log.Info("Upsert company " + companyProfileDelta.CompanyNumber + " to alphabetical index successful")

	log.Info("Attemting to upsert company " + companyProfileDelta.CompanyNumber + " to advanced index")
	advancedResponse, err := upsert.HTTPClient.Do(advancedRequest)
	if err != nil {
		return err
	}
	if err := advancedResponse.Body.Close(); err != nil {
		return err
	}
	if advancedResponse.StatusCode != http.StatusOK {
		return ErrInvalidResponse
	}
	log.Info("Upsert company " + companyProfileDelta.CompanyNumber + " to advanced index successful")
	return nil
}
