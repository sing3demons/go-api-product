package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sing3demons/app/database"
	"github.com/sing3demons/app/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	docs "github.com/sing3demons/app/docs"
	swagger "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

// @title Swagger GO-API-PRODUCT API
// @version 1.0
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host https://go-kpsing.herokuapp.com
// @BasePath /
func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	database.ConnectDB()
	// seeds.Load()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{os.Getenv("ENV_URL")}
	corsConfig.AddAllowHeaders("Authorization")

	// config.AllowOrigins = []string{
	// 	"http://localhost:8080",
	// }
	// config.AllowHeaders = []string{
	// 	"Origin",
	// 	"Authorization",
	// 	"TransactionID",
	// }

	r := gin.Default()
	r.Use(cors.New(corsConfig)) //cors
	r.Static("/uploads", "./uploads")

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagger.Handler))

	r.GET("/healthz", health)

	r.GET("/x", buildX)

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(r)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

// @Accept  json
// @Produce  json
// @Success 200
// @Router /healthz [get]
func health(c *gin.Context) {
	c.Status(200)
}

// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]any
// @Router /x [get]
func buildX(c *gin.Context) {
	c.JSON(200, gin.H{
		"build_commit": buildcommit,
		"build_time":   buildtime,
	})
}
