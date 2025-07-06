package utils

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateRSAKey() (*rsa.PrivateKey, error) {
	// Gera uma chave RSA de 2048 bits
	return rsa.GenerateKey(rand.Reader, 2048)
}
