package auth

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidatePassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}