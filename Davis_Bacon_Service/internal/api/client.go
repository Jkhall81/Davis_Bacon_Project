package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://sam.gov/api/prod"

type APIClient struct {
	client  *http.Client
	headers map[string]string
}

func NewClient() *APIClient {
	headers := map[string]string{
		"User-Agent":        "Mozilla/5.0 (DavisBacon-GoScraper)",
		"Accept":            "application/json, text/javascript, */*; q=0.01",
		"X-Requested-Width": "XMLHttpRequest",
	}
	return &APIClient{
		client:  &http.Client{Timeout: 30 * time.Second},
		headers: headers,
	}
}

// Sends GET request and decodes JSON response.
func (a *APIClient) GetJSON(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range a.headers {
		req.Header.Set(k, v)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
