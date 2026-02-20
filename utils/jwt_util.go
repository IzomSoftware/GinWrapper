package jwt_util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/configuration"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUser struct {
	Username string `json:"username"`
	JWTType  string `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTPair struct {
	AccessJWT    string `json:"access_token"`
	RefreshJWT   string `json:"refresh_token"`
	JWTType      string `json:"token_type"`
	JWTExpiresIn int64  `json:"expires_in"`
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
		return nil, fmt.Errorf("empty secret")
	}

	return base64.URLEncoding.DecodeString(secret)
}

func GenerateJWT(id string, tokenType string, exp time.Duration) (string, error) {
	claims := JWTUser{
		Username: id,
		JWTType:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func GenerateJWTPair(id string) (*JWTPair, error) {
	accessToken, err := GenerateJWT(id, "access", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateJWT(id, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &JWTPair{
		AccessJWT:    accessToken,
		RefreshJWT:   refreshToken,
		JWTType:      "Bearer",
		JWTExpiresIn: 900,
	}, nil
}

func ParseJWT(tokenString string) (*JWTUser, error) {
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

func ValidateJWT(tokenString string) (*JWTUser, error) {
    claims, err := ParseJWT(tokenString)
    if err != nil {
        return nil, err
    }

    if claims.JWTType != "access" {
        return nil, fmt.Errorf("invalid token type")
    }

    return claims, nil
}

func ValidateRefreshToken(tokenString string) (*JWTUser, error) {
    claims, err := ParseJWT(tokenString)
    if err != nil {
        return nil, err
    }

    if claims.JWTType != "refresh" {
        return nil, fmt.Errorf("invalid token type")
    }

    return claims, nil
}