package main

import (
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

	// Creating 2000 payloads
	for i := 0; i < 2000; i++ {
		message := fmt.Sprintf("Hello World %v!", i)
		payload := apns2.Payload{
			Alert: apns2.Alert{
				Body: message},
		}
		payloads = append(payloads, payload)
	}

	payloadsProcessed = 0
	totalPayloads = len(payloads)

	// goroutines
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
	var pemFilename = "../certs/ck.pem"
	var password = "pushchat"

	cert, err := certificate.ReadPemFile(pemFilename, password)
	if err != nil {
		log.Fatal(err)
	}

	// Setup a new http client
	client, err := apns2.NewClient(cert, apns2.Development)

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
