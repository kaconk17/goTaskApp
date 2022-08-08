package main

import (
	"database/sql"
	"net/http"

	//"html"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func connectDatabase(){
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
	//defer db.Close()
	DB = db
}

func main() {
	

	router := gin.Default()
	connectDatabase()
	router.LoadHTMLFiles("views/index.html")
   port := os.Getenv("PORT")
   router.GET("/",func(ctx *gin.Context) {
	   ctx.HTML(http.StatusOK, "index.html", gin.H{
		   "title":"Task App",
		   "baseurl":ctx.Request.URL.Path,
		})
	})

	router.POST("/add",postHandler)
	router.GET("/all",getAll)

	
	if port == "" {
		port = "8000"
	}
   router.Run(":"+port)
   
}
type Task struct {
	ID int `json:"id"`
	Nama string `json:"nama"`
	Isi string `json:"isi"`
	Tanggal string `json:"tanggal"`
}

func postHandler(ctx *gin.Context) {
	var newTask Task
	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message" : err.Error(),
		})
		return
	}
	
	if newTask.Isi != "" {
		_, err := DB.Exec("INSERT into tb_task (name, isi, tanggal) VALUES ($1,$2,$3)", newTask.Nama,newTask.Isi, newTask.Tanggal)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}
	ctx.JSON(http.StatusCreated,gin.H{"success":true, "pesan": "berhasil"})
}

func getAll(ctx *gin.Context){
	
	var tasks []Task
	var val Task
	rows, err := DB.Query("select * from tb_task")
	
	if err != nil {
		log.Fatalln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":false,
			"message":err,
		})
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&val.ID,&val.Nama,&val.Isi,&val.Tanggal)
		tasks = append(tasks, val)
		
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success" : true,
		"data" : tasks,
	})
}
