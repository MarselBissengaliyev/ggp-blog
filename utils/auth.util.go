package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic("error occured while hashing password: ", err)
	}

	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(providedPassword))
	isValid := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		isValid = false
	}

	return isValid, msg
}
