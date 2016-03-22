package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/sger/go-apns2"
	"github.com/sger/go-apns2/certificate"
)

var applicationStatus bool
var payloads []apns2.Payload
var payloadsProcessed int
var totalPayloadCount int
var apns []*apns2.ApnsResponse

func main() {
	applicationStatus = true
	statusChannel := make(chan int)
	payloadChannel := make(chan *apns2.ApnsResponse)
	totalPayloadCount = 0

	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Hello World %v!", i)
		payload := apns2.Payload{
			Alert: apns2.Alert{
				Body: message},
		}
		payloads = append(payloads, payload)
	}

	payloadsProcessed = 0
	totalPayloadCount = len(payloads)

	go sendPayloads(statusChannel, payloadChannel)
	go processPayloadResponses(payloadChannel)

	for {
		if applicationStatus == false {

			for _, id := range apns {
				fmt.Println(id)
			}

			fmt.Println("Done sending ", totalPayloadCount, " payloads")
			break
		}
		select {
		case sC := <-statusChannel:
			fmt.Println("Payload received on StatusChannel", sC)

			payloadsProcessed++

			if payloadsProcessed == totalPayloadCount {
				fmt.Println("Received all Payloads")
				applicationStatus = false
				close(statusChannel)
				close(payloadChannel)
			}
		}
	}
}

func sendPayloads(statusChannel chan int, payloadChannel chan *apns2.ApnsResponse) {
	time.Sleep(time.Millisecond * 1)
	fmt.Println("Grabbing ", len(payloads), "payloads")

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

	for i := 0; i < totalPayloadCount; i++ {

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
