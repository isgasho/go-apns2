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

	fmt.Println(resp)
}
