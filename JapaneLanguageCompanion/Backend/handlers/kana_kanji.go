package handlers
import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
    "github.com/JosephGabrie/JapaneLanguageCompanion/Backend/models"
)
JapaneLanguageCompanion/Backend/models
func GetLearnKanaHandler(c *fiber.Ctx, db *sql.DB) error {
	var kanaKanjiID int
	learnKanaKanjiList := make([]KanaKanji, 0, 6)

	nextSet := db.QueryRow("SELECT kanakanji_id FROM userprogress WHERE lastlearned = true").Scan(&kanaKanjiID)
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
		var kanaKanji KanaKanji
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

func GetReviewKanaKanjiHandler(c *fiber.Ctx, db *sql.DB) error {
	currentTime := time.Now()
	var reviewKanaKanjiList []KanaKanji

	//Get all jabajabhu records that are due for review
	rows, err := db.Query(`
        SELECT kk.kanakanji_id, kk.character, kk.romanization
        FROM kanakanji kk
        INNER JOIN userprogress up ON kk.kanakanji_id = up.kanakanji_id
        WHERE up.next_time_review <= $1 `, currentTime)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to scan row"})
	}
	defer rows.Close()

	for rows.Next() {
		var kanaKanji KanaKanji
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

func updateUserProgress(c *fiber.Ctx, db *sql.DB) error {
	var newProgress Progress
	if err := c.BodyParser(&newProgress); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})

	}

	userTypedAnswer := true
	newUserMastery := calculateUserMastery(&newProgress, userTypedAnswer)

	fmt.Println("mastery level in updateUserProgress", newUserMastery)

	newProgress.NextTimeReview = setNextTime(newUserMastery)
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

func deleteKanaKanjiHandler(c *fiber.Ctx, db *sql.DB) error {
	kanaKanjiToDelete := c.Query("")
	db.Exec("DELETE from todos WHERE item=$1", kanaKanjiToDelete)
	return c.SendString("deleted")
}