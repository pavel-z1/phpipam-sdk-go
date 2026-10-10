// Package customers provides types and methods for working with the
// customers controller (/tools/customers/).
package customers

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Customer represents a PHPIPAM customer.
type Customer struct {
	// The customer ID.
	ID int `json:"id,string,omitempty"`

	// The customer name.
	Title string `json:"title,omitempty"`

	// The street address of the customer.
	Address string `json:"address,omitempty"`

	// The postal code of the customer.
	Postcode string `json:"postcode,omitempty"`

	// The city of the customer.
	City string `json:"city,omitempty"`

	// The state or region of the customer. PHPIPAM stores this in the "state"
	// column but exposes it through the API as "tag".
	State string `json:"tag,omitempty"`

	// The latitude of the customer's address.
	Latitude string `json:"lat,omitempty"`

	// The longitude of the customer's address.
	Longitude string `json:"long,omitempty"`

	// The name of the contact person.
	ContactPerson string `json:"contact_person,omitempty"`

	// The phone number of the contact person.
	ContactPhone string `json:"contact_phone,omitempty"`

	// The email address of the contact person.
	ContactMail string `json:"contact_mail,omitempty"`

	// Free-form notes about the customer.
	Note string `json:"note,omitempty"`

	// The customer status. One of "Active", "Reserved" or "Inactive".
	Status string `json:"status,omitempty"`
}

// Controller is the base client for the Customers controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Customers controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListCustomers lists all customers.
func (c *Controller) ListCustomers() (out []Customer, err error) {
	err = c.SendRequest("GET", "/tools/customers/", &struct{}{}, &out)
	return
}

// CreateCustomer creates a customer by sending a POST request.
func (c *Controller) CreateCustomer(in Customer) (message string, err error) {
	err = c.SendRequest("POST", "/tools/customers/", &in, &message)
	return
}

// GetCustomerByID GETs a customer via its ID in the PHPIPAM database.
func (c *Controller) GetCustomerByID(id int) (out Customer, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/tools/customers/%d/", id), &struct{}{}, &out)
	return
}

// UpdateCustomer updates a customer by sending a PATCH request.
//
// The tools controller takes the object ID from the URL; an "id" field in the
// request body would override the subcontroller name, so it is not sent.
func (c *Controller) UpdateCustomer(in Customer) (message string, err error) {
	id := in.ID
	in.ID = 0
	err = c.SendRequest("PATCH", fmt.Sprintf("/tools/customers/%d/", id), &in, &message)
	return
}

// DeleteCustomer deletes a customer by its ID.
func (c *Controller) DeleteCustomer(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/tools/customers/%d/", id), &struct{}{}, &message)
	return
}
