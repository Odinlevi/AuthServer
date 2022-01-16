package rsamail

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	log "github.com/sirupsen/logrus"
)

func Decrypt(key *rsa.PrivateKey, cipherText string) string {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, key, []byte(cipherText), nil)
	if err != nil {
		log.Error(err)
	}
	return string(plaintext)
}
