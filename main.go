package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "certs/PushChatKey.p12"
	var password = "pushchat"

	cert, key, err := readFile(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	t := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}

	// Create http client

	config := &tls.Config{
		Certificates: []tls.Certificate{t},
	}
	config.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: config}

	if err := http2.ConfigureTransport(transport); err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Transport: transport}

	url := fmt.Sprintf("%v/3/device/%v", "https://api.development.push.apple.com", deviceToken)

	payload := []byte(`{ "aps" : { "alert" : "Hello world", "badge" : "10" } }`)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("NewRequest error %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	fmt.Printf("resp %v", resp)
	if err != nil {
		log.Fatal(err)
	}
}
