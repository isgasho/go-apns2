package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

const (
	Development = "https://api.development.push.apple.com"
	Production  = ""
)

func main() {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "certs/PushChatKey.p12"
	var password = "pushchat"

	// POST URL
	url := fmt.Sprintf("%v/3/device/%v", Development, deviceToken)

	// Setup payload
	payload := []byte(`{ "aps" : { "alert" : "Hello world" } }`)

	cert, key, err := readFile(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	t := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{t},
	}
	config.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: config}

	if err := http2.ConfigureTransport(transport); err != nil {
		log.Fatal(err)
	}

	// Create http client
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	fmt.Printf("resp %v", resp)

	if err != nil {
		log.Fatal(err)
	}
}
