package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

var secreteKey = []byte(config.SecreteKey)

type CustomClaims struct {
	ChannelID string
	OwnerID   string
	jwt.RegisteredClaims
}

func GenerateTokenJWT(channelID, ownerID string) (string, error) {
	claims := CustomClaims{
		ChannelID: channelID,
		OwnerID:   ownerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secreteKey)
}

func ValidateTokenJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secreteKey, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
