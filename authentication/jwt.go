package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = fmt.Errorf("Invalid token")
var ErrInvalidSigning = fmt.Errorf("Invalid signing method")
var ErrInvalidTokenType = fmt.Errorf("Invalid token type")

type JWTPair struct {
	AccessJWT  string    `json:"access_jwt"`
	RefreshJWT string    `json:"refresh_jwt"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type JWTClaims struct {
	Uuid      string `json:"uuid"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	encryption         bool
	secret             string
	issuer             string
	CookieName         string
	RefreshCookieName  string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func GenerateRandomSecret(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewJWTManager(encryption bool, secret string, issuer string, cookieName string, refreshCookieName string, accessExpiry time.Duration, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		encryption: encryption,
		secret:             secret,
		issuer:             issuer,
		CookieName:         cookieName,
		RefreshCookieName:  refreshCookieName,
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

func (J *JWTManager) GenerateJWTPair(uuid string, username string) (*JWTPair, error) {
	currentTime := time.Now()
	accessExpiry := currentTime.Add(J.accessTokenExpiry)

	claims := JWTClaims{
		Uuid:      uuid,
		Username:  username,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			Issuer:    J.issuer,
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := accessJWT.SignedString([]byte(J.secret))
	if err != nil {
		return nil, err
	}

	claims = JWTClaims{
		Uuid:      uuid,
		Username:  username,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(J.refreshTokenExpiry)),
			Issuer:    J.issuer,
		},
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshStr, err := refreshJWT.SignedString([]byte(J.secret))
	if err != nil {
		return nil, err
	}

	return &JWTPair{
		AccessJWT:  accessStr,
		RefreshJWT: refreshStr,
		ExpiresAt:  accessExpiry,
	}, nil
}

func (J *JWTManager) ValidateJWTSigningMethod(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidSigning
	}
	return []byte(J.secret), nil
}

func (J *JWTManager) ValidateJWT(jwtStr string, tokenType string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &JWTClaims{}, J.ValidateJWTSigningMethod, jwt.WithIssuer(J.issuer))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != tokenType {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

func (J *JWTManager) RefreshToken(refreshStr string) (*JWTPair, error) {
	claims, err := J.ValidateJWT(refreshStr, "refresh")
	if err != nil {
		return nil, err
	}

	return J.GenerateJWTPair(claims.Uuid, claims.Username)
}
func (J *JWTManager) SetAuthCookies(c *gin.Context, pair *JWTPair, maxAge int) {
	c.SetCookie(J.CookieName, pair.AccessJWT, maxAge, "/", "", J.encryption, true)
	c.SetCookie(J.RefreshCookieName, pair.RefreshJWT, int(J.refreshTokenExpiry.Seconds()), "/", "", J.encryption, true)
}

func (J *JWTManager) ClearAuthCookies(c *gin.Context) {
	c.SetCookie(J.CookieName, "", -1, "/", "", J.encryption, true)
	c.SetCookie(J.RefreshCookieName, "", -1, "/", "", J.encryption, true)
}

func (J *JWTManager) GetAuthCookie(c *gin.Context) string {
	cookie, _ := c.Cookie(J.CookieName)
	return cookie
}

func (J *JWTManager) GetRefreshCookie(c *gin.Context) string {
	cookie, _ := c.Cookie(J.RefreshCookieName)
	return cookie
}
