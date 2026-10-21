package authentication

import (
	"crypto/rand"
	"encoding/binary"
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
	buff := make([]byte, 4)
	for {
		if _, err := rand.Read(buff); err != nil {
			return "", err
		}
		n := binary.BigEndian.Uint32(buff)
		if n < 4294963200 {
			return fmt.Sprintf("%06d", n%1000000), nil
		}
	}
}

func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
