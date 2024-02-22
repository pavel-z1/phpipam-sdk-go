// Package sections provides types and methods for working with the sections
// controller.
package sections

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Section represents a PHPIPAM section.
type SectionDTO struct {
	// The section ID.
	ID phpipam.JSONIntString `json:"id,omitempty"`

	// The section's name.
	Name string `json:"name,omitempty"`

	// The section's description.
	Description string `json:"description,omitempty"`

	// The ID of the section's parent, if nested.
	MasterSection phpipam.JSONIntString `json:"masterSection,omitempty"`

	// A JSON object, stringified, that represents the permissions for this
	// section.
	Permissions string `json:"permissions,omitempty"`

	// Whether or not to check consistency for subnets and IP addresses.
	StrictMode phpipam.BoolIntString `json:"strictMode,omitempty"`

	// How to order subnets in this section when viewing.
	SubnetOrdering string `json:"subnetOrdering,omitempty"`

	// The order position of this section when displaying sections.
	Order phpipam.JSONIntString `json:"order,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// Whether or not to show VLANs in the subnet listing of this section.
	ShowVLAN phpipam.BoolIntString `json:"showVLAN,omitempty"`

	// Whether or not to show VRF information in the subnet listing of this
	// section.
	ShowVRF phpipam.BoolIntString `json:"showVRF,omitempty"`

	// Whether or not to show only supernets in the subnet listing of this
	// section.
	ShowSupernetOnly phpipam.BoolIntString `json:"showSupernetOnly,omitempty"`

	// The ID of the DNS resolver to be used for this section.
	DNS phpipam.JSONIntString `json:"DNS,omitempty"`
}

type Section struct {
	// The section ID.
	ID int

	// The section's name.
	Name string

	// The section's description.
	Description string

	// The ID of the section's parent, if nested.
	MasterSection int

	// A JSON object, stringified, that represents the permissions for this
	// section.
	Permissions string

	// Whether or not to check consistency for subnets and IP addresses.
	StrictMode phpipam.BoolIntString

	// How to order subnets in this section when viewing.
	SubnetOrdering string

	// The order position of this section when displaying sections.
	Order int

	// The date of the last edit to this resource.
	EditDate string

	// Whether or not to show VLANs in the subnet listing of this section.
	ShowVLAN phpipam.BoolIntString

	// Whether or not to show VRF information in the subnet listing of this
	// section.
	ShowVRF phpipam.BoolIntString

	// Whether or not to show only supernets in the subnet listing of this
	// section.
	ShowSupernetOnly phpipam.BoolIntString

	// The ID of the DNS resolver to be used for this section.
	DNS int
}

func (s *Section) FromDTO(sectionDTO *SectionDTO) {
	s.ID = int(sectionDTO.ID)
	s.Name = sectionDTO.Name
	s.Description = sectionDTO.Description
	s.MasterSection = int(sectionDTO.MasterSection)
	s.Permissions = sectionDTO.Permissions
	s.StrictMode = sectionDTO.StrictMode
	s.SubnetOrdering = sectionDTO.SubnetOrdering
	s.Order = int(sectionDTO.Order)
	s.EditDate = sectionDTO.EditDate
	s.ShowVLAN = sectionDTO.ShowVLAN
	s.ShowVRF = sectionDTO.ShowVRF
	s.ShowSupernetOnly = sectionDTO.ShowSupernetOnly
	s.DNS = int(sectionDTO.DNS)
}

func (s *Section) ToDTO() *SectionDTO {
	return &SectionDTO{
		ID:               phpipam.JSONIntString(s.ID),
		Name:             s.Name,
		Description:      s.Description,
		MasterSection:    phpipam.JSONIntString(s.MasterSection),
		Permissions:      s.Permissions,
		StrictMode:       s.StrictMode,
		SubnetOrdering:   s.SubnetOrdering,
		Order:            phpipam.JSONIntString(s.Order),
		EditDate:         s.EditDate,
		ShowVLAN:         s.ShowVLAN,
		ShowVRF:          s.ShowVRF,
		ShowSupernetOnly: s.ShowSupernetOnly,
		DNS:              phpipam.JSONIntString(s.DNS),
	}
}

// Controller is the base client for the Sections controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Sections controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListSections lists all sections.
func (c *Controller) ListSections() (out []Section, err error) {
	var dto []SectionDTO
	err = c.SendRequest("GET", "/sections/", &struct{}{}, &dto)
	for _, sectionDTO := range dto {
		var section Section
		section.FromDTO(&sectionDTO)
		out = append(out, section)
	}
	return
}

// CreateSection creates a section by sending a POST request.
func (c *Controller) CreateSection(in Section) (message string, err error) {
	err = c.SendRequest("POST", "/sections/", in.ToDTO(), &message)
	return
}

// GetSectionByID GETs a section via its ID.
func (c *Controller) GetSectionByID(id int) (out Section, err error) {
	var dto SectionDTO
	err = c.SendRequest("GET", fmt.Sprintf("/sections/%d/", id), &struct{}{}, &dto)
	out.FromDTO(&dto)
	return
}

// GetSectionByName GETs a section via its name.
func (c *Controller) GetSectionByName(name string) (out Section, err error) {
	var dto SectionDTO
	err = c.SendRequest("GET", fmt.Sprintf("/sections/%s/", name), &struct{}{}, &dto)
	out.FromDTO(&dto)
	return
}

// GetSubnetsInSection GETs the subnets in a section by section ID.
func (c *Controller) GetSubnetsInSection(id int) (out []subnets.Subnet, err error) {
	var dtos []subnets.SubnetDTO
	err = c.SendRequest("GET", fmt.Sprintf("/sections/%d/subnets/", id), &struct{}{}, &dtos)
	for _, subnetDTO := range dtos {
		var subnet subnets.Subnet
		subnet.FromDTO(&subnetDTO)
		out = append(out, subnet)
	}
	return
}

// UpdateSection updates a section by sending a PATCH request.
func (c *Controller) UpdateSection(in Section) (err error) {
	err = c.SendRequest("PATCH", "/sections/", in.ToDTO(), &struct{}{})
	return
}

// DeleteSection deletes a section by sending a DELETE request. All subnets and
// addresses in the section will be deleted as well.
func (c *Controller) DeleteSection(id int) (err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/sections/%d/", id), &struct{}{}, &struct{}{})
	return
}
