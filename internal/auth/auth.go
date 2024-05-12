package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidatePassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func GenerateConfirmationToken(uuid uuid.UUID) string {
	hash := sha256.Sum256(uuid.Bytes())
	return string(hex.EncodeToString(hash[:]))
}