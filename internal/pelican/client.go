// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
	"time"
)

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
