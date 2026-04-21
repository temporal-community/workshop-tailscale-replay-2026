// ABOUTME: Temporal activities for the hello-tsnet exercise.
// ABOUTME: GetIP fetches the public IP; GetLocationInfo geolocates it.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetIP(ctx context.Context) (string, error) {
	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

func GetLocationInfo(ctx context.Context, ip string) (string, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		City       string `json:"city"`
		RegionName string `json:"regionName"`
		Country    string `json:"country"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, %s, %s", data.City, data.RegionName, data.Country), nil
}
