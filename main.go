package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/pkcs12"
	"golang.org/x/net/http2"
)

type response struct {
	// Reason for failure
	Reason string `json:"reason"`
	// Timestamp for 410 StatusGone (ErrUnregistered)
	Timestamp int64 `json:"timestamp"`
}

// Headers sent with a push to control the notification (optional)
type Headers struct {
	// ID for the notification. Apple generates one if ommitted.
	// This should be a UUID with 32 lowercase hexadecimal digits.
	// TODO: use a UUID type.
	ID string

	// Apple will retry delivery until this time. The default behavior only tries once.
	Expiration time.Time

	// Allow Apple to group messages to together to reduce power consumption.
	// By default messages are sent immediately.
	LowPriority bool

	// Topic for certificates with multiple topics.
	Topic string
}

type APS struct {
	Alert Alert
}

type Alert struct {
	Body string `json:"body,omitempty"`
}

func main() {
	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "certs/PushChatKey.p12"
	var password = "pushchat"

	// Load the .p12 file
	p12, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Printf("Can't load %s: %v", filename, err)
	}

	// Decode the .p12 file
	privateKey, cert, err := pkcs12.Decode(p12, password)

	if err != nil {
		fmt.Printf("Can't decode %v", err)
	}

	// Perform a verification check if certificate is valid
	_, err = cert.Verify(x509.VerifyOptions{})
	if err != nil {
		fmt.Printf("Error while verifying: %v", err)
	}

	switch e := err.(type) {
	case x509.CertificateInvalidError:
		switch e.Reason {
		case x509.Expired:
			fmt.Println("expired")
		default:
			fmt.Printf("error %v", err)
		}
	case x509.UnknownAuthorityError:
		fmt.Println("UnknownAuthorityError")
	default:
	}

	// assert that private key is RSA
	priv, ok := privateKey.(*rsa.PrivateKey)
	fmt.Println("priv ", priv)
	if !ok {
		fmt.Println("ok: ", ok)
	}

	//get cert, prev

	// Create certification
	t := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  priv,
		Leaf:        cert,
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{t},
	}

	config.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: config}

	if err := http2.ConfigureTransport(transport); err != nil {
		fmt.Println("err: ", err)
	}

	// create new client
	client := &http.Client{Transport: transport}

	//fmt.Println("client: ", client)

	//send data
	urlStr := fmt.Sprintf("%v/3/device/%v", "https://api.development.push.apple.com", deviceToken)
	fmt.Println("urlStr: ", urlStr)
	//
	payload := APS{
		Alert: Alert{Body: "test message"},
	}

	b, err := json.Marshal(payload)

	fmt.Println("json payload", b)

	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(b))
	if err != nil {
		fmt.Println("error POST", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apns-priority", "5")
	req.Header.Set("apns-expiration", strconv.FormatInt(time.Now().Unix(), 10))

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error client", err)
	}

	//defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("header ", resp.Header)
		fmt.Println("apns-id", resp.Header.Get("apns-id"))

	}

	body, err := ioutil.ReadAll(resp.Body)

	var response response
	json.Unmarshal(body, &response)

	fmt.Println("response", &response)

}
