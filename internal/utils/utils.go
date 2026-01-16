package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"regexp"
	"strings"
)

func GenerateRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func RemoveHTMLTags(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

func NormalizePort(p string) string {
	if p == "" {
		return ":7000"
	}
	if !strings.HasPrefix(p, ":") {
		return ":" + p
	}
	return p
}
