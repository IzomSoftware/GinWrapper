package https_utils

import (
	"time"

	"github.com/SyNdicateFoundation/GinWrapper/common/configuration"
	"github.com/dgrijalva/jwt-go"
)

var secret = []byte(configuration.ConfigHolder.Tokenizer.TokenizerSecret)

type User struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(id string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, User{
		Username: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	})
	return token.SignedString(secret)
}

func ParseToken(token string) (string, error) {
	tk, err := jwt.ParseWithClaims(token, &User{}, func(tk *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if claim, ok := tk.Claims.(*User); ok && tk.Valid {
		return claim.Username, nil
	}
	return "", err
}
