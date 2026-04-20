// ABOUTME: Temporal Activities for the Go weather agent.
// ABOUTME: Includes OpenAI API calls and tool implementations (weather, IP, location).

package main

import (
	"context"
)

// Activities holds shared state for activity implementations.
type Activities struct {
	// TODO: Add OpenAI base URL and any shared configuration
	OpenAIBaseURL string
}

// CreateCompletion calls the OpenAI API (through Aperture) with function calling.
func (a *Activities) CreateCompletion(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Implement OpenAI Responses API call via Aperture
	//
	// 1. Build the HTTP request to OPENAI_BASE_URL/responses
	// 2. Include the model, instructions, input, and tools
	// 3. Parse the response to determine if the LLM wants to call a tool
	//
	// The Aperture endpoint is OpenAI-compatible, so use the same API format.

	return nil, nil
}

// GetWeatherAlerts fetches active weather alerts for a US state from the NWS API.
func (a *Activities) GetWeatherAlerts(ctx context.Context, state string) (string, error) {
	// TODO: Call https://api.weather.gov/alerts/active/area/{state}
	return "", nil
}

// GetIPAddress returns the public IP of the machine running this activity.
func (a *Activities) GetIPAddress(ctx context.Context) (string, error) {
	// TODO: Call https://icanhazip.com
	return "", nil
}

// GetLocationInfo geolocates an IP address and returns city, region, country.
func (a *Activities) GetLocationInfo(ctx context.Context, ip string) (string, error) {
	// TODO: Call http://ip-api.com/json/{ip}
	return "", nil
}
