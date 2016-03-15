package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sger/go-apns2"
	"github.com/sger/go-apns2/certificate"
)

func main() {

	payload1 := []byte(`{ "aps" : { "alert" : "Hello world 1" } }`)
	payload2 := []byte(`{ "aps" : { "alert" : "Hello world 2" } }`)
	payload3 := []byte(`{ "aps" : { "alert" : "Hello world 3" } }`)
	payload4 := []byte(`{ "aps" : { "alert" : "Hello world 4" } }`)
	payload5 := []byte(`{ "aps" : { "alert" : "Hello world 5" } }`)
	payload6 := []byte(`{ "aps" : { "alert" : "Hello world 6" } }`)

	payloads := [][]byte{payload1, payload2, payload3, payload4, payload5, payload6}

	results := asyncHTTPPosts(payloads)
	for _, result := range results {
		if result != nil {
			fmt.Printf("status: %s\n", result.Status)
		}
	}
}

func asyncHTTPPosts(payloads [][]byte) []*http.Response {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "../certs/PushChatKey.p12"
	var password = "pushchat"

	ch := make(chan *http.Response)
	responses := []*http.Response{}

	cert, key, err := p12.ReadFile(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	certificate := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}

	// Setup a new http client
	client, err := apns2.NewClient(certificate)

	if err != nil {
		log.Fatal(err)
	}

	for _, payload := range payloads {
		go func(payload []byte) {
			fmt.Printf("Sending %s \n", payload)
			resp, err := client.Send(payload, deviceToken)
			if err != nil {
				log.Fatal(err)
			}
			ch <- resp
		}(payload)
	}

	for {
		select {
		case resp := <-ch:
			fmt.Printf("%T was fetched\n", resp)
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
