package main

import (
	"database/sql" // add this
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // add this

	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/template/html/v2"
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

func getKanaKanjiHandler(c *fiber.Ctx, db *sql.DB) error {
    var kanaKanjiList []KanaKanji

    rows, err := db.Query("SELECT kanakanji_id, character, romanization FROM userprogress")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get user progress"})
    }
    defer rows.Close()

    for rows.Next() {
        var kanaKanji KanaKanji
        if err := rows.Scan(&kanaKanji.KanaKanji_ID, &kanaKanji.Character, &kanaKanji.Romanization); err != nil {
           return c.Status(500).JSON(fiber.Map{"error": "Failed to scan row "})

        }
        kanaKanjiList = append(kanaKanjiList, kanaKanji)
    }

    return c.JSON(kanaKanjiList)
}
/*
TODO:
Make the name sound cooler
*/

//Select kanakanji_id, character, romanization FROM kanakanji
func learnKanaTimer(c *fiber.Ctx, db * sql.DB) error {
    var kanaKanjiID int 
    kanaKanjiList := make([]KanaKanji, 0, 6)
    
     nextSet:= db.QueryRow("SELECT kanakanji_id FROM userprogress WHERE lastlearned = true").Scan(&kanaKanjiID)
    fmt.Println(nextSet);
    rows, err := db.Query("SELECT * FROM kanakanji ORDER BY kanakanji_id ASC LIMIT 7 OFFSET $1", nextSet)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to query data1"})
    }
    defer rows.Close()

    // count := 0
    for rows.Next() {
        var kanaKanji KanaKanji
        if err := rows.Scan(&kanaKanji.KanaKanji_ID, &kanaKanji.Character, &kanaKanji.Romanization); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to scan row"})
        }
        kanaKanjiList = append(kanaKanjiList, kanaKanji)
        // count ++
        // if count >= 7 {
        //     break
        // }
    }

    return c.JSON(kanaKanjiList)
    }
    
 //Check the progression table for LastCompleted if its True grab the next 7 kana
 //Write an edge case that checks that it is in scope

 //Write a while loop
 






/*
Written by Joseph N Gabrie
8/12/24
This is the function that updates the user's flashcards
We expect a JSON list from the frontend that should display:


*/
//func kanaKanjiRetrievalHandler(c *fiber.Ctx, db *sql.DB) error {}

//func lessonUpdateHandler(c *fiber.Ctx, db *sql.DB) error {} 

// func postHandler(c *fiber.Ctx, db *sql.DB) error {
//    newTodo := todo{}
//    if err := c.BodyParser(&newTodo); err != nil {
//        log.Printf("An error occured: %v", err)
//        return c.SendString(err.Error())
//    }
//    fmt.Printf("%v", newTodo)
//    if newTodo.Item != "" {
//        _, err := db.Exec("INSERT into todos VALUES ($1)", newTodo.Item)
//        if err != nil {
//            log.Fatalf("An error occured while executing query: %v", err)
//        }
//    }

//    return c.Redirect("/")
// }

func main() {
   connStr := "postgres://postgres:Josephg57!@localhost:5432/kanaKanji?sslmode=disable"
   // Connect to database
   db, err := sql.Open("postgres", connStr)
   if err != nil {
       log.Fatal(err)
   }
   app := fiber.New()

   app.Get("/kanakanji", func(c *fiber.Ctx) error {
        return getKanaKanjiHandler(c, db)
   })
   
   app.Get("/", func(c *fiber.Ctx) error {
        return learnKanaTimer(c, db)
   })
   port := os.Getenv("PORT")
   if port == "" {
       port = "3000"
   }
   app.Static("/", "../Frontend/public")
   log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}