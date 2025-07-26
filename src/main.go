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
	DB     *gorm.DB
	RawDB  *sql.DB
)

func syncDatabase() {
	DB.AutoMigrate(&domain.User{})
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

	DB = db
	RawDB, _ = db.DB()

	log.Printf("Connected to DB!")
}

func main() {
	connectToDb()
	defer RawDB.Close()
	syncDatabase()

	server = gin.Default()

	userRepo := repo.NewUserRepository(DB)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)
	userRoute := *api.NewUserRoute(userHandler)

	server.GET("/ping", func(c *gin.Context) {
		err := RawDB.Ping()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pong"})
	})

	router := server.Group("/api")
	userRoute.UserRoute(router)

	server.Run()
}
