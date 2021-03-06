// Package address implements API functions residing under /ipam/address.
// This path contains methods for managing IPs.
package address

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathAddressPrefix        = "/api/ipam/v1/address.json"
	pathReserveAddressPrefix = "/api/ipam/v1/address/reserve/ip/count.json"
)

// Address contains all the information about a specific address.
type Address struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	DescriptionInternal string `json:"description_internal"`
	Role                string `json:"role"`
	Version             int    `json:"version"`
	Status              string `json:"status"`
	VLANID              string `json:"vlan"`
	PrefixID            string `json:"prefix"`
}

// Summary is the address information returned by a listing.
type Summary struct {
	ID                  string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	Role                string `json:"role"`
}

// Update contains fields to change on a prefix.
type Update struct {
	Name                string `json:"name,omitempty"`
	DescriptionCustomer string `json:"description_customer,omitempty"`
	Role                string `json:"role,omitempty"`
}

// Create defines meta data of an address to create.
type Create struct {
	PrefixID            string `json:"prefix"`
	Address             string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
	Role                string `json:"role"`
	Organization        string `json:"organization"`
}

// ReserveRandom defines metadata of addresses to reserve randomly.
type ReserveRandom struct {
	LocationID string `json:"location_identifier"`
	VlanID     string `json:"vlan_identifier"`
	Count      int    `json:"count"`
}

// ReserveRandomSummary is the reserved IPs information returned by list request.
type ReserveRandomSummary struct {
	Limit      int          `json:"limit"`
	Page       int          `json:"page"`
	TotalItems int          `json:"total_items"`
	TotalPages int          `json:"total_pages"`
	Data       []ReservedIP `json:"data"`
}

// ReservedIP returns details about reserved ip.
type ReservedIP struct {
	ID      string `json:"identifier"`
	Address string `json:"text"`
	Prefix  string `json:"prefix"`
}

type listResponse struct {
	Data struct {
		Data []Summary `json:"data"`
	} `json:"data"`
}

// NewCreate creates a new address definition with required vlaues.
func NewCreate(prefixID string, address string) Create {
	return Create{
		PrefixID: prefixID,
		Address:  address,
		Role:     "Default",
	}
}

func (a api) List(ctx context.Context, page, limit int) ([]Summary, error) {
	url := fmt.Sprintf(
		"%s%s?page=%v&limit=%v",
		a.client.BaseURL(),
		pathAddressPrefix, page, limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create address list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute address list request: %w", err)
	}
	var responsePayload listResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("could not decode address list response: %w", err)
	}

	return responsePayload.Data.Data, err
}

func (a api) Get(ctx context.Context, id string) (Address, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathAddressPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Address{}, fmt.Errorf("could not create address get request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Address{}, fmt.Errorf("could not execute address get request: %w", err)
	}
	var responsePayload Address
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return Address{}, fmt.Errorf("could not decode address get response: %w", err)
	}

	return responsePayload, err
}

func (a api) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathAddressPrefix,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create address delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute address delete request: %w", err)
	}

	return httpResponse.Body.Close()
}

func (a api) Create(ctx context.Context, create Create) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathAddressPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(create); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan creation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create vlan post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute vlan post request: %w", err)
	}
	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Summary{}, fmt.Errorf("could not decode vlan post response: %w", err)
	}

	return summary, nil
}

func (a api) Update(ctx context.Context, id string, update Update) (Summary, error) {
	url := fmt.Sprintf(
		"%s%s/%s",
		a.client.BaseURL(),
		pathAddressPrefix, id,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(update); err != nil {
		panic(fmt.Sprintf("could not create request data for vlan update: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &requestData)
	if err != nil {
		return Summary{}, fmt.Errorf("could not create vlan update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Summary{}, fmt.Errorf("could not execute vlan update request: %w", err)
	}
	var summary Summary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return summary, fmt.Errorf("could not decode vlan update response: %w", err)
	}

	return summary, err
}

func (a api) ReserveRandom(ctx context.Context, reserve ReserveRandom) (ReserveRandomSummary, error) {
	url := fmt.Sprintf(
		"%s%s",
		a.client.BaseURL(),
		pathReserveAddressPrefix,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(reserve); err != nil {
		panic(fmt.Sprintf("could not create request data for IP address reservation: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return ReserveRandomSummary{}, fmt.Errorf("could not create IP address reserve random post request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return ReserveRandomSummary{}, fmt.Errorf("could not execute IP address reserve random post request: %w", err)
	}
	var summary ReserveRandomSummary
	err = json.NewDecoder(httpResponse.Body).Decode(&summary)
	_ = httpResponse.Body.Close()
	if err != nil {
		return ReserveRandomSummary{}, fmt.Errorf("could not decode IP address reserve random post response: %w", err)
	}

	return summary, nil
}
