package tags

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

const testTagJSON = `{"id":"2","type":"Used","showtag":"0","bgcolor":"#a9c9a4","fgcolor":"#ffffff","compress":"No","locked":"Yes","updateTag":"1"}`

var testTagExpected = Tag{
	ID:        2,
	Type:      "Used",
	ShowTag:   false,
	BGColor:   "#a9c9a4",
	FGColor:   "#ffffff",
	Compress:  "No",
	Locked:    "Yes",
	UpdateTag: true,
}

func TestListTags(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":[`+testTagJSON+`]}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).ListTags()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "GET", "/0123456789abcdefgh/tools/tags/")
	if expected := []Tag{testTagExpected}; !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetTagByID(t *testing.T) {
	for name, body := range map[string]string{
		"stringified": testTagJSON,
		"native":      `{"id":2,"type":"Used","showtag":0,"bgcolor":"#a9c9a4","fgcolor":"#ffffff","compress":"No","locked":"Yes","updateTag":1}`,
	} {
		t.Run(name, func(t *testing.T) {
			var last testRequest
			ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"data":`+body+`}`, &last)
			defer ts.Close()

			actual, err := NewController(fullSessionConfig(ts.URL)).GetTagByID(testTagExpected.ID)
			if err != nil {
				t.Fatalf("Bad: %s", err)
			}
			assertRequest(t, last, "GET", fmt.Sprintf("/0123456789abcdefgh/tools/tags/%d/", testTagExpected.ID))
			if !reflect.DeepEqual(testTagExpected, actual) {
				t.Fatalf("Expected %#v, got %#v", testTagExpected, actual)
			}
		})
	}
}

func TestCreateTag(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusCreated, `{"code":201,"success":true,"data":"tags object created"}`, &last)
	defer ts.Close()

	actual, err := NewController(fullSessionConfig(ts.URL)).CreateTag(Tag{Type: "Decommissioned", BGColor: "#888888"})
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "POST", "/0123456789abcdefgh/tools/tags/")
	if expected := map[string]interface{}{"type": "Decommissioned", "bgcolor": "#888888"}; !reflect.DeepEqual(expected, last.Body) {
		t.Fatalf("Expected body %#v, got %#v", expected, last.Body)
	}
	if actual != "tags object created" {
		t.Fatalf("Expected %q, got %q", "tags object created", actual)
	}
}

func TestUpdateTag(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"tags object updated"}`, &last)
	defer ts.Close()

	in := Tag{ID: 2, FGColor: "#000000"}
	if _, err := NewController(fullSessionConfig(ts.URL)).UpdateTag(in); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "PATCH", fmt.Sprintf("/0123456789abcdefgh/tools/tags/%d/", in.ID))
	if _, ok := last.Body["id"]; ok {
		t.Fatalf("Expected no id in request body, got %#v", last.Body)
	}
	if len(last.Body) != 1 {
		t.Fatalf("Expected a single field in request body, got %#v", last.Body)
	}
}

func TestDeleteTag(t *testing.T) {
	var last testRequest
	ts := newTestServer(http.StatusOK, `{"code":200,"success":true,"message":"tags object deleted"}`, &last)
	defer ts.Close()

	if _, err := NewController(fullSessionConfig(ts.URL)).DeleteTag(testTagExpected.ID); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	assertRequest(t, last, "DELETE", fmt.Sprintf("/0123456789abcdefgh/tools/tags/%d/", testTagExpected.ID))
}
