package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = fmt.Errorf("Invalid token")
var ErrInvalidSigning = fmt.Errorf("Invalid signing method")
var ErrInvalidTokenType = fmt.Errorf("Invalid token type")
var ErrRevoked = fmt.Errorf("Token has been revoked")

type RevokedChecker interface {
	Revoked(jti string) bool
}

type JWTPair struct {
	AccessJWT  string    `json:"access_jwt"`
	RefreshJWT string    `json:"refresh_jwt"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type JWTClaims struct {
	Uuid      string `json:"uuid"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	JTI       string `json:"jti"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	httpOnly           bool
	secureCookies      bool
	secret             string
	domain             string
	issuer             string
	CookieName         string
	RefreshCookieName  string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	sameSite           http.SameSite
	revoked            RevokedChecker
}

func GenerateRandomSecret(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewJWTManager(httpOnly bool, sameSite http.SameSite, secureCookies bool, secret string, issuer string, cookieName string, refreshCookieName string, accessExpiry time.Duration, refreshExpiry time.Duration) *JWTManager {
	secretBytes, err := hex.DecodeString(secret)
	if err != nil || len(secretBytes) == 0 {
		secretBytes = []byte(secret)
	}
	return &JWTManager{
		httpOnly:           httpOnly,
		sameSite:           sameSite,
		secureCookies:      secureCookies,
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
		JTI:       GenerateUUID(),
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
		JTI:       GenerateUUID(),
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

	if J.revoked != nil && claims.JTI != "" && J.revoked.Revoked(claims.JTI) {
		return nil, ErrRevoked
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

func (J *JWTManager) SetCookie(c *gin.Context, name string, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   J.domain,
		MaxAge:   maxAge,
		Secure:   J.secureCookies,
		HttpOnly: J.httpOnly,
		SameSite: J.sameSite,
	})
}

func (J *JWTManager) SetAuthCookies(c *gin.Context, pair *JWTPair, maxAge int) {
	J.SetCookie(c, J.CookieName, pair.AccessJWT, maxAge)
	J.SetCookie(c, J.RefreshCookieName, pair.RefreshJWT, int(J.refreshTokenExpiry.Seconds()))
}

func (J *JWTManager) ClearAuthCookies(c *gin.Context) {
	J.SetCookie(c, J.CookieName, "", -1)
	J.SetCookie(c, J.RefreshCookieName, "", -1)
}

func (J *JWTManager) GetAuthCookie(c *gin.Context) string {
	cookie, _ := c.Cookie(J.CookieName)
	return cookie
}

func (J *JWTManager) GetRefreshCookie(c *gin.Context) string {
	cookie, _ := c.Cookie(J.RefreshCookieName)
	return cookie
}

func (J *JWTManager) RefreshTokenExpiry() time.Duration {
	return J.refreshTokenExpiry
}
