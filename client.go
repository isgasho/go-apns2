package apns2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

// Apple Development and Production URLs
const (
	Development = "https://api.development.push.apple.com"
	Production  = "https://api.push.apple.com"
)

type ApnsResponse struct {
	ApnsID string `json:"apns-id,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type ErrorResponse struct {
	Reason    string `json:"reason,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

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
func (c *Client) Send(payload []byte, deviceToken string, headers *Headers) (*http.Response, error) {

	url := fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// Send JSON Headers
	headers.Set(req.Header)

	// Do the request
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
	}

	defer resp.Body.Close()

	return resp, nil
}

// Send a push notification with payload []byte and device token
func (c *Client) SendPush(payload []byte, deviceToken string, headers *Headers) (*ApnsResponse, error) {

	url := fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// Send JSON Headers
	headers.Set(req.Header)

	// Do the request
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	output := ApnsResponse{}

	if resp.StatusCode == http.StatusOK {
		output.ApnsID = resp.Header.Get("apns-id")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var errorResponse ErrorResponse
	json.Unmarshal(body, &errorResponse)

	if errorResponse.Reason != "" {
		output.Reason = errorResponse.Reason
	}

	/*b, err := json.Marshal(output)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}*/

	return &output, nil
}
