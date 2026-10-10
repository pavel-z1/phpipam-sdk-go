// Package devices provides types and methods for working with the devices
// controller.
package devices

import (
	"fmt"
	"net/url"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Device represents a PHPIPAM device.
type Device struct {
	// The device ID.
	ID int `json:"id,string,omitempty"`

	// The device hostname. Required when creating a device.
	Hostname string `json:"hostname,omitempty"`

	// The IP address of the device.
	IPAddress string `json:"ip,omitempty"`

	// The ID of the device type.
	Type int `json:"type,string,omitempty"`

	// A detailed description of the device.
	Description string `json:"description,omitempty"`

	// A semicolon-separated list of section IDs this device belongs to.
	Sections string `json:"sections,omitempty"`

	// The ID of the rack the device is mounted in.
	Rack int `json:"rack,string,omitempty"`

	// The starting rack unit of the device.
	RackStart int `json:"rack_start,string,omitempty"`

	// The size of the device in rack units.
	RackSize int `json:"rack_size,string,omitempty"`

	// The ID of the location of the device.
	Location int `json:"location,string,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration.
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// Controller is the base client for the devices controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the devices
// controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListDevices lists all devices.
func (c *Controller) ListDevices() (out []Device, err error) {
	err = c.SendRequest("GET", "/devices/", &struct{}{}, &out)
	return
}

// CreateDevice creates a device by sending a POST request.
func (c *Controller) CreateDevice(in Device) (message string, err error) {
	err = c.SendRequest("POST", "/devices/", &in, &message)
	return
}

// GetDeviceByID GETs a device via its ID in the PHPIPAM database.
func (c *Controller) GetDeviceByID(id int) (out Device, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/devices/%d/", id), &struct{}{}, &out)
	return
}

// SearchDevices searches for devices whose hostname, IP address, description
// or custom fields match the supplied search string.
func (c *Controller) SearchDevices(searchString string) (out []Device, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/devices/search/%s/", url.PathEscape(searchString)), &struct{}{}, &out)
	return
}

// GetAddressesOnDevice GETs the IP addresses assigned to a device.
func (c *Controller) GetAddressesOnDevice(id int) (out []addresses.Address, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/devices/%d/addresses/", id), &struct{}{}, &out)
	return
}

// GetSubnetsOnDevice GETs the subnets assigned to a device.
func (c *Controller) GetSubnetsOnDevice(id int) (out []subnets.Subnet, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/devices/%d/subnets/", id), &struct{}{}, &out)
	return
}

// UpdateDevice updates a device by sending a PATCH request.
func (c *Controller) UpdateDevice(in Device) (message string, err error) {
	err = c.SendRequest("PATCH", "/devices/", &in, &message)
	return
}

// DeleteDevice deletes a device by its ID.
func (c *Controller) DeleteDevice(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/devices/%d/", id), &struct{}{}, &message)
	return
}
