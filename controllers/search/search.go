// Package search provides types and methods for working with the search
// controller.
package search

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/url"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// subnetSearchResult is used internally to parse subnet search responses.
type subnetSearchResult struct {
	Subnets struct {
		Code int             `json:"code,string"`
		Data json.RawMessage `json:"data"`
	} `json:"subnets"`
}

// Controller is the base client for the Search controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Search controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// SearchSubnets searches for subnets in the phpIPAM database. An empty slice
// is returned when nothing matches.
//
// This requires PHPIPAM 1.6 or higher.
func (c *Controller) SearchSubnets(searchString string) (out []subnets.Subnet, err error) {
	var result subnetSearchResult
	err = c.SendRequest("GET", fmt.Sprintf("/search/%s/?addresses=0&subnets=1&vlans=0&vrfs=0", url.PathEscape(searchString)), &struct{}{}, &result)
	if err != nil {
		return
	}

	// phpIPAM reports an empty search result with code 404.
	if result.Subnets.Code == 404 {
		return []subnets.Subnet{}, nil
	}

	if result.Subnets.Code != 200 {
		var errorMsg string
		if err := json.Unmarshal(result.Subnets.Data, &errorMsg); err == nil {
			return nil, fmt.Errorf("search API returned code %d: %s", result.Subnets.Code, errorMsg)
		}
		return nil, fmt.Errorf("search API returned non-200 code: %d", result.Subnets.Code)
	}

	if err = json.Unmarshal(result.Subnets.Data, &out); err != nil {
		return nil, fmt.Errorf("failed to parse subnet data: %w", err)
	}

	// For some reason, subnet addresses are returned as integers: https://github.com/phpipam/phpipam/issues/4159
	for i := range out {
		ip, err := integerToIP(out[i].SubnetAddress)
		if err != nil {
			return out, err
		}
		out[i].SubnetAddress = ip
	}

	return
}

// maxIPv4 is the largest integer representation of an IPv4 address.
var maxIPv4 = big.NewInt(0xFFFFFFFF)

// integerToIP converts the decimal integer representation of an IP address,
// as returned by the search controller, to its textual form. Following
// phpIPAM's own identify_address logic, values up to 2^32-1 are treated as
// IPv4 and larger values as IPv6. Addresses that are already in textual form
// are returned unchanged.
func integerToIP(s string) (string, error) {
	if ip := net.ParseIP(s); ip != nil {
		return ip.String(), nil
	}
	n, ok := new(big.Int).SetString(s, 10)
	if !ok || n.Sign() < 0 || n.BitLen() > 128 {
		return "", fmt.Errorf("invalid integer IP address %q", s)
	}
	if n.Cmp(maxIPv4) <= 0 {
		b := make([]byte, 4)
		return net.IP(n.FillBytes(b)).String(), nil
	}
	b := make([]byte, 16)
	return net.IP(n.FillBytes(b)).String(), nil
}
