package conn

import (
	"net"
	"net/http"
	"time"
)

// DefaultOutboundTimeout is used whenever a schedule does not define its own timeout.
// http.Client{Timeout: 0} means "wait forever", which would pin a job goroutine
// on an unresponsive endpoint, so zero is never passed through.
const DefaultOutboundTimeout = 30 * time.Second

// NewOutboundClient builds the HTTP client used for scheduled requests and webhooks.
// Redirects are not followed: a cron target answering 302 should be recorded as such
// rather than silently resolved somewhere else.
func NewOutboundClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = DefaultOutboundTimeout
	}
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: timeout,
			IdleConnTimeout:       60 * time.Second,
			MaxIdleConns:          20,
			MaxIdleConnsPerHost:   5,
		},
	}
}
