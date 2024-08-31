package handlers

import (
	"database/sql"
	"fmt"
	"japlearning/models"
	"japlearning/utils"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
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

		// Create JWT token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = "1"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Sign the token with your secret key
		tokenString, err := token.SignedString([]byte(SecretKey))
		//tokenString, err := token.SignedString(utils.JWTSecret)
		fmt.Print(tokenString)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		return c.JSON(fiber.Map{
			"token": tokenString,
		})
	}	

	/*
		This will run if the user gives their email to sign
	*/
	if err := c.BodyParser(&userData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
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
		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": liveUserID,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		})

		// Sign the token with your secret key
		tokenString, err := token.SignedString(SecretKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		return c.JSON(fiber.Map{
			"token": tokenString,
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
const SecretKey = "Your-secret-key"

func VerifyJwt(c *fiber.Ctx) error {
	
	authHeader := c.Get("Authorization")
if authHeader == "" {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Missing or invalid token1",
	})
}

 tokenString := authHeader[7:]
 token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
 		return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
 	}
 	return []byte(SecretKey), nil
 })
 if err != nil || !token.Valid{
 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
 		"error": "Invalid or expired token",
 	})
 }
 if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
 	// Check expiration time
 	if exp, ok := claims["exp"].(float64); ok {
 		if time.Unix(int64(exp), 0).Before(time.Now()) {
 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
 				"error": "Token has expired",
 			})
 		}
 	}
 	// Extract user information from claims if needed
 	userIDStr := claims["sub"].(string)
	userID, err := strconv.Atoi( userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
 	c.Locals("user_id", userID)
 } else {
 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
 		"error": "Invalid token claims",
 	})
 }
 // If everything is okay, continue with the request
 return c.Next()
 }
