package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	api "bookem-user-service/api"
	domain "bookem-user-service/domain"
	repo "bookem-user-service/repo"
	service "bookem-user-service/service"

	"github.com/gin-gonic/gin"
)

func Add(a int, b int) int {
	return a + b
}

var (
	server *gin.Engine
	dB     *gorm.DB
	rawDB  *sql.DB
)

func syncDatabase() {
	dB.AutoMigrate(&domain.User{})
}

func connectToDb() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	dB = db
	rawDB, _ = db.DB()

	log.Printf("Connected to DB!")
}

func main() {
	connectToDb()
	defer rawDB.Close()
	syncDatabase()

	server = gin.Default()

	repo := repo.NewRepository(dB)
	service := service.NewService(repo)
	handler := api.NewHandler(service)
	route := *api.NewRoute(handler)

	server.GET("/ping", func(c *gin.Context) {
		err := rawDB.Ping()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pong"})
	})

	rg := server.Group("/api")
	route.Route(rg)

	server.Run()
}
