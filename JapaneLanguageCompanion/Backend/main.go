package main

import (
	"database/sql" // add this
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // add this

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type KanaKanji struct {
    KanaKanji_ID           int    `json:"kanakanji_id"`
    Character    string `json:"character"`
    Romanization string `json:"romanization"`
}

type Progress struct {
    DateCompleted   time.Time `json:"date_completed"`
    NextDate        time.Time `json:"next_date"`        // Consider whether this is the same as `next_time_review`
    MasteryLevel    int       `json:"mastery_level"`
    NextTimeReview  time.Time `json:"next_time_review"` // Renamed for consistency
    LastLearned     bool      `json:"last_learned"`     // Fixed JSON tag and field name
}


type Users struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

//Select kanakanji_id, character, romanization FROM kanakanji
func learnKana(c *fiber.Ctx, db * sql.DB) error {
    var kanaKanjiID int 
    kanaKanjiList := make([]KanaKanji, 0, 6)
    
    nextSet:= db.QueryRow("SELECT kanakanji_id FROM userprogress WHERE lastlearned = true").Scan(&kanaKanjiID)
   
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
        kanaKanjiList = append(kanaKanjiList, kanaKanji)

    }

    // return c.JSON(kanaKanjiList)
     return c.Render("index", fiber.Map{
        "kanaList": kanaKanjiList,
     })
    }

    func deleteKanaKanji(c *fiber.Ctx, db *sql.DB) error {
    kanaKanjiToDelete := c.Query("")
    db.Exec("DELETE from todos WHERE item=$1", kanaKanjiToDelete)
   return c.SendString("deleted")
}

func update_user_kanaKanji(c *fiber.Ctx, db * sql.DB) error {
    old
}
func main() {
   connStr := "postgres://postgres:Josephg57!@localhost:5432/kanaKanji?sslmode=disable"
   // Connect to database
   db, err := sql.Open("postgres", connStr)
   if err != nil {
       log.Fatal(err)
   }
   engine := html.New("../Frontend/views", ".html")
   app := fiber.New(fiber.Config{
    Views: engine,
   })

   app.Get("/", func(c *fiber.Ctx) error {
        return learnKana(c, db)
   })

   app.Delete("/", func(c *fiber.Ctx) error {
        return deleteKanaKanji(c, db)
   })
   port := os.Getenv("PORT")
   if port == "" {
       port = "3000"
   }

   app.Static("/", "../Frontend/public") // add this before starting the app
   log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}