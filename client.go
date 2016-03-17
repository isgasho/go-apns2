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

// ApnsResponse contains apns-id, reason, status code, status code description.
type ApnsResponse struct {
	StatusCode            int
	StatusCodeDescription string
	ApnsID                string `json:"apns-id,omitempty"`
	Reason                string `json:"reason,omitempty"`
}

// ErrorResponse contains reason, timestamp
type ErrorResponse struct {
	Reason    string `json:"reason,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// Client struct with HTTPClient, Certificate, Host as parameters.
type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

// NewClient constructor tls.Certificate parameter.
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

// SendPush a push notification with payload ([]byte), device token, *Headers
// returns ApnsResponse struct
func (c *Client) SendPush(payload []byte, deviceToken string, headers *Headers) (*ApnsResponse, error) {

	url := fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// Send headers to request
	headers.Set(req.Header)

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	apnsResponse := ApnsResponse{}
	apnsResponse.StatusCode = resp.StatusCode
	apnsResponse.StatusCodeDescription = statusCode[resp.StatusCode]

	if resp.StatusCode == http.StatusOK {
		apnsResponse.ApnsID = resp.Header.Get("apns-id")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var errorResponse ErrorResponse
	json.Unmarshal(body, &errorResponse)

	if errorResponse.Reason != "" {
		apnsResponse.Reason = errorReason[errorResponse.Reason]
	}

	return &apnsResponse, nil
}
