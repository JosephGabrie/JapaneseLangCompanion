package utils

import (
	"database/sql"
	"japlearning/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)
func HashUserPasswords(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err ==nil
}

func CheckRepeatedUsername(username string) (bool, error) {

}

func CheckRepeatedEmail(){
var existingEmail string

}