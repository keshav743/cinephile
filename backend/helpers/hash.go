package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func ComparePasswords(password string, hashedPassword string) (error) {
	passwordsMatched := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
	return passwordsMatched
}