package rsamail

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	log "github.com/sirupsen/logrus"
)

func Encrypt(key *rsa.PrivateKey, message string) string {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, &key.PublicKey, []byte(message), nil)
	if err != nil {
		log.Error(err)
	}
	return string(ciphertext)
}
