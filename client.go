package apns2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/http2"
)

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
func (c *Client) SendPush(payload interface{}, deviceToken string, headers *Headers) (*ApnsResponse, error) {

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v/3/device/%v", c.Host, deviceToken)

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
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
