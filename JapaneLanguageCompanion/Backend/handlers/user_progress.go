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
	userID := c.Locals("user_id").(int)

	query := `INSERT INTO userprogress 
		(timecompleted, next_time_review, mastery_level, lastlearned, user_id, kanakanji_id)
		VALUES ($1, $2, 1, $3, $4, $5)
		ON CONFLICT (user_id, kanakanji_id) DO UPDATE
		SET timecompleted = EXCLUDED.timecompleted,
		    next_time_review = EXCLUDED.next_time_review,
		    mastery_level = EXCLUDED.mastery_level,
		    lastlearned = EXCLUDED.lastlearned`


	result, err := db.Exec(query, timeCompleted, nextTimeReview, createProgress.LastLearned, 
		userID, createProgress.KanaKanjiID)

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
	createProgress.UserID = userID
	return c.Status(fiber.StatusOK).JSON(createProgress)
}

func UpdateUserProgress(c *fiber.Ctx, db *sql.DB) error {
	var newProgress models.Progress
	var nextTimeReview time.Time
	timeCompleted := time.Now()
	userID := c.Locals("user_id").(int)



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
		newProgress.LastLearned, userID, newProgress.KanaKanjiID)

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

	newProgress.UserID = userID
	return c.Status(fiber.StatusOK).JSON(newProgress)
}

