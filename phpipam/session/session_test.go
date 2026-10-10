package session

import (
	"reflect"
	"testing"
	"time"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

func phpipamConfig() phpipam.Config {
	return phpipam.Config{
		AppID:    "0123456789abcdefgh",
		Endpoint: "http://localhost/api",
		Password: "changeit",
		Username: "nobody",
	}
}

func fullSessionConfig() *Session {
	return &Session{
		Config: phpipamConfig(),
		Token: Token{
			String: "foobarbazboop",
		},
	}
}

func TestNewSession(t *testing.T) {
	cfg := phpipamConfig()

	expected := &Session{
		Config: phpipamConfig(),
	}

	actual := NewSession(cfg)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected session to be %#v, got %#v", expected, actual)
	}
}

func TestNewSessionMergesConfigs(t *testing.T) {
	override := phpipam.Config{
		Username: "somebody",
		Insecure: true,
		Timeout:  5 * time.Second,
	}
	expected := phpipamConfig()
	expected.Username = "somebody"
	expected.Insecure = true
	expected.Timeout = 5 * time.Second

	actual := NewSession(phpipamConfig(), override)

	if !reflect.DeepEqual(expected, actual.Config) {
		t.Fatalf("Expected config to be %#v, got %#v", expected, actual.Config)
	}
}

func TestHTTPClient(t *testing.T) {
	s := fullSessionConfig()
	c := s.HTTPClient()
	if c != s.HTTPClient() {
		t.Fatalf("Expected HTTPClient to return the same client on each call")
	}
	if c.Timeout != phpipam.DefaultTimeout {
		t.Fatalf("Expected default timeout %s, got %s", phpipam.DefaultTimeout, c.Timeout)
	}

	s = fullSessionConfig()
	s.Config.Timeout = 5 * time.Second
	if actual := s.HTTPClient().Timeout; actual != 5*time.Second {
		t.Fatalf("Expected timeout 5s, got %s", actual)
	}
}
