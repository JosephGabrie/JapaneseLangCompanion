package main

import (
   "database/sql" // add this
   "fmt"
   "log"
   "os"
   "encoding/json"
   

   _"github.com/lib/pq" // add this
   "github.com/gofiber/fiber/v2"
   "github.com/gofiber/template/html/v2"
)


type FlashcardRequest struct {
    UserID          int    `json:"user_id"`
    MinMasteryLevel int    `json:"min_mastery_level"`
    MaxMasteryLevel int    `json:"max_mastery_level"`
    CharacterType   string `json:"character_type"`
}

type Flashcard struct {
    Character    string `json:"character"`
    Romanization string `json:"romanization"`
    MasteryLevel int    `json:"mastery_level"`
}

func getFlashcardsHandler(c *fiber.Ctx, db *sql.DB) error {
    var FlashcardRequest
}
func



/*
TO DO:
Change port before going to production
*/
func main() {
    connStr := "postgres://postgres:Josephg57!@localhost:5432/kanaKanji?sslmode=disable"
   // Connect to database
   db, err := sql.Open("postgress", connStr)
   if err != nil {
    log.Fatal(err)
   }
   app := fiber.New()

  
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
