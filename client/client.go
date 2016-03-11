package client

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/net/http2"
)

// APNS
const (
	Development = "https://api.development.push.apple.com"
	Production  = "https://api.push.apple.com"
)

type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
}

func NewClient(certificate tls.Certificate) (*Client, error) {
	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	config.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: config}

	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}

	client := &Client{
		HTTPClient:  &http.Client{Transport: transport},
		Certificate: certificate,
	}

	return client, nil
}
