package search

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

const testSearchSubnetsOnlyOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "subnets": {
      "code": 200,
      "data": [
        {
          "id": "3",
          "subnet": "168099840",
          "mask": "24",
		  "description": "Customer 1",
          "sectionId": "1",
          "masterSubnetId": "2",
          "allowRequests": "1",
          "vlanId": "0",
          "showName": "1",
          "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
          "pingSubnet": "0",
          "discoverSubnet": "0",
          "isFolder": "0",
          "isFull": "0"
        }
      ]
    }
  }
}
`

var testSearchSubnetsOnlyOutputExpected = []subnets.Subnet{
	{
		ID:             3,
		SubnetAddress:  "10.5.0.0",
		Mask:           24,
		SectionID:      1,
		Description:    "Customer 1",
		MasterSubnetID: 2,
		AllowRequests:  true,
		ShowName:       true,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
	},
}

func newHTTPTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(f))
	return ts
}

func httpOKTestServer(output string) *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, output, http.StatusOK)
	})
}

func fullSessionConfig() *session.Session {
	return &session.Session{
		Config: phpipam.Config{
			AppID:    "0123456789abcdefgh",
			Password: "changeit",
			Username: "nobody",
		},
		Token: session.Token{
			String: "foobarbazboop",
		},
	}
}

func TestSearch(t *testing.T) {
	ts := httpOKTestServer(testSearchSubnetsOnlyOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testSearchSubnetsOnlyOutputExpected
	actual, err := client.SearchSubnets("Customer 1")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestIntegerToIP(t *testing.T) {
	cases := map[string]string{
		"0":                                      "0.0.0.0",
		"168099840":                              "10.5.0.0",
		"4294967295":                             "255.255.255.255",
		"42540766411282592856903984951653826560": "2001:db8::",
		"340282366920938463463374607431768211455": "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"10.5.0.0":   "10.5.0.0",
		"2001:db8::": "2001:db8::",
	}
	for in, expected := range cases {
		actual, err := integerToIP(in)
		if err != nil {
			t.Fatalf("%s: unexpected error: %s", in, err)
		}
		if actual != expected {
			t.Fatalf("%s: expected %s, got %s", in, expected, actual)
		}
	}

	for _, in := range []string{"", "abc", "-1", "340282366920938463463374607431768211456"} {
		if _, err := integerToIP(in); err == nil {
			t.Fatalf("%q: expected error, got none", in)
		}
	}
}

func TestSearchNoResults(t *testing.T) {
	ts := httpOKTestServer(`{"code":200,"success":true,"data":{"subnets":{"code":404,"data":"No subnets found"}}}`)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	actual, err := client.SearchSubnets("nothing")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
	if actual == nil || len(actual) != 0 {
		t.Fatalf("Expected empty result, got %#v", actual)
	}
}

func TestSearchEscapesSearchString(t *testing.T) {
	var uri string
	ts := newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		uri = r.RequestURI
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, testSearchSubnetsOnlyOutputJSON, http.StatusOK)
	})
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	if _, err := client.SearchSubnets("10.5.0.0/24 #1?"); err != nil {
		t.Fatalf("Bad: %s", err)
	}
	expected := "/0123456789abcdefgh/search/10.5.0.0%2F24%20%231%3F/?addresses=0&subnets=1&vlans=0&vrfs=0"
	if uri != expected {
		t.Fatalf("Expected request to %s, got %s", expected, uri)
	}
}
