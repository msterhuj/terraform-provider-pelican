// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// https://dashflo.net/docs/api/pterodactyl/v1/

const HostURL = "https://127.0.0.1:8080" // Default url to pelican server

type Pelican struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(hostURL, token string) (*Pelican, error) {
	if hostURL == "" {
		hostURL = HostURL
	}

	if token == "" {
		return nil, fmt.Errorf("token must be provided")
	}

	return &Pelican{
		HostURL: hostURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Token: token,
	}, nil
}

type ApiPelicanListMetaPagination struct {
	Total       int                    `json:"total"`
	Count       int                    `json:"count"`
	PerPage     int                    `json:"per_page"`
	CurrentPage int                    `json:"current_page"`
	TotalPages  int                    `json:"total_pages"`
	Links       map[string]interface{} `json:"links"`
}

func (p *Pelican) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.Token))

	res, err := p.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// print body for debugging
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("unexpected status code: %d, Body: %s, RequestURI: %s", res.StatusCode, body, req.URL)
	}

	if v != nil {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}
