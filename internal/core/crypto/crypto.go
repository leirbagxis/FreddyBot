package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

var (
	ErrInvalidKey     = errors.New("encryption key must be 32 bytes")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
)

func deriveKey(masterKey []byte, userID int64) []byte {
	hk := hkdf.New(sha256.New, masterKey, []byte(fmt.Sprintf("user-session-%d", userID)), nil)
	key := make([]byte, 32)
	if _, err := io.ReadFull(hk, key); err != nil {
		panic(err)
	}
	return key
}

func Encrypt(plaintext []byte, masterKey []byte, userID int64) (string, error) {
	if len(masterKey) != 32 {
		return "", ErrInvalidKey
	}

	key := deriveKey(masterKey, userID)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(cipherHex string, masterKey []byte, userID int64) ([]byte, error) {
	if len(masterKey) != 32 {
		return nil, ErrInvalidKey
	}

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return nil, err
	}

	key := deriveKey(masterKey, userID)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

func HashPhone(phone string, masterKey []byte) string {
	mac := hmac.New(sha256.New, masterKey)
	mac.Write([]byte(phone))
	return hex.EncodeToString(mac.Sum(nil))
}
