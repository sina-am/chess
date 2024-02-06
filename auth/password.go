package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	if err != nil {
		panic(err)
	}
	return string(bytePassword)
}

func VerifyPassword(plainPassword, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
