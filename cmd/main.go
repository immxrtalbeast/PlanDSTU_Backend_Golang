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
	"github.com/immxrtalbeast/plandstu/internal/task"
	"github.com/immxrtalbeast/plandstu/internal/usecase/llm"
	"github.com/immxrtalbeast/plandstu/internal/usecase/roadmap"
	"github.com/immxrtalbeast/plandstu/internal/usecase/tests"
	"github.com/immxrtalbeast/plandstu/internal/usecase/user"
	"github.com/immxrtalbeast/plandstu/internal/worker"
	"github.com/immxrtalbeast/plandstu/storage/psql"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// go run .\cmd\main.go --config=./config/local.yaml
func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	log.Info("starting application", slog.Any("config", cfg))
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("postgresql://postgres.rcrqslgziljieocjppwc:%s@aws-0-eu-central-1.pooler.supabase.com:6543/postgres", os.Getenv("DB_PASS"))
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Info("db connected")
	db.AutoMigrate(&domain.User{}, &domain.History{}, &domain.RoadmapHistory{}, &domain.RoadmapTest{})
	if err := db.Exec("DEALLOCATE ALL").Error; err != nil {
		panic(err)
	}
	authMiddleware := middleware.AuthMiddleware(os.Getenv("APP_SECRET"))
	usrRepo := psql.NewUserRepository(db)
	userINT := user.NewUserInteractor(usrRepo, cfg.TokenTTL, os.Getenv("APP_SECRET"))
	userController := controller.NewUserController(userINT, cfg.TokenTTL, os.Getenv("APP_SECRET"))
	LLMRepo := psql.NewLLMRepository(db)
	LLMINT := llm.NewLLMInteractor(LLMRepo)
	LLMController := controller.NewLLMController(os.Getenv("LLM_URL"), LLMINT)
	RoadmapRepo := psql.NewRoadmapRepository(db)

	TestRepository := psql.NewTestRepository(db)

	RoadmapINT := roadmap.NewRoadmapInteractor(RoadmapRepo, TestRepository)
	RoadmapController := controller.NewRoadmapController(RoadmapINT)

	TestINT := tests.NewTestInteractor(TestRepository, os.Getenv("LLM_URL"), RoadmapRepo)
	TestsController := controller.NewTestsController(os.Getenv("LLM_URL"), RoadmapINT, TestINT, os.Getenv("REDIS_URL"))
	parserController := controller.NewParserController(os.Getenv("PARSER_URL"))

	task.Init(os.Getenv("REDIS_URL"))
	worker := worker.NewWorker(os.Getenv("REDIS_URL"), 10, TestINT)
	go func() {
		if err := worker.Start(); err != nil {
			panic("worker failed")
		}
	}()
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
	llm := api.Group("/llm")
	llm.Use(authMiddleware)
	{
		llm.GET("/chat", LLMController.Chat)
		llm.GET("/history", LLMController.History)
		llm.GET("/clear_history", LLMController.ClearHistory)
	}
	api.POST("/llm/save-history", LLMController.SaveHistory)

	// api.GET("/roadmap/history/:link", RoadmapController.History).Use(authMiddleware)
	// api.POST("/roadmap/send-report").Use(authMiddleware)
	tests := api.Group("/tests")
	tests.Use(authMiddleware)
	{
		tests.GET("/history", RoadmapController.History)
		tests.POST("/first-test", TestsController.FirstTest)
		tests.POST("/answers", TestsController.Answers)
		tests.POST("/default-test", TestsController.CreateTest)
		tests.POST("/report", RoadmapController.Report)
		tests.GET("/status", TestsController.GetTaskStatus)
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
