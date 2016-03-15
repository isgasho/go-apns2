package apns2

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

// Apple endpoints
const (
	Development = "https://api.development.push.apple.com"
	Production  = "https://api.push.apple.com"
)

// Client struct with HTTPClient and Certificate as parameters
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
}

// NewClient constructor tls.Certificate parameter
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

func (c *Client) Send(payload []byte, deviceToken string) (*http.Response, error) {

	url := fmt.Sprintf("%v/3/device/%v", Development, deviceToken)

	// Sending the request with valid PAYLOAD (must starts with aps)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// Send JSON Header
	// TODO
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	return resp, nil
}
