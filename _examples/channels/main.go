package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sger/go-apns2"
	"github.com/sger/go-apns2/certificate"
)

func main() {

	payloads := []apns2.Payload{}

	for i := 0; i < 200; i++ {
		message := fmt.Sprintf("Hello World %v!", i)
		payload := apns2.Payload{
			Alert: apns2.Alert{
				Body: message},
		}
		payloads = append(payloads, payload)
	}

	results := asyncHTTPPosts(payloads)

	for _, result := range results {
		if result != nil {
			fmt.Println(result)
		}
	}
}

func asyncHTTPPosts(payloads []apns2.Payload) []*apns2.ApnsResponse {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "../certs/PushChatKey.p12"
	var password = "pushchat"

	ch := make(chan *apns2.ApnsResponse)
	responses := []*apns2.ApnsResponse{}

	cert, err := certificate.ReadP12File(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	// Setup a new http client
	client, err := apns2.NewClient(cert, apns2.Development)

	if err != nil {
		log.Fatal(err)
	}

	for _, payload := range payloads {
		go func(payload apns2.Payload) {
			fmt.Printf("Sending %v \n", payload)
			resp, err := client.SendPush(payload, deviceToken, &apns2.Headers{})
			if err != nil {
				log.Fatal(err)
			}
			ch <- resp
		}(payload)
	}

	for {
		select {
		case resp := <-ch:
			fmt.Printf("%v was received \n", resp)
			responses = append(responses, resp)
			if len(responses) == len(payloads) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}
