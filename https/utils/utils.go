package https_utils

import (
    "time"

    "github.com/IzomSoftware/GinWrapper/common/configuration"
    "github.com/golang-jwt/jwt/v5"
)

var secret = []byte(configuration.ConfigHolder.Tokenizer.TokenizerSecret)

type User struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
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
    return token.SignedString(secret)
}

func ParseToken(tokenString string) (string, error) {
    token, err := jwt.ParseWithClaims(tokenString, &User{}, func(token *jwt.Token) (interface{}, error) {
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