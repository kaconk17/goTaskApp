package main

import (
	"database/sql"
	"fmt"

	//"html"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	  }
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DB_PORT")
	database := os.Getenv("DATABASE")
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASSWORD")

	
	connStr := "postgresql://"+dbuser+":"+dbpass+"@"+host+":"+dbPort+"/"+database+"?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c,db)
	})
	app.Get("/all", func(c *fiber.Ctx) error {
		return getAll(c,db)
	})
	app.Post("/add", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})
	
   port := os.Getenv("PORT")
   if port == "" {
       port = "8000"
   }
  
   log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	
	var res string
	var tasks []string

	rows, err := db.Query("select * from tb_task")
	
	if err != nil {
		log.Fatalln(err)
		c.JSON("Terjadi kesalahan")
	}

	for rows.Next() {
		rows.Scan(&res)
		tasks = append(tasks, res)
	}
	return c.Render("index", fiber.Map{
		"Tasks":tasks,
	})
	
}

func getAll(c *fiber.Ctx, db *sql.DB) error {

	return c.JSON("data")
}
type newtask struct {
	nama string;
	isi string;
	tanggal string;
}
func postHandler(c *fiber.Ctx, db *sql.DB) error {
	task := newtask{}
	type setatus struct {
		success bool;
		pesan string;
	}
	
	hasil := setatus{
		success: false,
		pesan: "Error",
	}
	if err := c.BodyParser(&task); err != nil {
		log.Printf("An error occured: %v", err)
		return c.JSON(err.Error())
	}
	if task.isi != "" {
		_, err := db.Exec("INSERT into tb_task VALUES ($1,$2,$3)", task.nama, task.isi, task.tanggal)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
		
		hasil = setatus{
			success: true,
			pesan: "Insert success",
		}

	}
	
	return c.JSON(hasil)
 }

 func putHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello")
 }

 func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello")
 }