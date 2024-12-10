package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/onlydarel/basic-react-golang-app/internal/api"
	"github.com/onlydarel/basic-react-golang-app/internal/driver"
	"log"
	"os"
)

var db *sql.DB

func main() {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")
	DB_NAME := os.Getenv("DBNAME")
	DB_USER := os.Getenv("DBUSER")
	DB_PASS := os.Getenv("DBPASS")

	// Set the app config using fiber
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  PORT,
		AppName:       "Inspire Property Backend",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173/",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Connect to db
	log.Println("Connecting to db...")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=localhost port=5432 dbname=%s user=%s password=%s", DB_NAME, DB_USER, DB_PASS))
	if err != nil {
		log.Fatal(err)
	}

	api.SetDatabase(db.SQL)

	app.Get("/api/todos", api.GetTodos)
	app.Post("/api/todos", api.AddTodos)
	app.Patch("/api/todos/:id", api.UpdateTodo)
	app.Delete("api/todos/:id", api.DeleteTodo)

	log.Fatal(app.Listen(PORT))
}
