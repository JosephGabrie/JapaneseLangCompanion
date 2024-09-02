package main

import (
	"database/sql" // add this
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"japlearning/handlers"
)

//	func update_user_kanaKanji(c *fiber.Ctx, db * sql.DB) error {
//	    old.kana
//	}
func main() {
	// Connect to database

	connStr := "user=postgres.ncizsrxcufmvbesuijkj password=Coldshoulder2345! host=aws-0-us-west-1.pooler.supabase.com port=6543 dbname=postgres sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:5500",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	/*
		app.Post("/signin", func(c *fiber.Ctx) error {
			return handlers.GetUsersSignIn(c, db)
		})

		app.Post("/signup", func(c *fiber.Ctx) error {
			return handlers.PostUsers(c, db)
		})
	*/
	//app.Get("/",/* handlers.VerifyJwt,*/ func(c *fiber.Ctx) error {
	//	return handlers.GetLearnKana(c, db)
	//})
	/*
		app.Get("/reviewKanaKanji", func(c *fiber.Ctx) error {
			return handlers.GetReviewKanaKanji(c, db)
		})

		app.Put("/update", func(c *fiber.Ctx) error {
			return handlers.UpdateUserProgress(c, db)
		})
	*/

	app.Post("/", func(c *fiber.Ctx) error {

		return handlers.CreateSchema(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
