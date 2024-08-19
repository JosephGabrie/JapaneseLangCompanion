package handlers

import (
	"fmt"
	"time"
	"database/sql"
	"japlearning/models"


)

func GetUserSignIn(c *fiber.Ctx, db *sql.DB){
	var existingUsername string
	var userData models.RegistrationData
	
	err := db.QueryRow ("SELECT username FROM users WHERE username = $1", userData.Username).Scan(&existingUsername)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if existingUsername != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username already taken",
		})
	}
}

func PostUserSignUp()
	var createUser models.RegistrationData
	var signUpTime time.Time
	var userpassword string

	hash, _ := utils.HashUserPasswords(createUser.Password)
	
	
	query = `
	INSERT INTO users(username, email, password) VALUES $1, $2, $3
	`

	

func EditUser()

func DeleteUser()

