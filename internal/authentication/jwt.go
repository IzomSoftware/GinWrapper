package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var InvalidTokenError = fmt.Errorf("Invalid token")
var UnexpectedSigningMethodError = fmt.Errorf("Unexpected signing method")

type JWTPair struct {
	AccessJWT  string    `json:"access_jwt"`
	RefreshJWT string    `json:"refresh_jwt"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type JWTClaims struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret             string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTManager(secret string, accessExpiry time.Duration, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secret:             secret,
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

func (J *JWTManager) GenerateJWTPair(uuid string, username string) (*JWTPair, error) {
	currentTime := time.Now()
	accessExpiry := currentTime.Add(J.accessTokenExpiry)

	claims := JWTClaims{
		Uuid:     uuid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			Issuer:    "GinWrapper",
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := accessJWT.SignedString(J.secret)
	if err != nil {
		return nil, err
	}

	claims = JWTClaims{
		Uuid:     uuid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(J.refreshTokenExpiry)),
		},
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshStr, err := refreshJWT.SignedString(J.secret)
	if err != nil {
		return nil, err
	}

	return &JWTPair{
		AccessJWT:  accessStr,
		RefreshJWT: refreshStr,
		ExpiresAt:  accessExpiry,
	}, nil
}

func (J *JWTManager) ValidateJWTSigningMethod(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, UnexpectedSigningMethodError
	}
	return J.secret, nil
}

func (J *JWTManager) ValidateJWT(jwtStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &JWTClaims{}, J.ValidateJWTSigningMethod)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, InvalidTokenError
	}
	jwtClaims, _ := token.Claims.(*JWTClaims)
	return jwtClaims, err
}