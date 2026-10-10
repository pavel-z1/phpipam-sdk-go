// Package locations provides types and methods for working with the
// locations controller (/tools/locations/).
package locations

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/devices"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Location represents a PHPIPAM location.
type Location struct {
	// The location ID.
	ID int `json:"id,string,omitempty"`

	// The location name.
	Name string `json:"name,omitempty"`

	// A detailed description of the location.
	Description string `json:"description,omitempty"`

	// The postal address of the location. If set without coordinates,
	// PHPIPAM looks up the coordinates via OpenStreetMap.
	Address string `json:"address,omitempty"`

	// The latitude of the location.
	Latitude string `json:"lat,omitempty"`

	// The longitude of the location.
	Longitude string `json:"long,omitempty"`
}

// Controller is the base client for the Locations controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Locations controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListLocations lists all locations.
func (c *Controller) ListLocations() (out []Location, err error) {
	err = c.SendRequest("GET", "/tools/locations/", &struct{}{}, &out)
	return
}

// CreateLocation creates a location by sending a POST request.
func (c *Controller) CreateLocation(in Location) (message string, err error) {
	err = c.SendRequest("POST", "/tools/locations/", &in, &message)
	return
}

// GetLocationByID GETs a location via its ID in the PHPIPAM database.
func (c *Controller) GetLocationByID(id int) (out Location, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/locations/%d/", id), &struct{}{}, &out)
	return
}

// GetSubnetsInLocation GETs the subnets assigned to a location.
func (c *Controller) GetSubnetsInLocation(id int) (out []subnets.Subnet, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/locations/%d/subnets/", id), &struct{}{}, &out)
	return
}

// GetDevicesInLocation GETs the devices assigned to a location.
func (c *Controller) GetDevicesInLocation(id int) (out []devices.Device, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/locations/%d/devices/", id), &struct{}{}, &out)
	return
}

// GetAddressesInLocation GETs the IP addresses assigned to a location.
func (c *Controller) GetAddressesInLocation(id int) (out []addresses.Address, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/locations/%d/ipaddresses/", id), &struct{}{}, &out)
	return
}

// UpdateLocation updates a location by sending a PATCH request.
//
// The tools controller takes the object ID from the URL; an "id" field in the
// request body would override the subcontroller name, so it is not sent.
func (c *Controller) UpdateLocation(in Location) (message string, err error) {
	id := in.ID
	in.ID = 0
	err = c.SendRequest("PATCH", fmt.Sprintf("/tools/locations/%d/", id), &in, &message)
	return
}

// DeleteLocation deletes a location by its ID.
func (c *Controller) DeleteLocation(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/tools/locations/%d/", id), &struct{}{}, &message)
	return
}
