package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ValidateHash(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateOTP() (string, error) {
	var bytes [3]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", (uint32(bytes[0])<<16|uint32(bytes[1])<<8|uint32(bytes[2]))%1000000), nil
}

func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
