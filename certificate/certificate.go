package certificate

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

// BlockType PEM formatted block (certificate, private key etc)
type BlockType string

func (bt BlockType) String() string {
	return string(bt)
}

// Type of BlockType
const (
	PrivateKey    BlockType = "PRIVATE KEY"
	PublicKey     BlockType = "PUBLIC KEY"
	RSAPrivateKey BlockType = "RSA PRIVATE KEY"
	ECPrivateKey  BlockType = "EC PRIVATE KEY"
	Certificate   BlockType = "CERTIFICATE"
)

// ReadP12File reading a .p12 file
func ReadP12File(filename string, password string) (tls.Certificate, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("Error while loading %s: %v", filename, err)
	}

	// Decode the certification
	privateKey, cert, err := pkcs12.Decode(file, password)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Verify the certification
	_, err = cert.Verify(x509.VerifyOptions{})
	if err == nil {
		return tls.Certificate{}, err
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
		return tls.Certificate{}, fmt.Errorf("Error with private key")
	}

	certificate := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  priv,
		Leaf:        cert,
	}

	//return cert, priv, nil
	return certificate, nil
}

// ReadPemFile parse .pem file returns tls.Certificate, error
func ReadPemFile(filename string, password string) (tls.Certificate, error) {

	var certification tls.Certificate
	var block *pem.Block

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return tls.Certificate{}, err
	}

	if len(bytes) > 0 {
		for {
			block, bytes = pem.Decode(bytes)
			if block == nil {
				break
			}
			switch BlockType(block.Type) {
			case PrivateKey:
				// PrivateKey
			case PublicKey:
				// PublicKey
			case Certificate:
				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					return tls.Certificate{}, err
				}
				certification.Leaf = cert
				certification.Certificate = append(certification.Certificate, block.Bytes)
			case RSAPrivateKey:
				if x509.IsEncryptedPEMBlock(block) {
					bytes, err := x509.DecryptPEMBlock(block, []byte(password))
					if err != nil {
						return tls.Certificate{}, errors.New("Failed to decrypt private key")
					}
					key, err := x509.ParsePKCS1PrivateKey(bytes)
					if err != nil {
						return tls.Certificate{}, errors.New("Failed to parse PKCS1 private key")
					}
					certification.PrivateKey = key
				} else {
					key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
					if err != nil {
						return tls.Certificate{}, errors.New("Failed to parse PKCS1 private key")
					}
					certification.PrivateKey = key
				}
			case ECPrivateKey:
				//ECPrivateKey
			default:
				return tls.Certificate{}, fmt.Errorf("Decode Pem file: encountered unknown block type %s", block.Type)
			}
		}
	}
	return certification, nil
}
