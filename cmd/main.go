package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/plandstu/internal/config"
	"github.com/immxrtalbeast/plandstu/internal/controller"
	"github.com/immxrtalbeast/plandstu/internal/domain"
	"github.com/immxrtalbeast/plandstu/internal/middleware"
	"github.com/immxrtalbeast/plandstu/internal/usecase/user"
	"github.com/immxrtalbeast/plandstu/storage/psql"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("starting application", slog.Any("config", cfg))
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("postgresql://postgres.rcrqslgziljieocjppwc:%s@aws-0-eu-central-1.pooler.supabase.com:5432/postgres", os.Getenv("DB_PASS"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{}, &domain.History{}, &domain.Message{})
	if err := db.Exec("DEALLOCATE ALL").Error; err != nil {
		panic(err)
	}
	usrRepo := psql.NewUserRepository(db)
	userINT := user.NewUserInteractor(usrRepo, cfg.TokenTTL, os.Getenv("APP_SECRET"))
	userController := controller.NewUserController(userINT, cfg.TokenTTL, os.Getenv("APP_SECRET"))

	parserController := controller.NewParserController(os.Getenv("PARSER_URL"))
	authMiddleware := middleware.AuthMiddleware(os.Getenv("APP_SECRET"))

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{
		"Authorization",
		"Content-Type",
		"Origin",
		"Accept",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.ExposeHeaders = []string{"Set-Cookie"}
	router.Use(cors.New(config))
	api := router.Group("/api/v1")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)
	}
	parser := api.Group("/parser")
	parser.Use(authMiddleware)
	{
		parser.GET("/faculties", parserController.Faculties)
		parser.GET("/faculties/:id", parserController.FacultyByID)
		parser.GET("/disciplines/:direction", parserController.Disciplines)
		parser.GET("/roadmap/:discipline/:link", parserController.Roadmap)
	}
	router.Run(":8080")
}
func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return log
}
