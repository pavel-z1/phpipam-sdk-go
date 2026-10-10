// Package nameservers provides types and methods for working with the
// nameservers controller (/tools/nameservers/).
package nameservers

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Nameserver represents a PHPIPAM nameserver set.
type Nameserver struct {
	// The nameserver set ID.
	ID int `json:"id,string,omitempty"`

	// The name of the nameserver set.
	Name string `json:"name,omitempty"`

	// A semicolon-separated list of nameserver addresses.
	Nameservers string `json:"namesrv1,omitempty"`

	// A detailed description of the nameserver set.
	Description string `json:"description,omitempty"`

	// A semicolon-separated list of section IDs the nameserver set is
	// available in.
	Sections string `json:"permissions,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`
}

// Controller is the base client for the Nameservers controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Nameservers controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListNameservers lists all nameserver sets.
func (c *Controller) ListNameservers() (out []Nameserver, err error) {
	err = c.SendRequest("GET", "/tools/nameservers/", &struct{}{}, &out)
	return
}

// CreateNameserver creates a nameserver set by sending a POST request.
func (c *Controller) CreateNameserver(in Nameserver) (message string, err error) {
	err = c.SendRequest("POST", "/tools/nameservers/", &in, &message)
	return
}

// GetNameserverByID GETs a nameserver set via its ID in the PHPIPAM database.
func (c *Controller) GetNameserverByID(id int) (out Nameserver, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/nameservers/%d/", id), &struct{}{}, &out)
	return
}

// UpdateNameserver updates a nameserver set by sending a PATCH request.
//
// The tools controller takes the object ID from the URL; an "id" field in the
// request body would override the subcontroller name, so it is not sent.
func (c *Controller) UpdateNameserver(in Nameserver) (message string, err error) {
	id := in.ID
	in.ID = 0
	err = c.SendRequest("PATCH", fmt.Sprintf("/tools/nameservers/%d/", id), &in, &message)
	return
}

// DeleteNameserver deletes a nameserver set by its ID.
func (c *Controller) DeleteNameserver(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/tools/nameservers/%d/", id), &struct{}{}, &message)
	return
}
