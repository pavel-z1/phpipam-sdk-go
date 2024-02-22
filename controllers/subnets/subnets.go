// Package subnets provides types and methods for working with the subnets
// controller.
package subnets

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Subnet represents a PHPIPAM subnet.
type SubnetDTO struct {
	// The subnet ID.
	ID phpipam.JSONIntString `json:"id,omitempty"`

	// The subnet address, in dotted quad format (i.e. A.B.C.D).
	SubnetAddress string `json:"subnet,omitempty"`

	// The subnet's mask in number of bits (i.e. 24).
	Mask phpipam.JSONIntString `json:"mask,omitempty"`

	// A detailed description of the subnet.
	Description string `json:"description,omitempty"`

	// The section ID to add the subnet to (required when adding).
	SectionID phpipam.JSONIntString `json:"sectionId,omitempty"`

	// The ID of a linked IPv6 subnet.
	LinkedSubnet phpipam.JSONIntString `json:"linked_subnet,omitempty"`

	// The ID of the VLAN that this subnet belongs to.
	VLANID phpipam.JSONIntString `json:"vlanId,omitempty"`

	// The ID of the VRF this subnet belongs to.
	VRFID phpipam.JSONIntString `json:"vrfId,omitempty"`

	// The parent subnet ID if this is a nested subnet.
	MasterSubnetID phpipam.JSONIntString `json:"masterSubnetId,omitempty"`

	// The ID of the nameserver to attache the subnet to.
	NameserverID phpipam.JSONIntString `json:"nameserverId,omitempty"`

	// The ID and IPs of the nameservers for the subnet
	Nameservers map[string]interface{} `json:"nameservers,omitempty"`

	// true if the name should be displayed in listing instead of the subnet
	// address.
	ShowName phpipam.BoolIntString `json:"showName,omitempty"`

	// A JSON object, stringified, that represents the permissions for this
	// section.
	Permissions string `json:"permissions,omitempty"`

	// Controls if PTR records should be created for the subnet.
	DNSRecursive phpipam.BoolIntString `json:"DNSrecursive,omitempty"`

	// Controls if DNS hostname records are displayed.
	DNSRecords phpipam.BoolIntString `json:"DNSrecords,omitempty"`

	// Controls if IP requests are allowed for the subnet.
	AllowRequests phpipam.BoolIntString `json:"allowRequests,omitempty"`

	// The ID of the scan agent to use for the subnet.
	ScanAgent phpipam.JSONIntString `json:"scanAgent,omitempty"`

	// Controls if the subnet should be included in status checks.
	PingSubnet phpipam.BoolIntString `json:"pingSubnet,omitempty"`

	// Controls if new hosts should be discovered for new host scans.
	DiscoverSubnet phpipam.BoolIntString `json:"discoverSubnet,omitempty"`

	// Controls if we are adding a subnet or folder.
	IsFolder phpipam.BoolIntString `json:"isFolder,omitempty"`

	// Marks the subnet as permitting allocation of the network and broadcast addresses.
	IsPool phpipam.BoolIntString `json:"isPool,omitempty"`

	// Marks the subnet as used.
	IsFull phpipam.BoolIntString `json:"isFull,omitempty"`

	// The threshold of the subnet.
	Threshold phpipam.JSONIntString `json:"threshold,omitempty"`

	// The location index of the subnet.
	Location phpipam.JSONIntString `json:"location,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// Gateway IP and ID of Gateway IP
	Gateway map[string]interface{} `json:"gateway,omitempty"`

	// Gateway IP ID
	GatewayID string `json:"gatewayId,omitempty"`

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration. If this is not
	// enabled, this map will be nil on GETs and POSTs and PATCHes with this
	// field set will fail. Use the explicit custom field functions instead.
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`

	// Controls enabling resolve DNS function.
	ResolveDNS phpipam.BoolIntString `json:"resolveDNS,omitempty"`
}

