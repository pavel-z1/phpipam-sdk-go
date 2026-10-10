// Package session provides session management utility and token storage.
package session

import (
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

// timeLayout represents the datetime format returned by the PHPIPAM api.
const timeLayout = "2006-01-02 15:04:05"

// Token represents a PHPIPAM session token.
type Token struct {
	// The token string.
	String string `json:"token"`
}

// Session represents a PHPIPAM session.
type Session struct {
	// The session's configuration.
	Config phpipam.Config

	// The session token.
	Token Token

	// The HTTP client shared by all requests made with this session. It is
	// created on first use by HTTPClient.
	httpClient     *http.Client
	httpClientOnce sync.Once
}

// NewSession creates a new session based off supplied configs. It is up to the
// client for each controller implementation to log in and refresh the token.
// This is provided in the base client.Client implementation.
//
// Non-zero fields of each supplied config override the defaults from
// phpipam.DefaultConfigProvider, in order.
func NewSession(configs ...phpipam.Config) *Session {
	s := &Session{
		Config: phpipam.DefaultConfigProvider(),
	}
	for _, v := range configs {
		mergeConfig(&s.Config, v)
	}

	return s
}

// mergeConfig copies the non-zero fields of src into dst.
func mergeConfig(dst *phpipam.Config, src phpipam.Config) {
	if src.AppID != "" {
		dst.AppID = src.AppID
	}
	if src.Endpoint != "" {
		dst.Endpoint = src.Endpoint
	}
	if src.Password != "" {
		dst.Password = src.Password
	}
	if src.Username != "" {
		dst.Username = src.Username
	}
	if src.Insecure {
		dst.Insecure = src.Insecure
	}
	if src.Timeout != 0 {
		dst.Timeout = src.Timeout
	}
}

// HTTPClient returns the HTTP client for this session, creating it on first
// use. Reusing a single client allows connections to the PHPIPAM API to be
// kept alive between requests.
//
// Note that by design, the client does not follow redirects - if you get a
// 302 error or some other sort of 300 error from the SDK, please check your
// API endpoints.
func (s *Session) HTTPClient() *http.Client {
	s.httpClientOnce.Do(func() {
		timeout := s.Config.Timeout
		if timeout == 0 {
			timeout = phpipam.DefaultTimeout
		}
		s.httpClient = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: s.Config.Insecure},
				Proxy:           http.ProxyFromEnvironment,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	})
	return s.httpClient
}
