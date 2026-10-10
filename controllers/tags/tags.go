// Package tags provides types and methods for working with the
// IP address tags controller (/tools/tags/).
package tags

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Tag represents a PHPIPAM IP address tag.
type Tag struct {
	// The tag ID. This is the value used in the Tag field of an address.
	ID int `json:"id,string,omitempty"`

	// The tag name, such as "Used", "Reserved" or "Offline".
	Type string `json:"type,omitempty"`

	// Whether the tag is shown in the address list.
	ShowTag phpipam.BoolIntString `json:"showtag,omitempty"`

	// The background color of the tag, as a hex color code.
	BGColor string `json:"bgcolor,omitempty"`

	// The foreground color of the tag, as a hex color code.
	FGColor string `json:"fgcolor,omitempty"`

	// Whether addresses with this tag are compressed into ranges in the
	// address list. One of "Yes" or "No".
	Compress string `json:"compress,omitempty"`

	// Whether the tag is locked (built-in) and cannot be deleted. One of "Yes"
	// or "No".
	Locked string `json:"locked,omitempty"`

	// Whether the tag of an address is updated by status scans.
	UpdateTag phpipam.BoolIntString `json:"updateTag,omitempty"`
}

// Controller is the base client for the Tags controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Tags controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListTags lists all IP address tags.
func (c *Controller) ListTags() (out []Tag, err error) {
	err = c.SendRequest("GET", "/tools/tags/", &struct{}{}, &out)
	return
}

// CreateTag creates an IP address tag by sending a POST request.
func (c *Controller) CreateTag(in Tag) (message string, err error) {
	err = c.SendRequest("POST", "/tools/tags/", &in, &message)
	return
}

// GetTagByID GETs an IP address tag via its ID in the PHPIPAM database.
func (c *Controller) GetTagByID(id int) (out Tag, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/tags/%d/", id), &struct{}{}, &out)
	return
}

// UpdateTag updates an IP address tag by sending a PATCH request.
//
// The tools controller takes the object ID from the URL; an "id" field in the
// request body would override the subcontroller name, so it is not sent.
func (c *Controller) UpdateTag(in Tag) (message string, err error) {
	id := in.ID
	in.ID = 0
	err = c.SendRequest("PATCH", fmt.Sprintf("/tools/tags/%d/", id), &in, &message)
	return
}

// DeleteTag deletes an IP address tag by its ID.
func (c *Controller) DeleteTag(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/tools/tags/%d/", id), &struct{}{}, &message)
	return
}
