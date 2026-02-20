package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Shared HTTP client with a global timeout for all outbound requests.
var Client = &http.Client{
	Timeout: 8 * time.Second,
}

// GetJSON performs a GET request to the given URL and decodes the JSON response
// into the provided target. It returns the upstream HTTP status for caller-side handling.
func GetJSON(url string, target any) (int, error) {
	resp, err := Client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("upstream returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return http.StatusBadGateway, fmt.Errorf("failed to decode response")
	}

	return resp.StatusCode, nil
}