type Subnet struct {
	// The subnet ID.
	ID int

	// The subnet address, in dotted quad format (i.e. A.B.C.D).
	SubnetAddress string

	// The subnet's mask in number of bits (i.e. 24).
	Mask phpipam.JSONIntString

	// A detailed description of the subnet.
	Description string

	// The section ID to add the subnet to (required when adding).
	SectionID int

	// The ID of a linked IPv6 subnet.
	LinkedSubnet int

	// The ID of the VLAN that this subnet belongs to.
	VLANID int

	// The ID of the VRF this subnet belongs to.
	VRFID int

	// The parent subnet ID if this is a nested subnet.
	MasterSubnetID int

	// The ID of the nameserver to attache the subnet to.
	NameserverID int

	// The ID and IPs of the nameservers for the subnet
	Nameservers map[string]interface{}

	// true if the name should be displayed in listing instead of the subnet
	// address.
	ShowName phpipam.BoolIntString

	// A JSON object, stringified, that represents the permissions for this
	// section.
	Permissions string

	// Controls if PTR records should be created for the subnet.
	DNSRecursive phpipam.BoolIntString

	// Controls if DNS hostname records are displayed.
	DNSRecords phpipam.BoolIntString

	// Controls if IP requests are allowed for the subnet.
	AllowRequests phpipam.BoolIntString

	// The ID of the scan agent to use for the subnet.
	ScanAgent int

	// Controls if the subnet should be included in status checks.
	PingSubnet phpipam.BoolIntString

	// Controls if new hosts should be discovered for new host scans.
	DiscoverSubnet phpipam.BoolIntString

	// Controls if we are adding a subnet or folder.
	IsFolder phpipam.BoolIntString

	// Marks the subnet as permitting allocation of the network and broadcast addresses.
	IsPool phpipam.BoolIntString

	// Marks the subnet as used.
	IsFull phpipam.BoolIntString

	// The threshold of the subnet.
	Threshold int

	// The location index of the subnet.
	Location int

	// The date of the last edit to this resource.
	EditDate string

	// Gateway IP and ID of Gateway IP
	Gateway map[string]interface{}

	// Gateway IP ID
	GatewayID string

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration. If this is not
	// enabled, this map will be nil on GETs and POSTs and PATCHes with this
	// field set will fail. Use the explicit custom field functions instead.
	CustomFields map[string]interface{}

	// Controls enabling resolve DNS function.
	ResolveDNS phpipam.BoolIntString
}

func (s *Subnet) FromDTO(subnetDTO *SubnetDTO) {
	s.ID = int(subnetDTO.ID)
	s.SubnetAddress = subnetDTO.SubnetAddress
	s.Mask = subnetDTO.Mask
	s.Description = subnetDTO.Description
	s.SectionID = int(subnetDTO.SectionID)
	s.LinkedSubnet = int(subnetDTO.LinkedSubnet)
	s.VLANID = int(subnetDTO.VLANID)
	s.VRFID = int(subnetDTO.VRFID)
	s.MasterSubnetID = int(subnetDTO.MasterSubnetID)
	s.NameserverID = int(subnetDTO.NameserverID)
	s.Nameservers = subnetDTO.Nameservers
	s.ShowName = subnetDTO.ShowName
	s.Permissions = subnetDTO.Permissions
	s.DNSRecursive = subnetDTO.DNSRecursive
	s.DNSRecords = subnetDTO.DNSRecords
	s.AllowRequests = subnetDTO.AllowRequests
	s.ScanAgent = int(subnetDTO.ScanAgent)
	s.PingSubnet = subnetDTO.PingSubnet
	s.DiscoverSubnet = subnetDTO.DiscoverSubnet
	s.IsFolder = subnetDTO.IsFolder
	s.IsPool = subnetDTO.IsPool
	s.IsFull = subnetDTO.IsFull
	s.Threshold = int(subnetDTO.Threshold)
	s.Location = int(subnetDTO.Location)
	s.EditDate = subnetDTO.EditDate
	s.Gateway = subnetDTO.Gateway
	s.GatewayID = subnetDTO.GatewayID
	s.CustomFields = subnetDTO.CustomFields
	s.ResolveDNS = subnetDTO.ResolveDNS
}

func (s *Subnet) ToDTO() *SubnetDTO {
	return &SubnetDTO{
		ID:             phpipam.JSONIntString(s.ID),
		SubnetAddress:  s.SubnetAddress,
		Mask:           s.Mask,
		Description:    s.Description,
		SectionID:      phpipam.JSONIntString(s.SectionID),
		LinkedSubnet:   phpipam.JSONIntString(s.LinkedSubnet),
		VLANID:         phpipam.JSONIntString(s.VLANID),
		VRFID:          phpipam.JSONIntString(s.VRFID),
		MasterSubnetID: phpipam.JSONIntString(s.MasterSubnetID),
		NameserverID:   phpipam.JSONIntString(s.NameserverID),
		Nameservers:    s.Nameservers,
		ShowName:       s.ShowName,
		Permissions:    s.Permissions,
		DNSRecursive:   s.DNSRecursive,
		DNSRecords:     s.DNSRecords,
		AllowRequests:  s.AllowRequests,
		ScanAgent:      phpipam.JSONIntString(s.ScanAgent),
		PingSubnet:     s.PingSubnet,
		DiscoverSubnet: s.DiscoverSubnet,
		IsFolder:       s.IsFolder,
		IsPool:         s.IsPool,
		IsFull:         s.IsFull,
		Threshold:      phpipam.JSONIntString(s.Threshold),
		Location:       phpipam.JSONIntString(s.Location),
		EditDate:       s.EditDate,
		Gateway:        s.Gateway,
		GatewayID:      s.GatewayID,
		CustomFields:   s.CustomFields,
		ResolveDNS:     s.ResolveDNS,
	}
}

