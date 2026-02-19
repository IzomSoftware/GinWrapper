package https_utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/common/configuration"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUser struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTTokenPair struct {
	JWTAccessToken  string `json:"access_token"`
	JWTRefreshToken string `json:"refresh_token"`
	JWTTokenType    string `json:"token_type"`
	JWTExpiresIn    int64  `json:"expires_in"`
}

func GenerateJWTRandomSecret(size int) (string, error) {
	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func getJWTSecret() ([]byte, error) {
	secret := configuration.ConfigHolder.Protections.JWTProtection.JWTSecret

	if secret == "" {
		return nil, fmt.Errorf("Empty Secret")
	}

	return base64.URLEncoding.DecodeString(secret)
}

func GenerateJWTToken(id string, exp time.Duration) (string, error) {
	tokenId := uuid.NewString()

	claims := JWTUser{
		Username: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId,
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	return token.SignedString(secret)
}

func GenerateJWTTokenPair(id string) (*JWTTokenPair, error) {
	accessToken, err := GenerateJWTToken(id, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateJWTToken(id, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &JWTTokenPair{
		JWTAccessToken:  accessToken,
		JWTRefreshToken: refreshToken,
		JWTTokenType:    "Bearer",
		JWTExpiresIn:    900,
	}, nil
}

func ParseToken(tokenString string) (*JWTUser, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTUser{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTUser); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
