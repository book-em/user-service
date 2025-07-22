package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func Add(a int, b int) int {
	return a + b
}

var DB *sql.DB

func connectToDb() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	DB_URL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	DB = db

	log.Printf("Connected to DB!")
}

func main() {
	connectToDb()
	defer DB.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		err := DB.Ping()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pong"})
	})

	r.Run()
}
