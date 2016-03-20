# go-apns2
	This project is under development

[![Build Status](https://travis-ci.org/sger/go-apns2.svg?branch=master)](https://travis-ci.org/sger/go-apns2)

Simple package for HTTP/2 [Apple Push Notification Service](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/ApplePushService.html).

## Installation

 Via `go-get`:

```sh
$ go get github.com/sger/go-apns2
```

## Usage

# Simple example

```go
package main

import (
	"crypto/tls"
	"encoding/json"
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
	}

	// Marshal the payload structure
	b, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b)

	//payload := []byte(`{ "aps" : { "alert" : "Hello world" } }`)
	//fmt.Println(payload)

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

	// Setup a new http client with Certificate and host environment (apns2.Development, apns2.Production)
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
```

# Channels example

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/sger/go-apns2"
	"github.com/sger/go-apns2/certificate"
)

func main() {

	payloads := [][]byte{}

	for i := 0; i < 200; i++ {
		message := fmt.Sprintf("Hello World %v!", i)
		payload := []byte(`{ "aps" : { "alert" : "` + message + `" } }`)
		payloads = append(payloads, payload)
	}

	results := asyncHTTPPosts(payloads)
	for _, result := range results {
		if result != nil {
			fmt.Println(result)
		}
	}
}

func asyncHTTPPosts(payloads [][]byte) []*apns2.ApnsResponse {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "../certs/PushChatKey.p12"
	var password = "pushchat"

	ch := make(chan *apns2.ApnsResponse)
	responses := []*apns2.ApnsResponse{}

	cert, key, err := certificate.ReadP12File(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	certificate := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}

	// Setup a new http client
	client, err := apns2.NewClient(certificate, apns2.Development)

	if err != nil {
		log.Fatal(err)
	}

	for _, payload := range payloads {
		go func(payload []byte) {
			fmt.Printf("Sending %s \n", payload)
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
```

## TODO
- [ ] Error Handling
- [ ] Payload Struct

Author
-----

__Spiros Gerokostas__ 

- [![](https://img.shields.io/badge/twitter-sger-brightgreen.svg)](https://twitter.com/sger) 
- :email: spiros.gerokostas@gmail.com

License
-----

go-apns2 is available under the MIT license. See the LICENSE file for more info.

