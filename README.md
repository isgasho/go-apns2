[![Build Status](https://travis-ci.org/sger/go-apns2.svg?branch=master)](https://travis-ci.org/sger/go-apns2)

# Go Apns2
	This project is under development

Go package for HTTP/2 [Apple Push Notification Service](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/ApplePushService.html).

## Installation

 Via `go-get`:

```sh
$ go get github.com/sger/go-apns2
$ cd go-apns2/_examples
$ cd basic
$ go build
$ ./basic
```

## Documentation

```sh
$ godoc .
$ godoc -http=:6060
```

## Usage

## Simple example

```go
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
		Badge: 5,
	}

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
```

## Goroutines and channels example

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

var status bool
var payloads []apns2.Payload
var payloadsProcessed int
var totalPayloads int
var apns []*apns2.ApnsResponse

func main() {
	status = true
	statusChannel := make(chan int)
	payloadChannel := make(chan *apns2.ApnsResponse)
	totalPayloads = 0

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Hello World %v!", i)
		payload := apns2.Payload{
			Alert: apns2.Alert{
				Body: message},
		}
		payloads = append(payloads, payload)
	}

	payloadsProcessed = 0
	totalPayloads = len(payloads)

	go sendPayloads(statusChannel, payloadChannel)
	go processPayloadResponses(payloadChannel)

	for {
		if status == false {
			for _, id := range apns {
				fmt.Println(id)
			}
			fmt.Println("Done sending ", totalPayloads, " payloads")
			break
		}
		select {
		case sC := <-statusChannel:
			fmt.Println("Payload received on StatusChannel", sC)
			payloadsProcessed++
			if payloadsProcessed == totalPayloads {
				fmt.Println("Received all Payloads")
				status = false
				close(statusChannel)
				close(payloadChannel)
			}
		}
	}
}

func sendPayloads(statusChannel chan int, payloadChannel chan *apns2.ApnsResponse) {
	time.Sleep(time.Millisecond * 1)
	fmt.Println("Sending", len(payloads), "payloads")

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "../certs/PushChatKey.p12"
	var password = "pushchat"

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

	for i := 0; i < totalPayloads; i++ {
		fmt.Println("sending payload ", i, payloads[i])
		resp, err := client.SendPush(payloads[i], deviceToken, &apns2.Headers{})
		if err != nil {
			log.Fatal(err)
		}
		payloadChannel <- resp
		statusChannel <- 0
	}
}

func processPayloadResponses(payloadChannel chan *apns2.ApnsResponse) {
	for {
		select {
		case pC := <-payloadChannel:
			apns = append(apns, pC)
		}
	}
}
```

## TODO
- [ ] Error Handling
- [ ] Payload Struct
- [ ] Badge
- [ ] Pem Support
- [ ] Tests

Author
-----

__Spiros Gerokostas__ 

- [![](https://img.shields.io/badge/twitter-sger-brightgreen.svg)](https://twitter.com/sger) 
- :email: spiros.gerokostas@gmail.com

License
-----

go-apns2 is available under the MIT license. See the LICENSE file for more info.

