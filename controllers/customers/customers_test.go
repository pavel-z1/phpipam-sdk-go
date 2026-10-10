package customers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// testRequest records the details of the last request received by a test
// server.
type testRequest struct {
	Method string
	URI    string
	Body   map[string]interface{}
}

func newTestServer(status int, output string, last *testRequest) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last.Method = r.Method
		last.URI = r.URL.RequestURI()
		last.Body = nil
		if b, _ := io.ReadAll(r.Body); len(b) > 0 {
			json.Unmarshal(b, &last.Body)
		}
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, output, status)
	}))
}

func fullSessionConfig(endpoint string) *session.Session {
	return &session.Session{
		Config: phpipam.Config{
			AppID:    "0123456789abcdefgh",
			Endpoint: endpoint,
			Password: "changeit",
			Username: "nobody",
		},
		Token: session.Token{
			String: "foobarbazboop",
		},
	}
}

func assertRequest(t *testing.T, last testRequest, method, uri string) {
	t.Helper()
	if last.Method != method || last.URI != uri {
		t.Fatalf("Expected request %s %s, got %s %s", method, uri, last.Method, last.URI)
	}
}

const testCustomerJSON = `{"id":"5","title":"ACME","address":"1 Main St","postcode":"1000","city":"Sofia","tag":"Sofia-city","lat":null,"long":null,"contact_person":"Jane Doe","contact_phone":"+359000000","contact_mail":"jane@example.com","note":null,"status":"Active"}`

var testCustomerExpected = Customer{
	ID:            5,
	Title:         "ACME",
	Address:       "1 Main St",
	Postcode:      "1000",
	City:          "Sofia",
	State:         "Sofia-city",
	ContactPerson: "Jane Doe",
	ContactPhone:  "+359000000",
	ContactMail:   "jane@example.com",
	Status:        "Active",
}

func TestListCustomers(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testCustomerJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).ListCustomers()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/customers/")
	if expected := []Customer{testCustomerExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetCustomerByID(t *testing.T) {
	for name, body := range map[string]string{
		"stringified": testCustomerJSON,
		"native":      `{"id":5,"title":"ACME","address":"1 Main St","postcode":"1000","city":"Sofia","tag":"Sofia-city","lat":null,"long":null,"contact_person":"Jane Doe","contact_phone":"+359000000","contact_mail":"jane@example.com","note":null,"status":"Active"}`,
	} {
		t.Run(name, func(t *testing.T) {
			var last testRequest
			ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":`+body+`}`, &last)
			defer ts.Close()

			actual, err := NewController(fullSessionConfig(ts.URL)).GetCustomerByID(testCustomerExpected.ID)
			if err != nil {
				t.Fatalf("Bad: %s", err)
			}
			assertRequest(t, last, "GET", fmt.Sprintf("/0123456789abcdefgh/tools/customers/%d/", testCustomerExpected.ID))
			if !reflect.DeepEqual(testCustomerExpected, actual) {
				t.Fatalf("Expected %#v, got %#v", testCustomerExpected, actual)
			}
		})
	}
}

func TestCreateCustomer(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusCreated, `{"code":201,"success":true,"data":"customers object created"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).CreateCustomer(Customer{Title: "ACME", State: "Sofia-city"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "POST", "/0123456789abcdefgh/tools/customers/")
	if expected := map[string]interface{}{"title": "ACME", "tag": "Sofia-city"}; !reflect.DeepEqual(expected, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expected, last.Body)
	}
	if actual != "customers object created" {
		t.Fatalf("Expected %q, got %q", "customers object created", actual)
	}
}

func TestUpdateCustomer(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"customers object updated"}`, &last)
	defer ts.Close()

	in := Customer{ID: 5, Status: "Inactive"}
	if _, err := NewController(fullSessionConfig(ts.URL)).UpdateCustomer(in); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "PATCH", fmt.Sprintf("/0123456789abcdefgh/tools/customers/%d/", in.ID))
	if _, ok := last.Body["id"]; ok {
		t.Fatalf("Expected no id in request body, got %#v", last.Body)
	}
	if len(last.Body) != 1 {
		t.Fatalf("Expected a single field in request body, got %#v", last.Body)
	}
}

func TestDeleteCustomer(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"customers object deleted"}`, &last)
	defer ts.Close()

	if _, err := NewController(fullSessionConfig(ts.URL)).DeleteCustomer(testCustomerExpected.ID); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "DELETE", fmt.Sprintf("/0123456789abcdefgh/tools/customers/%d/", testCustomerExpected.ID))
}
