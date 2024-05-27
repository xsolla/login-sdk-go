package login

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func NewDefaultTransport(ignoreSslErrors bool) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: ignoreSslErrors,
		},
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   1 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
