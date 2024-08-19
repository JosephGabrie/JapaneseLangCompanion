package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"japlearning/models"
	"japlearning/utils"
	"time"
)


func PostKanaKanji(c *fiber.Ctx, db *sql.DB) error {
	var createProgress models.Progress
	timeCompleted := time.Now()
	nextTimeReview := timeCompleted.Add(24 * time.Hour)

	query := `INSERT INTO userprogress 
	SET timecompleted = $1,
	next_time_review = $2
	mastery_level = 1,
	lastlearned = $3
	WHERE user_id = $4 AND kanaKanji_id = $5
	`
	result, err := db.Exec(query, timeCompleted, nextTimeReview, createProgress.LastLearned, 
		createProgress.UserID, createProgress.KanaKanjiID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to execute update query: %v", err),
		})
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to retrieve affected rows: %v", err),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No rows were updated. Check if the user_id and kanakanji_id exist.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(createProgress)
}

func UpdateUserProgress(c *fiber.Ctx, db *sql.DB) error {
	var newProgress models.Progress
	var nextTimeReview time.Time
	timeCompleted := time.Now()


	if err := c.BodyParser(&newProgress); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})

	}


	newUserMastery := utils.CalculateUserMastery(&newProgress, newProgress.UserTypedAnswer)
	nextTimeReview = utils.SetNextTime(newUserMastery)

	query := `UPDATE userprogress
	SET timecompleted = $1,
	next_time_review = $2,
	masterylevel = $3,
	lastlearned = $4
	WHERE user_id = $5 AND kanakanji_id = $6
	`

	result, err := db.Exec(query, timeCompleted, nextTimeReview, newProgress.MasteryLevel, 
		newProgress.LastLearned, newProgress.UserID, newProgress.KanaKanjiID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to execute update query: %v", err),
		})
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to retrive affected rows: %v", err),
		})
	}

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No rows were updated. Check if the user_id and kanakanji_id exist.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(newProgress)
}
