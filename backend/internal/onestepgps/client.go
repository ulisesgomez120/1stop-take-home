package onestepgps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const requestTimeout = 10 * time.Second

// Client fetches device data from the OneStepGPS API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

// FetchDevices calls the OneStepGPS device endpoint and maps the response
// into the domain Device type.
func (c *Client) FetchDevices(ctx context.Context) ([]Device, error) {
	url := fmt.Sprintf("%s?latest_point=true&api-key=%s", c.baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch devices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch devices: unexpected status %d", resp.StatusCode)
	}

	var raw rawResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return mapDevices(raw), nil
}
