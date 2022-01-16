package rsamail

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func GenerateKeyPair(publicKeyPath string, privateKeyPath string) *rsa.PrivateKey {
	privateBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(privateBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey
}