// Controller is the base client for the Subnets controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Subnets controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// CreateSubnet creates a subnet by sending a POST request.
func (c *Controller) CreateSubnet(in Subnet) (message string, err error) {
	err = c.SendRequest("POST", "/subnets/", in.ToDTO(), &message)
	return
}

// CreateFirstFreeSubnet creates a first free child subnet inside subnet with specified mask by sending a POST request.
func (c *Controller) CreateFirstFreeSubnet(id int, mask int, in Subnet) (message string, err error) {
	err = c.SendRequest("POST", fmt.Sprintf("/subnets/%d/first_subnet/%d/", id, mask), in.ToDTO(), &message)
	return
}

// GetSubnetByID GETs a subnet via its ID.
func (c *Controller) GetSubnetByID(id int) (out Subnet, err error) {
	var dto SubnetDTO
	err = c.SendRequest("GET", fmt.Sprintf("/subnets/%d/", id), &struct{}{}, &dto)
	out.FromDTO(&dto)
	return
}

// GetSubnetsByCIDR GETs a subnet via its CIDR (i.e. 10.10.1.0/24).
//
// The function's name reflects the fact that an array of subnets is returned
// through the API, although it remains unclear how to actually query this
// method in a way that would return multiple results. Using a broader CIDR
// will not return multiple results, and using the CIDR of a master subnet will
// return that subnet only.
func (c *Controller) GetSubnetsByCIDR(cidr string) (out []Subnet, err error) {
	var dtos []SubnetDTO
	err = c.SendRequest("GET", fmt.Sprintf("/subnets/cidr/%s/", cidr), &struct{}{}, &dtos)
	for _, dto := range dtos {
		var subnet Subnet
		subnet.FromDTO(&dto)
		out = append(out, subnet)
	}
	return
}

func (c *Controller) GetSubnetsByCIDRAndSection(cidr string, section_id int) (out []Subnet, err error) {
	var dtos []SubnetDTO
	err = c.SendRequest("GET", fmt.Sprintf("/subnets/cidr/%s/?filter_by=sectionId&filter_value=%d", cidr, section_id), &struct{}{}, &dtos)
	for _, dto := range dtos {
		var subnet Subnet
		subnet.FromDTO(&dto)
		out = append(out, subnet)
	}
	return
}

// GetFirstFreeSubnet GETs the first free child subnet inside subnet with specified mask
func (c *Controller) GetFirstFreeSubnet(id int, mask int) (message string, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/subnets/%d/first_subnet/%d/", id, mask), &struct{}{}, &message)
	return
}

// GetFirstFreeAddress GETs the first free IP address in a subnet and returns
// it as a string. This can be used to automatically determine the next address
// you should use. If there are no more available addresses, the string will be
// blank.
//
// Note that marking a subnet as used does not prevent this function from
// returning data.
func (c *Controller) GetFirstFreeAddress(id int) (out string, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/subnets/%d/first_free/", id), &struct{}{}, &out)
	return
}

// GetAddressesInSubnet GETs the IP addresses for a specific subnet, via a
// supplied subnet ID.
func (c *Controller) GetAddressesInSubnet(id int) (out []addresses.Address, err error) {
	var dtos []addresses.AddressDTO
	err = c.SendRequest("GET", fmt.Sprintf("/subnets/%d/addresses/", id), &struct{}{}, &dtos)
	for _, dto := range dtos {
		var address addresses.Address
		address.FromDTO(&dto)
		out = append(out, address)
	}
	return
}

// GetSubnetCustomFieldsSchema GETs the custom fields for the subnets controller via
// client.GetCustomFieldsSchema.
func (c *Controller) GetSubnetCustomFieldsSchema() (out map[string]phpipam.CustomField, err error) {
	out, err = c.Client.GetCustomFieldsSchema("subnets")
	return
}

// GetSubnetCustomFields GETs the custom fields for a subnet via
// client.GetCustomFields.
func (c *Controller) GetSubnetCustomFields(id int) (out map[string]interface{}, err error) {
	out, err = c.Client.GetCustomFields(id, "subnets")
	return
}

// UpdateSubnet updates a subnet by sending a PATCH request.
//
// Note you cannot use this function to update a subnet's CIDR - to split,
// grow, or renumber a subnet, you need to use other methods that are currently
// not implemented in this SDK. See the API spec for more details.
func (c *Controller) UpdateSubnet(in Subnet) (message string, err error) {
	err = c.SendRequest("PATCH", "/subnets/", in.ToDTO(), &message)
	return
}

// UpdateSubnetCustomFields PATCHes the subnet's custom fields via
// client.UpdateCustomFields.
func (c *Controller) UpdateSubnetCustomFields(id int, in map[string]interface{}) (message string, err error) {
	message, err = c.Client.UpdateCustomFields(id, in, "subnets")
	return
}

// DeleteSubnet deletes a subnet by its ID.
func (c *Controller) DeleteSubnet(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/subnets/%d/", id), &struct{}{}, &message)
	return
}
