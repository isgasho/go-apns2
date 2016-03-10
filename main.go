package main

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
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
	privateKey, cert, err := pkcs12.Decode(p12, password)

	if err != nil {
		fmt.Errorf("Can't decode %v", err)
	}

	// check if certificate is valid

	v, err := cert.Verify(x509.VerifyOptions{})

	switch e := err.(type) {
	case x509.CertificateInvalidError:
		switch e.Reason {
		case x509.Expired:
			fmt.Println("expired")
		default:
			fmt.Errorf("error %v", err)
		}
	case x509.UnknownAuthorityError:
	default:
	}

	fmt.Println("v: ", v)

	priv, ok := privateKey.(*rsa.PrivateKey)

	if !ok {
		fmt.Errorf("error %v", ok)
	}

	fmt.Println("priv: ", priv)

	t := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}

	fmt.Println("t: ", t)

	// create new client

}
