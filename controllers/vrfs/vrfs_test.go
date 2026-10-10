package vrfs

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

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

const testVRFJSON = `{
  "id": "2",
  "vrfId": "2",
  "name": "customer-a",
  "rd": "65000:100",
  "description": "Customer A",
  "sections": "1;2",
  "editDate": null
}`

var testVRFExpected = VRF{
	ID:          2,
	Name:        "customer-a",
	RD:          "65000:100",
	Description: "Customer A",
	Sections:    "1;2",
}

func TestListVRFs(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testVRFJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).ListVRFs()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/vrfs/")
	if expected := []VRF{testVRFExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetVRFByID(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":`+testVRFJSON+`}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetVRFByID(2)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/vrfs/2/")
	if !reflect.DeepEqual(testVRFExpected, actual) {
		t.Fatalf("Expected %#v, got %#v", testVRFExpected, actual)
	}
}

func TestGetVRFByIDNativeTypes(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":{"id":2,"vrfId":2,"name":"customer-a","rd":"65000:100","description":"Customer A","sections":"1;2","editDate":null}}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetVRFByID(2)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if !reflect.DeepEqual(testVRFExpected, actual) {
		t.Fatalf("Expected %#v, got %#v", testVRFExpected, actual)
	}
}

func TestGetVRFsByName(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testVRFJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetVRFsByName("customer a")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/vrfs/?filter_by=name&filter_value=customer+a")
	if expected := []VRF{testVRFExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSubnetsInVRF(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[{"id":"3","subnet":"10.10.1.0","mask":"24","vrfId":"2"}]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).GetSubnetsInVRF(2)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/vrfs/2/subnets/")
	expected := []subnets.Subnet{{ID: 3, SubnetAddress: "10.10.1.0", Mask: 24, VRFID: 2}}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestCreateVRF(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusCreated, `{"code":201,"success":true,"data":"VRF created"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).CreateVRF(VRF{Name: "customer-a", RD: "65000:100"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "POST", "/0123456789abcdefgh/vrfs/")
	expectedBody := map[string]interface{}{"name": "customer-a", "rd": "65000:100"}
	if !reflect.DeepEqual(expectedBody, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expectedBody, last.Body)
	}
	if actual != "VRF created" {
		t.Fatalf("Expected %q, got %q", "VRF created", actual)
	}
}

func TestUpdateVRF(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":"VRF updated"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).UpdateVRF(VRF{ID: 2, Description: "updated"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "PATCH", "/0123456789abcdefgh/vrfs/")
	expectedBody := map[string]interface{}{"id": "2", "description": "updated"}
	if !reflect.DeepEqual(expectedBody, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expectedBody, last.Body)
	}
	if actual != "VRF updated" {
		t.Fatalf("Expected %q, got %q", "VRF updated", actual)
	}
}

func TestDeleteVRF(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":"VRF deleted"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).DeleteVRF(2)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "DELETE", "/0123456789abcdefgh/vrfs/2/")
	if actual != "VRF deleted" {
		t.Fatalf("Expected %q, got %q", "VRF deleted", actual)
	}
}
