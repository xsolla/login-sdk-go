package login

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Adapter struct {
	client           *http.Client
	baseUrl          string
	extraHeaderName  string
	extraHeaderValue string
}

func NewAdapter(baseUrl string, ignoreSslErrors bool, timeout time.Duration, extraHeaderName, extraHeaderValue string) *Adapter {
	//nolint:gosec
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: ignoreSslErrors,
		},
	}

	return &Adapter{
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
		baseUrl:          baseUrl,
		extraHeaderName:  extraHeaderName,
		extraHeaderValue: extraHeaderValue,
	}
}
