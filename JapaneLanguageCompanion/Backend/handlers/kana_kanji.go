package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"japlearning/models"
	"time"
)

func GetLearnKana(c *fiber.Ctx, db *sql.DB) error {

	userID := c.Locals("user_id").(int)
	var kanaKanjiID int
	learnKanaKanjiList := make([]models.KanaKanji, 0, 6)

	nextSet := db.QueryRow("SELECT kanakanji_id FROM userprogress WHERE lastlearned = true and user_id = $1", userID).Scan(&kanaKanjiID)
	err := nextSet
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "last learn has no true values"})
	}

	rows, err := db.Query("SELECT * FROM kanakanji ORDER BY kanakanji_id ASC LIMIT 7 OFFSET $1", nextSet)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query data1"})
	}
	defer rows.Close()

	for rows.Next() {
		var kanaKanji models.KanaKanji
		if err := rows.Scan(&kanaKanji.KanaKanji_ID, &kanaKanji.Character, &kanaKanji.Romanization); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan row"})
		}
		learnKanaKanjiList = append(learnKanaKanjiList, kanaKanji)

	}

	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Errro occured during row iteration"})
	}
	// instead of returning a webpage we are going to return a json list so that in the future we can make dynamic changes to the content of the json file
	return c.JSON(learnKanaKanjiList)

}

func GetReviewKanaKanji(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id").(int)

	currentTime := time.Now()
	var reviewKanaKanjiList []models.KanaKanji

	//Get all  records that are due for review
	rows, err := db.Query(`
        SELECT kk.kanakanji_id, kk.character, kk.romanization
        FROM kanakanji kk
        INNER JOIN userprogress up ON kk.kanakanji_id = up.kanakanji_id
        WHERE up.next_time_review <= $1 AND up.user_id = $2 `, currentTime, userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to scan row"})
	}
	defer rows.Close()

	for rows.Next() {
		var kanaKanji models.KanaKanji
		if err := rows.Scan(&kanaKanji.KanaKanji_ID, &kanaKanji.Character, &kanaKanji.Romanization); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan rows"})
		}
		reviewKanaKanjiList = append(reviewKanaKanjiList, kanaKanji)
	}
	if err := rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error occured during row iteration"})
	}
	return c.JSON(reviewKanaKanjiList)
}

func CreateSchema(c *fiber.Ctx, db *sql.DB) error {
	// Define the schema name
	schemaName := "test_schema"

	// Prepare the SQL statement for schema creation
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Error creating schema: %v", err))
	}

	return c.SendString("Schema created successfully!")
}
