package handlers

import (
	"database/sql"
	"fmt"
	"japlearning/models"
	"japlearning/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetUsersSignIn(c *fiber.Ctx, db *sql.DB) error{
	var userData models.Users
	var liveUser string
	var liveUserID int
	var liveUserEmail string
	var hashedPassword string

		/*
		This will run if the user gives their username to sign
	*/

	if err := c.BodyParser(&userData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	err := db.QueryRow ("SELECT user_id, username, password FROM users WHERE username = $1", userData.Username).Scan(&liveUserID, &liveUser, &hashedPassword)
	if err != nil && err != sql.ErrNoRows{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if liveUser != "" {
			if !utils.CheckPasswordHash(userData.Password, hashedPassword){
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "passwords did not match",
				})
			} 
			return c.JSON(fiber.Map{
				"user_id": liveUserID,
			})
		}
	/*
		This will run if the user gives their email to sign
	*/
	err = db.QueryRow("SELECT user_id, email, password FROM users WHERE email = $1", userData.Email).Scan(&liveUserID, &liveUserEmail, &hashedPassword)
	if err != nil && err != sql.ErrNoRows{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if liveUserEmail != "" {
		if !utils.CheckPasswordHash(userData.Password, hashedPassword){
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "passwords did not match",
			})
		} 
		return c.JSON(fiber.Map{
			"user_id": liveUserID,
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Invalid username or email",
	})

}
 func PostUsers (c *fiber.Ctx, db *sql.DB) error{
	var userData models.RegistrationData
	var existingUsername string
	var existingEmail string
	 timeSignedUp := time.Now()

	 if err := c.BodyParser(&userData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	 }

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
	err = db.QueryRow ("SELECT email FROM users WHERE email = $1", userData.Email).Scan(&existingEmail)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if existingEmail != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already taken",
		})
	}
	
	hashedPassword, err := utils.HashUserPasswords(userData.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	_,err = db.Exec("INSERT INTO Users (username, email, password, time_signed_up) VALUES ($1, $2, $3, $4) ", userData.Username, userData.Email, hashedPassword, timeSignedUp)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to upload user"),
	})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered succesfully",
	})
}


func EditUser(){

}

func DeleteUser(){

}

