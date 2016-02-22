package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	var filename = "certs/PushChatKey.p12"
	var password = "pushchat"

	// load the .p12 file
	p12, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Errorf("Can't load %s: %v", filename, err)
	}

	// decode the .p12 file
}
