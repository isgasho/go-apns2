package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RobotsAndPencils/buford/payload/badge"

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

// APS is Apple's reserved namespace.
type APS struct {
	// Alert dictionary.
	Alert Alert

	// Badge to display on the app icon.
	// Set to badge.Preserve (default), badge.Clear
	// or a specific value with badge.New(n).
	Badge badge.Badge

	// The name of a sound file to play as an alert.
	Sound string

	// Content available for silent notifications.
	// With no alert, sound, or badge.
	ContentAvailable bool

	// Category identifier for custom actions in iOS 8 or newer.
	Category string
}

// Alert dictionary.
type Alert struct {
	// Title is a short string shown briefly on Apple Watch in iOS 8.2 or newer.
	Title        string   `json:"title,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`

	// Body text of the alert message.
	Body    string   `json:"body,omitempty"`
	LocKey  string   `json:"loc-key,omitempty"`
	LocArgs []string `json:"loc-args,omitempty"`

	// Key for localized string for "View" button.
	ActionLocKey string `json:"action-loc-key,omitempty"`

	// Image file to be used when user taps or slides the action button.
	LaunchImage string `json:"launch-image,omitempty"`
}

func main() {
	var deviceToken = "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	var filename = "certs/PushChatKey.p12"
	var password = "pushchat"

	cert, key, err := readFile(filename, password)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("cert %v key %s", cert, key)

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
		fmt.Printf("error %s", err)
	}

	client := &http.Client{Transport: transport}

	url := fmt.Sprintf("%v/3/device/%v", "https://api.development.push.apple.com", deviceToken)

	payload := APS{
		Alert: Alert{Body: "test message"},
	}

	fmt.Println("json payload", payload)

	b, err := json.Marshal(payload)

	fmt.Println("json payload", b)

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		fmt.Printf("NewRequest error %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Client error %s", err)
	}

	fmt.Printf("resp %v", resp)
}
