package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/sger/go-apns2"
	"github.com/sger/go-apns2/certificate"
)

func main() {
	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "../certs/PushChatKey.p12"
	var password = "pushchat"

	// Setup payload must contains an aps root label and alert message
	payload := apns2.Payload{
		Alert: apns2.Alert{
			Body: "Testing HTTP 2"},
		Badge: 50,
	}

	// Marshal the payload structure
	b, err := payload.MarshalJSON()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b)

	// Simple create the payload as []byte.Omit whitespace and line breaks
	// to reduce the size of the payload, improving network performance.
	// payload := []byte(`{"aps":{"alert":"Hello world"}}`)
	// fmt.Println(payload)

	// Parse the certificate
	cert, key, err := certificate.ReadP12File(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	// Create the certificate
	certificate := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}

	// Setup a new http client with pass the Certificate
	// and host environment (apns2.Development, apns2.Production)
	client, err := apns2.NewClient(certificate, apns2.Development)

	if err != nil {
		log.Fatal(err)
	}

	// Send the Push Notification
	resp, err := client.SendPush(b, deviceToken, &apns2.Headers{})

	if err != nil {
		log.Fatal(err)
	}

	// Returns ApnsResponse struct
	/*
		type ApnsResponse struct {
		StatusCode            int
		StatusCodeDescription string
		ApnsID                string `json:"apns-id,omitempty"`
		Reason                string `json:"reason,omitempty"`
	}*/
	fmt.Println(resp)
}
