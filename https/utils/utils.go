package https_utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/IzomSoftware/GinWrapper/common/configuration"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateRandomSecret(size int) (string, error) {
    bytes := make([]byte, size)

    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
	
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func getSecret() ([]byte, error) {
    secret := configuration.ConfigHolder.Protections.Tokenizer.TokenizerSecret

    if secret == "" {
        return nil, fmt.Errorf("Empty Secret")
    }

    return base64.URLEncoding.DecodeString(secret)
}

func GenerateToken(id string, exp time.Duration) (string, error) {
    claims := User{
        Username: id,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := getSecret()
    if err != nil {
        return "", err
    }

    return token.SignedString(secret)
}

func ParseToken(tokenString string) (string, error) {
	secret, err := getSecret()

	if err != nil {
		return "", err
	}

    token, err := jwt.ParseWithClaims(tokenString, &User{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, nil
        }
        return secret, nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(*User); ok && token.Valid {
        return claims.Username, nil
    }

    return "", jwt.ErrSignatureInvalid
}