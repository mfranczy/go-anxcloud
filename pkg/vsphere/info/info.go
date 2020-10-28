// Package info implements API functions residing under /info.
// This path contains methods to query information about created VMs.
package info

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPrefix = "/api/vsphere/v1/info.json"
)

// Info contains meta information of a VM.
type Info struct {
	Name             string `json:"name"`
	CustomName       string `json:"custom_name"`
	Identifier       string `json:"identifier"`
	GuestOS          string `json:"guest_os"`
	LocationCode     string `json:"location_code"`
	LocationCountry  string `json:"location_country"`
	LocationName     string `json:"location_name"`
	Status           string `json:"status"`
	VersionTools     string `json:"version_tools"`
	GuestToolsStatus string `json:"guest_tools_status"`
	RAM              int    `json:"ram"`
	CPU              int    `json:"cpu"`
	Cores            int    `json:"cores"`
	Disks            int    `json:"disks"`
	DiskInfo         []struct {
		DiskType     string `json:"disk_type"`
		StorageType  string `json:"storage_type"`
		BusType      string `json:"bus_type"`
		BusTypeLabel string `json:"bus_type_label"`
		DiskGB       int    `json:"disk_gb"`
		DiskID       int    `json:"disk_id"`
		IOPS         int    `json:"iops"`
		Latency      int    `json:"latence"`
	} `json:"disk_info"`
	Network []struct {
		NIC        int      `json:"nic"`
		ID         int      `json:"id"`
		VLAN       string   `json:"vlan"`
		MACAddress string   `json:"mac_address"`
		IPv4       []string `json:"ips_v4"`
		IPv6       []string `json:"ips_v6"`
	} `json:"network"`
}

// Get returns additional information to a given VM identifier.
//
// ctx is attached to the request and will cancel it on cancelation.
// identifier is the ID of the VM to query.
func (a api) Get(ctx context.Context, identifier string) (Info, error) {
	url := fmt.Sprintf(
		"%s%s/%s/info",
		a.client.BaseURL(),
		pathPrefix,
		identifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Info{}, fmt.Errorf("could not create VM info request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Info{}, fmt.Errorf("could not execute VM info request: %w", err)
	}
	var responsePayload Info
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()

	if err != nil {
		return Info{}, fmt.Errorf("could not decode VM info response: %w", err)
	}

	return responsePayload, err
}