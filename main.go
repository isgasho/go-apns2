package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

// Apple APNS 2
const (
	Development = "https://api.development.push.apple.com"
	Production  = ""
)

func main() {

	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var p12Filename = "certs/PushChatKey.p12"
	var password = "pushchat"

	// POST URL
	url := fmt.Sprintf("%v/3/device/%v", Development, deviceToken)

	// Setup payload must contains an aps root label and alert message
	payload := []byte(`{ "aps" : { "alert" : "Hello world" } }`)

	cert, key, err := readP12File(p12Filename, password)
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

	// Create http client with Transport with Go 1.6 Transport supports HTTP/2
	client := &http.Client{Transport: transport}

	// Sending the request with valid PAYLOAD (must starts with aps)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	// Send JSON Header
	// TODO
	req.Header.Set("Content-Type", "application/json")

	// Do the request
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Read the response
	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(body))
}
