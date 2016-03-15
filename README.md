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

```go
package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sger/go-apns2"
	"github.com/sger/go-apns2/certificate"
)

func main() {
	var deviceToken = ""
	var filename = "key.p12"
	var password = ""

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
	client, err := apns2.NewClient(certificate, apns2.Development)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Send(payload, deviceToken, &apns2.Headers{})

	if err != nil {
		log.Fatal(err)
	}

	// Print the response
	fmt.Println(resp)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Body %s\n", string(body))
}
```

## Status

  - ◦ error handling
  - ◦ payload struct

Author
-----

__Spiros Gerokostas__ 

- [![](https://img.shields.io/badge/twitter-sger-brightgreen.svg)](https://twitter.com/sger) 
- :email: spiros.gerokostas@gmail.com

License
-----

go-apns2 is available under the MIT license. See the LICENSE file for more info.

