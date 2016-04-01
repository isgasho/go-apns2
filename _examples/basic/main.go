package main

import (
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
		Badge: 5,
	}

	// Parse the certificate
	cert, err := certificate.ReadP12File2(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	// Setup a new http client with pass the Certificate
	// and host environment (apns2.Development, apns2.Production)
	client, err := apns2.NewClient(cert, apns2.Development)

	if err != nil {
		log.Fatal(err)
	}

	// Send the Push Notification
	resp, err := client.SendPush(payload, deviceToken, &apns2.Headers{})

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
