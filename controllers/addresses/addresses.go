// Package addresses provides types and methods for working with the addresses
// controller.
package addresses

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Address represents an IP address resource within PHPIPAM.
type AddressDTO struct {
	// The ID of the IP address entry within PHPIPAM.
	ID phpipam.JSONIntString `json:"id,omitempty"`

	// The ID of the subnet that the address belongs to.
	SubnetID phpipam.JSONIntString `json:"subnetId,omitempty"`

	// The IP address, without a CIDR subnet mask.
	IPAddress string `json:"ip,omitempty"`

	// true if this IP address is a gateway address.
	IsGateway phpipam.BoolIntString `json:"is_gateway,omitempty"`

	// A detailed description of the IP address entry.
	Description string `json:"description,omitempty"`

	// A hostname for the IP address.
	Hostname string `json:"hostname,omitempty"`

	// The MAC address for the IP.
	MACAddress string `json:"mac,omitempty"`

	// The address owner (customer, hostname, application, etc).
	Owner string `json:"owner,omitempty"`

	// The tag ID for the IP address.
	Tag phpipam.JSONIntString `json:"tag,omitempty"`

	// true if PTR records should not be created for this IP address.
	PTRIgnore phpipam.BoolIntString `json:"PTRIgnore,omitempty"`

	// The ID of a PowerDNS PTR record.
	PTRRecordID phpipam.JSONIntString `json:"PTR,omitempty"`

	// An ID of a device that this address belongs to.
	DeviceID phpipam.JSONIntString `json:"deviceId,omitempty"`

	// A switchport number/label that this IP address belongs to.
	Port string `json:"port,omitempty"`

	// A note for this IP address, detailing state information not sutiable for
	// entering in the description.
	Note string `json:"note,omitempty"`

	// A timestamp for when the address was last seen with ping.
	LastSeen string `json:"lastSeen,omitempty"`

	// true if you want to exclude this address from ping scans.
	ExcludePing phpipam.BoolIntString `json:"excludePing,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration. If this is not
	// enabled, this map will be nil on GETs and POSTs and PATCHes with this
	// field set will fail. Use the explicit custom field functions instead.
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type Address struct {
	// The ID of the IP address entry within PHPIPAM.
	ID int `json:"id,omitempty"`

	// The ID of the subnet that the address belongs to.
	SubnetID int `json:"subnetId,omitempty"`

	// The IP address, without a CIDR subnet mask.
	IPAddress string `json:"ip,omitempty"`

	// true if this IP address is a gateway address.
	IsGateway phpipam.BoolIntString `json:"is_gateway,omitempty"`

	// A detailed description of the IP address entry.
	Description string `json:"description,omitempty"`

	// A hostname for the IP address.
	Hostname string `json:"hostname,omitempty"`

	// The MAC address for the IP.
	MACAddress string `json:"mac,omitempty"`

	// The address owner (customer, hostname, application, etc).
	Owner string `json:"owner,omitempty"`

	// The tag ID for the IP address.
	Tag int `json:"tag,omitempty"`

	// true if PTR records should not be created for this IP address.
	PTRIgnore phpipam.BoolIntString `json:"PTRIgnore,omitempty"`

	// The ID of a PowerDNS PTR record.
	PTRRecordID int `json:"PTR,omitempty"`

	// An ID of a device that this address belongs to.
	DeviceID int `json:"deviceId,omitempty"`

	// A switchport number/label that this IP address belongs to.
	Port string `json:"port,omitempty"`

	// A note for this IP address, detailing state information not sutiable for
	// entering in the description.
	Note string `json:"note,omitempty"`

	// A timestamp for when the address was last seen with ping.
	LastSeen string `json:"lastSeen,omitempty"`

	// true if you want to exclude this address from ping scans.
	ExcludePing phpipam.BoolIntString `json:"excludePing,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration. If this is not
	// enabled, this map will be nil on GETs and POSTs and PATCHes with this
	// field set will fail. Use the explicit custom field functions instead.
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

func (a *Address) FromDTO(dto *AddressDTO) {
	a.ID = int(dto.ID)
	a.SubnetID = int(dto.SubnetID)
	a.IPAddress = dto.IPAddress
	a.IsGateway = dto.IsGateway
	a.Description = dto.Description
	a.Hostname = dto.Hostname
	a.MACAddress = dto.MACAddress
	a.Owner = dto.Owner
	a.Tag = int(dto.Tag)
	a.PTRIgnore = dto.PTRIgnore
	a.PTRRecordID = int(dto.PTRRecordID)
	a.DeviceID = int(dto.DeviceID)
	a.Port = dto.Port
	a.Note = dto.Note
	a.LastSeen = dto.LastSeen
	a.ExcludePing = dto.ExcludePing
	a.EditDate = dto.EditDate
	a.CustomFields = dto.CustomFields
}

func (a *Address) ToDTO() AddressDTO {
	return AddressDTO{
		ID:          phpipam.JSONIntString(a.ID),
		SubnetID:    phpipam.JSONIntString(a.SubnetID),
		IPAddress:   a.IPAddress,
		IsGateway:   a.IsGateway,
		Description: a.Description,
		Hostname:    a.Hostname,
		MACAddress:  a.MACAddress,
		Owner:       a.Owner,
		Tag:         phpipam.JSONIntString(a.Tag),
		PTRIgnore:   a.PTRIgnore,
		PTRRecordID: phpipam.JSONIntString(a.PTRRecordID),
		DeviceID:    phpipam.JSONIntString(a.DeviceID),
		Port:        a.Port,
		Note:        a.Note,
		LastSeen:    a.LastSeen,
		ExcludePing: a.ExcludePing,
		EditDate:    a.EditDate,
		CustomFields: a.CustomFields,
	}
}

// Controller is the base client for the Addresses controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Addresses controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// CreateAddress creates an address by sending a POST request.
func (c *Controller) CreateAddress(in Address) (message string, err error) {
	err = c.SendRequest("POST", "/addresses/", in.ToDTO(), &message)
	return
}

// CreateAddress creates a first free in subnet address by sending a POST request.
func (c *Controller) CreateFirstFreeAddress(id int, in Address) (out string, err error) {
        err = c.SendRequest("POST", fmt.Sprintf("/addresses/first_free/%d/", id), in.ToDTO(), &out)
        return
}

// GetAddressByID GETs an address via its ID.
func (c *Controller) GetAddressByID(id int) (out Address, err error) {
	var dto AddressDTO
	err = c.SendRequest("GET", fmt.Sprintf("/addresses/%d/", id), &struct{}{}, &dto)
	out.FromDTO(&dto)
	return
}

// GetAddressesByIP searches for an address by its IP.
//
// According to the spec, this can return multiple addresses, however it's not
// entirely clear how to perform a search that would yield multiple results.
func (c *Controller) GetAddressesByIP(ipaddr string) (out []Address, err error) {
	var dtos []AddressDTO
	err = c.SendRequest("GET", fmt.Sprintf("/addresses/search/%s/", ipaddr), &struct{}{}, &dtos)
	for _, dto := range dtos {
		var a Address
		a.FromDTO(&dto)
		out = append(out, a)
	}
	return
}

// GetAddressesByIP searches for an address by its IP with in given subnet
// When having multiple subnets with same ip range this will return the address in the given subnet
// Those subnet may not talk to each other but still exist under on phpIPAM instance especially on ones migrated from previous versions 
func (c *Controller) GetAddressesByIpInSubnet(ipaddr string,subnetID int) (out Address, err error) {
	var dto AddressDTO
	err = c.SendRequest("GET", fmt.Sprintf("/addresses/%s/%d", ipaddr,subnetID), &struct{}{}, &dto)
	out.FromDTO(&dto)
	return
}

// GetAddressCustomFieldsSchema GETs the custom fields for the addresses controller via
// client.GetCustomFieldsSchema.
func (c *Controller) GetAddressCustomFieldsSchema() (out map[string]phpipam.CustomField, err error) {
	out, err = c.Client.GetCustomFieldsSchema("addresses")
	return
}

// GetAddressCustomFields GETs the custom fields for a subnet via
// client.GetCustomFields.
func (c *Controller) GetAddressCustomFields(id int) (out map[string]interface{}, err error) {
	out, err = c.Client.GetCustomFields(id, "addresses")
	return
}

// UpdateAddress updates an address by sending a PATCH request.
func (c *Controller) UpdateAddress(in Address) (message string, err error) {
	err = c.SendRequest("PATCH", "/addresses/", in.ToDTO(), &message)
	return
}

// UpdateAddressCustomFields PATCHes the subnet's custom fields via
// client.UpdateCustomFields.
func (c *Controller) UpdateAddressCustomFields(id int, in map[string]interface{}) (message string, err error) {
	message, err = c.Client.UpdateCustomFields(id, in, "addresses")
	return
}

// DeleteAddress deletes an address by ID. RemoveDNS can be set to true if you
// want to have any related DNS records deleted as well.
func (c *Controller) DeleteAddress(id int, RemoveDNS phpipam.BoolIntString) (message string, err error) {
	in := struct {
		RemoveDNS phpipam.BoolIntString `json:"remove_dns,omitempty"`
	}{
		RemoveDNS: RemoveDNS,
	}
	err = c.SendRequest("DELETE", fmt.Sprintf("/addresses/%d/", id), &in, &message)
	return
}
