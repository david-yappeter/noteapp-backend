package tools

import "golang.org/x/crypto/bcrypt"

var hashCost = 10

func PasswordHash(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

func PasswordCompare(hashed string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
