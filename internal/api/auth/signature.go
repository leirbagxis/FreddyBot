package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateSignature(channelID, ownerID, secreteKey string) string {
	sum := channelID + ownerID
	data := fmt.Sprintf("signaturePayload:%s", sum)
	mac := hmac.New(sha256.New, []byte(secreteKey))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func ValidateSignature(channelID, ownerID, receivedSig, secreteKey string) bool {

	expectedSig := GenerateSignature(channelID, ownerID, receivedSig)
	return hmac.Equal([]byte(expectedSig), []byte(receivedSig))
}
