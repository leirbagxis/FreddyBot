package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leirbagxis/FreddyBot/pkg/config"
)

var secreteKey = []byte(config.SecreteKey)

const issuer = "t.me/legendasbrbot"

type ChannelClaims struct {
	ChannelID int64 `json:"channel_id"`
	IsAdmin   bool  `json:"is_admin"`
	TV        int64 `json:"tv"` // token version, para invalidar tokens antigos se necessário
	jwt.RegisteredClaims
}

func GenerateChannelToken(channelID, userID int64, isAdmin bool, tv int64, ttl time.Duration) (string, error) {
	if channelID == 0 || userID == 0 {
		return "", errors.New("channelID and userID are required")
	}
	if ttl <= 0 {
		return "", errors.New("ttl must be > 0")
	}
	if len(secreteKey) < 32 {
		return "", fmt.Errorf("JWT secret too short: use 32+ bytes (got %d)", len(secreteKey))
	}

	now := time.Now()

	claims := ChannelClaims{
		ChannelID: channelID,
		IsAdmin:   isAdmin,
		TV:        tv,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        newJTI(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secreteKey)
}

func ValidateChannelToken(tokenStr string) (*ChannelClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is required")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &ChannelClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secreteKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ChannelClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// valida issuer
	if claims.Issuer != issuer {
		return nil, errors.New("invalid issuer")
	}

	// valida campos mínimos
	if claims.Subject == "" {
		return nil, errors.New("missing subject (userID)")
	}
	if claims.ChannelID == 0 {
		return nil, errors.New("missing channel_id")
	}
	if claims.ID == "" {
		return nil, errors.New("missing jti")
	}

	if claims.TV <= 0 {
		return nil, errors.New("invalid token version")
	}

	return claims, nil
}

func newJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}
