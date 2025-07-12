// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"net/http"
)

type ApiPelicanUserListResponse struct {
	Object string                   `json:"object"`
	Data   []ApiPelicanUserResponse `json:"data"`
	Meta   struct {
		Pagination ApiPelicanListMetaPagination `json:"pagination"`
	} `json:"meta"`
}

type ApiPelicanUserResponse struct {
	Object     string      `json:"object"`
	Attributes PelicanUser `json:"attributes"`
}

type PelicanUser struct {
	ID         int    `json:"id"`
	ExternalID string `json:"external_id"`
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language"`
	RootAdmin  bool   `json:"root_admin"`
	TwoFA      bool   `json:"2fa"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (p *Pelican) GetUsers() (*[]PelicanUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/application/users", p.HostURL), nil)
	if err != nil {
		return nil, err
	}
	apiPelicanUserResponse := &ApiPelicanUserListResponse{}
	if err := p.sendRequest(req, apiPelicanUserResponse); err != nil {
		return nil, err
	}
	users := make([]PelicanUser, len(apiPelicanUserResponse.Data))
	for i, userResponse := range apiPelicanUserResponse.Data {
		users[i] = userResponse.Attributes
	}
	return &users, nil
}

func (p *Pelican) GetUserByID(id int) (*PelicanUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/application/users/%d", p.HostURL, id), nil)
	if err != nil {
		return nil, err
	}
	apiPelicanUserResponse := &ApiPelicanUserResponse{}
	if err := p.sendRequest(req, apiPelicanUserResponse); err != nil {
		return nil, err
	}
	return &apiPelicanUserResponse.Attributes, nil
}
