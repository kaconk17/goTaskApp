package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	  }
	host := os.Getenv("HOST")
	database := os.Getenv("DATABASE")
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")

	connStr := "postgresql://"+user+":"+pass+"@"+host+"/"+database+"?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c,db)
	})
	
   port := os.Getenv("PORT")
   if port == "" {
       port = "8000"
   }
  
   log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello")
}