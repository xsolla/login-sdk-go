package login

import (
	"crypto/tls"
	"net/http"
	"time"
)

const timeout = 3 * time.Second

type Adapter struct {
	client  *http.Client
	baseUrl string
}

func NewAdapter(baseUrl string, ignoreSslErrors bool) *Adapter {
	//nolint:gosec
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: ignoreSslErrors,
		},
	}

	return &Adapter{
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
		baseUrl: baseUrl,
	}
}
