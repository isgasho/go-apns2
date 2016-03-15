package apns2

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

// Apple end points
const (
	Development = "https://api.development.push.apple.com"
	Production  = "https://api.push.apple.com"
)

// Client struct with HTTPClient and Certificate as parameters
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

// NewClient constructor tls.Certificate parameter
func NewClient(certificate tls.Certificate, host string) (*Client, error) {
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
		Host:        host,
	}

	return client, nil
}

// Send a push notification with payload []byte and device token
func (c *Client) Send(payload []byte, deviceToken string) (*http.Response, error) {

	url := fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// Send JSON Headers
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
	}

	defer resp.Body.Close()

	return resp, nil
}
