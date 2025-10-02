package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	api "bookem-user-service/api"
	"bookem-user-service/client/roomclient"
	domain "bookem-user-service/domain"
	repo "bookem-user-service/repo"
	service "bookem-user-service/service"
	utils "bookem-user-service/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

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
	ctx := context.Background()
	// shutdown := utils.InitTracer(
	// 	ctx,
	// 	os.Getenv("SERVICE_NAME"),
	// 	os.Getenv("DEPLOYMENT_ENV"),
	// )
	// defer shutdown(ctx)

	shutdown2 := utils.TEL.Init(
		ctx,
		os.Getenv("SERVICE_NAME"),
		os.Getenv("DEPLOYMENT_ENV"),
	)
	defer shutdown2(ctx)

	connectToDb()
	defer rawDB.Close()
	syncDatabase()

	server = gin.Default()

	server.Use(otelgin.Middleware(os.Getenv("SERVICE_NAME")))
	server.Use(utils.TEL.GetLoggingMiddleware())
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "http://bookem.local"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.GET("/healthz", func(ctx *gin.Context) {
		err := rawDB.Ping()
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, "Database not reachable")
			return
		}
		ctx.JSON(http.StatusOK, nil)
	})

	roomclient := roomclient.NewRoomClient()

	repo := repo.NewRepository(dB)
	service := service.NewService(repo, roomclient)
	handler := api.NewHandler(service)
	route := *api.NewRoute(handler)

	rg := server.Group("/api")
	route.Route(rg)

	server.Run()
}
