package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/plandstu/internal/config"
	"github.com/immxrtalbeast/plandstu/internal/controller"
	"github.com/immxrtalbeast/plandstu/internal/domain"
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
	userINT := user.NewUserInteractor(usrRepo, cfg.TokenTTL, cfg.AppSecret)
	userController := controller.NewUserController(userINT, cfg.TokenTTL, cfg.AppSecret)
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)
	}
	router.Run(":8000")
}
func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return log
}
