package main

import (
	"database/sql" // add this
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // add this

	"github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
)

type KanaKanji struct {
    KanaKanji_ID           int    `json:"kanakanji_id"`
    Character    string `json:"character"`
    Romanization string `json:"romanization"`
}

type Progress struct {
	UserID          int`json:"user_id"`// Capitalized field name
	KanaKanjiID     int`json:"kanakanji_id"`// Added missing colon and quotes
	TimeCompleted   time.Time `json:"timecompleted"`// Consistent naming and JSON tag
	NextTimeReview  time.Time `json:"next_time_review"`// Consistent naming and JSON tag
	MasteryLevel    int`json:"masterylevel"`// JSON tag fixed
	LastLearned     bool`json:"lastlearned"`// JSON tag fixed
    
}

type Users struct {
    User_ID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}


//Select kanakanji_id, character, romanization FROM kanakanji
func GetLearnKana(c *fiber.Ctx, db * sql.DB) error {
    var kanaKanjiID int 
    learnKanaKanjiList := make([]KanaKanji, 0, 6)
    
    nextSet:= db.QueryRow("SELECT kanakanji_id FROM userprogress WHERE lastlearned = true").Scan(&kanaKanjiID)
    err := nextSet
    if err != nil{
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
    /*
    This function is very similar to getlearnkana but the main difference is that it checks if the nexttime review is before our current time.
    and instead of grabing a limit of 7 its going to grab all elements that fit the condition.
    */
    
    func GetReviewKanaKanji(c *fiber.Ctx, db *sql.DB) error{
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
            if err := rows.Scan(&kanaKanji.KanaKanji_ID, &kanaKanji.Character, &kanaKanji.Romanization); err != nil{
                return c.Status(500).JSON(fiber.Map{"error": "Failed to scan rows"})
            }
            reviewKanaKanjiList = append(reviewKanaKanjiList, kanaKanji)
        }
        if err := rows.Err(); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error occured during row iteration"})   
        }
    return c.JSON(reviewKanaKanjiList)
    }
    type todo struct {
        Item string
     }
     
     func postLearnKana(c *fiber.Ctx, db *sql.DB) error {
        newTodo := todo{}
        if err := c.BodyParser(&newTodo); err != nil {
            log.Printf("An error occured: %v", err)
            return c.SendString(err.Error())
        }
        fmt.Printf("%v", newTodo)

        return c.Redirect("/")
     }

          
     func updateUserProgress(c *fiber.Ctx, db *sql.DB)error {
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

    func deleteKanaKanji(c *fiber.Ctx, db *sql.DB) error {
        kanaKanjiToDelete := c.Query("")
        db.Exec("DELETE from todos WHERE item=$1", kanaKanjiToDelete)
    return c.SendString("deleted")
    }


    func calculateUserMastery(progress *Progress, userTypedAnswer bool ) int {
        
        masterylevel := progress.MasteryLevel
        fmt.Println(masterylevel)
        if userTypedAnswer {
            masterylevel++
        } else if !userTypedAnswer{
            masterylevel--
        }

        progress.MasteryLevel = masterylevel
        fmt.Println("mastery level in calculate user mastery ", masterylevel)
         if masterylevel > 0 && masterylevel < 4{
            return 1
        } else if masterylevel >= 4 && masterylevel < 6{
            return 2
        } else if masterylevel >= 6 && masterylevel < 7{
            return 3
        } else if masterylevel >= 7{
            return 4
        }
        return 0
    }

    func setNextTime(userAnswer int) time.Time {
        currentTime := time.Now()
        var nextTime time.Time

        switch userAnswer {
        case 1:
            nextTime = currentTime.Add(24 * time.Hour) // 24 hours later
        case 2:
            nextTime = currentTime.Add(72 * time.Hour) // 72 hours (3 days) later
        case 3:
            nextTime = currentTime.AddDate(0, 0, 7) // 1 week later
        case 4:
            nextTime = currentTime.AddDate(0, 1, 14) // 1 month and 2 weeks later
        default:
            nextTime = currentTime // If userAnswer doesn't match any case, keep it as current time
        }

        return nextTime
    }

// func update_user_kanaKanji(c *fiber.Ctx, db * sql.DB) error {
//     old.kana
// }
func main() {
   connStr := "postgres://postgres:Josephg57!@localhost:5432/kanaKanji?sslmode=disable"
   // Connect to database
   db, err := sql.Open("postgres", connStr)
   if err != nil {
       log.Fatal(err)
   }

   app := fiber.New()


   app.Use(cors.New(cors.Config{
    AllowOrigins: "http://127.0.0.1:5500",
    AllowHeaders: "Origin, Content-Type, Accept",
   }))

   app.Get("/", func(c *fiber.Ctx) error {
        return GetLearnKana(c, db)
   })
   app.Get("/reviewKanaKanji", func(c *fiber.Ctx) error {
    return GetReviewKanaKanji(c, db)
   })
   app.Post("/", func(c *fiber.Ctx) error {
    return postLearnKana(c, db)
   })

   app.Put("/", func(c *fiber.Ctx) error {
    return updateUserProgress(c, db)
   })
   app.Delete("/", func(c *fiber.Ctx) error {
        return deleteKanaKanji(c, db)
   })

   port := os.Getenv("PORT")
   if port == "" {
       port = "3000"
   }
   log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}