package nameservers

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

const testNameserverJSON = `{"id":"1","name":"Google NS","namesrv1":"8.8.8.8;8.8.4.4","description":"Google public nameservers","permissions":"1;2","editDate":null}`

var testNameserverExpected = Nameserver{
	ID:          1,
	Name:        "Google NS",
	Nameservers: "8.8.8.8;8.8.4.4",
	Description: "Google public nameservers",
	Sections:    "1;2",
}

func TestListNameservers(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testNameserverJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).ListNameservers()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/nameservers/")
	if expected := []Nameserver{testNameserverExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetNameserverByID(t *testing.T) {
	for name, body := range map[string]string{
		"stringified": testNameserverJSON,
		"native":      `{"id":1,"name":"Google NS","namesrv1":"8.8.8.8;8.8.4.4","description":"Google public nameservers","permissions":"1;2","editDate":null}`,
	} {
		t.Run(name, func(t *testing.T) {
			var last testRequest
			ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":`+body+`}`, &last)
			defer ts.Close()

			actual, err := NewController(fullSessionConfig(ts.URL)).GetNameserverByID(testNameserverExpected.ID)
			if err != nil {
				t.Fatalf("Bad: %s", err)
			}
			assertRequest(t, last, "GET", fmt.Sprintf("/0123456789abcdefgh/tools/nameservers/%d/", testNameserverExpected.ID))
			if !reflect.DeepEqual(testNameserverExpected, actual) {
				t.Fatalf("Expected %#v, got %#v", testNameserverExpected, actual)
			}
		})
	}
}

func TestCreateNameserver(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusCreated, `{"code":201,"success":true,"data":"nameservers object created"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).CreateNameserver(Nameserver{Name: "Google NS", Nameservers: "8.8.8.8;8.8.4.4"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "POST", "/0123456789abcdefgh/tools/nameservers/")
	if expected := map[string]interface{}{"name": "Google NS", "namesrv1": "8.8.8.8;8.8.4.4"}; !reflect.DeepEqual(expected, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expected, last.Body)
	}
	if actual != "nameservers object created" {
		t.Fatalf("Expected %q, got %q", "nameservers object created", actual)
	}
}

func TestUpdateNameserver(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"nameservers object updated"}`, &last)
	defer ts.Close()

	in := Nameserver{ID: 1, Description: "updated"}
	if _, err := NewController(fullSessionConfig(ts.URL)).UpdateNameserver(in); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "PATCH", fmt.Sprintf("/0123456789abcdefgh/tools/nameservers/%d/", in.ID))
	if _, ok := last.Body["id"]; ok {
		t.Fatalf("Expected no id in request body, got %#v", last.Body)
	}
	if len(last.Body) != 1 {
		t.Fatalf("Expected a single field in request body, got %#v", last.Body)
	}
}

func TestDeleteNameserver(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"nameservers object deleted"}`, &last)
	defer ts.Close()

	if _, err := NewController(fullSessionConfig(ts.URL)).DeleteNameserver(testNameserverExpected.ID); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "DELETE", fmt.Sprintf("/0123456789abcdefgh/tools/nameservers/%d/", testNameserverExpected.ID))
}
