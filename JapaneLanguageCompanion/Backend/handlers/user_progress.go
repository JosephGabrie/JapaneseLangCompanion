package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"japlearning/models"
	"japlearning/utils"
)

func UpdateUserProgress(c *fiber.Ctx, db *sql.DB) error {
	var newProgress models.Progress
	if err := c.BodyParser(&newProgress); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})

	}

	userTypedAnswer := true
	newUserMastery := utils.CalculateUserMastery(&newProgress, userTypedAnswer)

	fmt.Println("mastery level in updateUserProgress", newUserMastery)

	newProgress.NextTimeReview = utils.SetNextTime(newUserMastery)
	fmt.Println(newProgress.NextTimeReview)
	query := `UPDATE userprogress
	SET timecompleted = $1,
	next_time_review = $2,
	masterylevel = $3,
	lastlearned = $4
	WHERE user_id = $5 AND kanakanji_id = $6
	`

	result, err := db.Exec(query, newProgress.TimeCompleted, newProgress.NextTimeReview, newProgress.MasteryLevel, newProgress.LastLearned, newProgress.UserID, newProgress.KanaKanjiID)
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
