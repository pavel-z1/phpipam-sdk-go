// Package vrfs provides types and methods for working with the VRF
// controller.
package vrfs

import (
	"fmt"
	"net/url"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// VRF represents a PHPIPAM VRF.
type VRF struct {
	// The VRF ID. This is the vrfId column in the PHPIPAM database.
	ID int `json:"id,string,omitempty"`

	// The VRF name.
	Name string `json:"name,omitempty"`

	// The route distinguisher of the VRF.
	RD string `json:"rd,omitempty"`

	// A detailed description of the VRF.
	Description string `json:"description,omitempty"`

	// A semicolon-separated list of section IDs this VRF belongs to.
	Sections string `json:"sections,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration. If this is not
	// enabled, this map will be nil on GETs and POSTs and PATCHes with this
	// field set will fail. Use the explicit custom field functions instead.
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// Controller is the base client for the VRF controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the VRF controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListVRFs lists all VRFs.
func (c *Controller) ListVRFs() (out []VRF, err error) {
	err = c.SendRequest("GET", "/vrfs/", &struct{}{}, &out)
	return
}

// CreateVRF creates a VRF by sending a POST request.
func (c *Controller) CreateVRF(in VRF) (message string, err error) {
	err = c.SendRequest("POST", "/vrfs/", &in, &message)
	return
}

// GetVRFByID GETs a VRF via its ID in the PHPIPAM database.
func (c *Controller) GetVRFByID(id int) (out VRF, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/vrfs/%d/", id), &struct{}{}, &out)
	return
}

// GetVRFsByName GETs all VRFs matching the supplied name.
//
// PHPIPAM does not provide a search method on the VRF controller, so this
// filters the full VRF list using the filter_by API parameter.
func (c *Controller) GetVRFsByName(name string) (out []VRF, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/vrfs/?filter_by=name&filter_value=%s", url.QueryEscape(name)), &struct{}{}, &out)
	return
}

// GetSubnetsInVRF GETs the subnets assigned to a VRF.
func (c *Controller) GetSubnetsInVRF(id int) (out []subnets.Subnet, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/vrfs/%d/subnets/", id), &struct{}{}, &out)
	return
}

// GetVRFCustomFieldsSchema GETs the custom fields for the vrfs controller via
// client.GetCustomFieldsSchema.
func (c *Controller) GetVRFCustomFieldsSchema() (out map[string]phpipam.CustomField, err error) {
	out, err = c.Client.GetCustomFieldsSchema("vrfs")
	return
}

// GetVRFCustomFields GETs the custom fields for a VRF via
// client.GetCustomFields.
func (c *Controller) GetVRFCustomFields(id int) (out map[string]interface{}, err error) {
	out, err = c.Client.GetCustomFields(id, "vrfs")
	return
}

// UpdateVRF updates a VRF by sending a PATCH request.
func (c *Controller) UpdateVRF(in VRF) (message string, err error) {
	err = c.SendRequest("PATCH", "/vrfs/", &in, &message)
	return
}

// UpdateVRFCustomFields PATCHes the VRF's custom fields via
// client.UpdateCustomFields.
func (c *Controller) UpdateVRFCustomFields(id int, in map[string]interface{}) (message string, err error) {
	message, err = c.Client.UpdateCustomFields(id, in, "vrfs")
	return
}

// DeleteVRF deletes a VRF by its ID.
func (c *Controller) DeleteVRF(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/vrfs/%d/", id), &struct{}{}, &message)
	return
}
