package devices

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
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

const testDeviceJSON = `{
  "id": "4",
  "hostname": "core-sw-01",
  "ip": "10.0.0.1",
  "type": "1",
  "description": "Core switch",
  "sections": "1",
  "rack": "2",
  "rack_start": "10",
  "rack_size": "1",
  "location": "3",
  "editDate": null
}`

var testDeviceExpected = Device{
	ID:          4,
	Hostname:    "core-sw-01",
	IPAddress:   "10.0.0.1",
	Type:        1,
	Description: "Core switch",
	Sections:    "1",
	Rack:        2,
	RackStart:   10,
	RackSize:    1,
	Location:    3,
}

func TestListDevices(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testDeviceJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).ListDevices()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/devices/")
	if expected := []Device{testDeviceExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetDeviceByID(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":`+testDeviceJSON+`}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetDeviceByID(4)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/devices/4/")
	if !reflect.DeepEqual(testDeviceExpected, actual) {
		t.Fatalf("Expected %#v, got %#v", testDeviceExpected, actual)
	}
}

func TestGetDeviceByIDNativeTypes(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":{"id":4,"hostname":"core-sw-01","ip":"10.0.0.1","type":1,"description":"Core switch","sections":"1","rack":2,"rack_start":10,"rack_size":1,"location":3,"editDate":null}}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetDeviceByID(4)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if !reflect.DeepEqual(testDeviceExpected, actual) {
		t.Fatalf("Expected %#v, got %#v", testDeviceExpected, actual)
	}
}

func TestSearchDevices(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testDeviceJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).SearchDevices("core-sw")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/devices/search/core-sw/")
	if expected := []Device{testDeviceExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetAddressesOnDevice(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[{"id":"11","subnetId":"3","ip":"10.10.1.10","deviceId":"4"}]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetAddressesOnDevice(4)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/devices/4/addresses/")
	expected := []addresses.Address{{ID: 11, SubnetID: 3, IPAddress: "10.10.1.10", DeviceID: 4}}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSubnetsOnDevice(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[{"id":"3","subnet":"10.10.1.0","mask":"24"}]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetSubnetsOnDevice(4)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/devices/4/subnets/")
	expected := []subnets.Subnet{{ID: 3, SubnetAddress: "10.10.1.0", Mask: 24}}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestCreateDevice(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusCreated, `{"code":201,"success":true,"data":"Device created"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).CreateDevice(Device{Hostname: "core-sw-01", IPAddress: "10.0.0.1", Type: 1})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "POST", "/0123456789abcdefgh/devices/")
	expectedBody := map[string]interface{}{"hostname": "core-sw-01", "ip": "10.0.0.1", "type": "1"}
	if !reflect.DeepEqual(expectedBody, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expectedBody, last.Body)
	}
	if actual != "Device created" {
		t.Fatalf("Expected %q, got %q", "Device created", actual)
	}
}

func TestUpdateDevice(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":"Device updated"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).UpdateDevice(Device{ID: 4, Description: "updated"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "PATCH", "/0123456789abcdefgh/devices/")
	expectedBody := map[string]interface{}{"id": "4", "description": "updated"}
	if !reflect.DeepEqual(expectedBody, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expectedBody, last.Body)
	}
	if actual != "Device updated" {
		t.Fatalf("Expected %q, got %q", "Device updated", actual)
	}
}

func TestDeleteDevice(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":"Device deleted"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).DeleteDevice(4)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "DELETE", "/0123456789abcdefgh/devices/4/")
	if actual != "Device deleted" {
		t.Fatalf("Expected %q, got %q", "Device deleted", actual)
	}
}
