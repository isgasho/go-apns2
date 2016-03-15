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
	/*
		var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
		var filename = "certs/PushChatKey.p12"
		var password = "pushchat"

		// Setup payload must contains an aps root label and alert message
		payload := []byte(`{ "aps" : { "alert" : "Hello world" } }`)

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

		resp, err := client.Send(payload, deviceToken)

		if err != nil {
			log.Fatal(err)
		}

		// Read the response
		fmt.Println(resp)

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

	fmt.Printf("Body %s\n", string(body))*/

	payload1 := []byte(`{ "aps" : { "alert" : "Hello world 1" } }`)
	payload2 := []byte(`{ "aps" : { "alert" : "Hello world 2" } }`)
	payload3 := []byte(`{ "aps" : { "alert" : "Hello world 3" } }`)

	payloads := [][]byte{payload1, payload2, payload3}

	results := asyncHTTPPosts(payloads)
	for _, result := range results {
		if result != nil {
			fmt.Printf("status: %s\n", result.Status)
		}
	}
}

func asyncHTTPPosts(payloads [][]byte) []*http.Response {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "certs/PushChatKey.p12"
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

	for _, item := range payloads {
		go func(item []byte) {
			fmt.Printf("Sending %s \n", item)
			resp, err := client.Send(item, deviceToken)
			if err != nil {
				log.Fatal(err)
			}
			ch <- resp
		}(item)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r)
			responses = append(responses, r)
			if len(responses) == len(payloads) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}
