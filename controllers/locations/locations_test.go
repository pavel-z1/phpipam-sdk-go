package locations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/devices"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
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

const testLocationJSON = `{"id":"3","name":"DC1","description":"Primary datacenter","address":"1 Main St","lat":"42.69","long":"23.32"}`

var testLocationExpected = Location{
	ID:          3,
	Name:        "DC1",
	Description: "Primary datacenter",
	Address:     "1 Main St",
	Latitude:    "42.69",
	Longitude:   "23.32",
}

func TestListLocations(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testLocationJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).ListLocations()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/locations/")
	if expected := []Location{testLocationExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetLocationByID(t *testing.T) {
	for name, body := range map[string]string{
		"stringified": testLocationJSON,
		"native":      `{"id":3,"name":"DC1","description":"Primary datacenter","address":"1 Main St","lat":"42.69","long":"23.32"}`,
	} {
		t.Run(name, func(t *testing.T) {
			var last testRequest
			ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":`+body+`}`, &last)
			defer ts.Close()

			actual, err := NewController(fullSessionConfig(ts.URL)).GetLocationByID(testLocationExpected.ID)
			if err != nil {
				t.Fatalf("Bad: %s", err)
			}
			assertRequest(t, last, "GET", fmt.Sprintf("/0123456789abcdefgh/tools/locations/%d/", testLocationExpected.ID))
			if !reflect.DeepEqual(testLocationExpected, actual) {
				t.Fatalf("Expected %#v, got %#v", testLocationExpected, actual)
			}
		})
	}
}

func TestCreateLocation(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusCreated, `{"code":201,"success":true,"data":"locations object created"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).CreateLocation(Location{Name: "DC1", Address: "1 Main St"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "POST", "/0123456789abcdefgh/tools/locations/")
	if expected := map[string]interface{}{"name": "DC1", "address": "1 Main St"}; !reflect.DeepEqual(expected, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expected, last.Body)
	}
	if actual != "locations object created" {
		t.Fatalf("Expected %q, got %q", "locations object created", actual)
	}
}

func TestUpdateLocation(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"locations object updated"}`, &last)
	defer ts.Close()

	in := Location{ID: 3, Description: "updated"}
	if _, err := NewController(fullSessionConfig(ts.URL)).UpdateLocation(in); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "PATCH", fmt.Sprintf("/0123456789abcdefgh/tools/locations/%d/", in.ID))
	if _, ok := last.Body["id"]; ok {
		t.Fatalf("Expected no id in request body, got %#v", last.Body)
	}
	if len(last.Body) != 1 {
		t.Fatalf("Expected a single field in request body, got %#v", last.Body)
	}
}

func TestDeleteLocation(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"locations object deleted"}`, &last)
	defer ts.Close()

	if _, err := NewController(fullSessionConfig(ts.URL)).DeleteLocation(testLocationExpected.ID); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "DELETE", fmt.Sprintf("/0123456789abcdefgh/tools/locations/%d/", testLocationExpected.ID))
}

func TestLocationChildren(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[{"id":"9"}]}`, &last)
	defer ts.Close()
	c := NewController(fullSessionConfig(ts.URL))

	s, err := c.GetSubnetsInLocation(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/locations/3/subnets/")
	if expected := []subnets.Subnet{{ID: 9}}; !reflect.DeepEqual(expected, s) {
		t.Fatalf("Expected %#v, got %#v", expected, s)
	}

	d, err := c.GetDevicesInLocation(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/locations/3/devices/")
	if expected := []devices.Device{{ID: 9}}; !reflect.DeepEqual(expected, d) {
		t.Fatalf("Expected %#v, got %#v", expected, d)
	}

	a, err := c.GetAddressesInLocation(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/locations/3/ipaddresses/")
	if expected := []addresses.Address{{ID: 9}}; !reflect.DeepEqual(expected, a) {
		t.Fatalf("Expected %#v, got %#v", expected, a)
	}
}
