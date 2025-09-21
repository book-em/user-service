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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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

func initTracer() func(context.Context) error {
	ctx := context.Background()

	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create an OTLP HTTP exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("user-service"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp.Shutdown
}

func main() {
	ctx := context.Background()
	shutdown := initTracer()
	defer shutdown(ctx)

	connectToDb()
	defer rawDB.Close()
	syncDatabase()

	server = gin.Default()

	server.Use(otelgin.Middleware("user-service"))
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
