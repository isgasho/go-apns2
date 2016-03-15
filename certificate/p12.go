package p12

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

// ReadFile reading a .p12 file
func ReadFile(filename string, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while loading %s: %v", filename, err)
	}

	// Decode the certification
	privateKey, cert, err := pkcs12.Decode(file, password)
	if err != nil {
		return nil, nil, err
	}

	// Verify the certification
	_, err = cert.Verify(x509.VerifyOptions{})
	if err == nil {
		return nil, nil, err
	}

	switch e := err.(type) {
	case x509.CertificateInvalidError:
		switch e.Reason {
		case x509.Expired:
			// TODO Better support for error
		default:
		}
	case x509.UnknownAuthorityError:
		// TODO Better support for error
	default:
	}

	// check if private key is correct
	priv, b := privateKey.(*rsa.PrivateKey)
	if !b {
		return nil, nil, fmt.Errorf("Error with private key")
	}
	return cert, priv, nil
}
