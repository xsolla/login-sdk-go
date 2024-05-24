package login

import (
	"net/http"
	"time"
)

type Adapter struct {
	client           *http.Client
	baseUrl          string
	extraHeaderName  string
	extraHeaderValue string
}

func NewAdapter(
	baseUrl string,
	timeout time.Duration,
	extraHeaderName,
	extraHeaderValue string,
	transport *http.Transport,
) *Adapter {
	//nolint:gosec
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
