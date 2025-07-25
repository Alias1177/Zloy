package main

import (
	"log"
	"net/http"

	"github.com/Alias1177/Zloy/internal/config"
	handlers "github.com/Alias1177/Zloy/internal/delivery/http"
	"github.com/Alias1177/Zloy/internal/delivery/middleware"
	"github.com/Alias1177/Zloy/internal/repository"
	"github.com/Alias1177/Zloy/internal/usecase"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Инициализация репозиториев
	userRepo, err := repository.NewPostgresUserRepository(cfg)
	if err != nil {
		log.Fatal("Failed to initialize PostgreSQL repository:", err)
	}

	reportRepo, err := repository.NewMongoReportRepository(cfg)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB repository:", err)
	}

	// Инициализация use cases
	authUC := usecase.NewAuthUseCase(cfg.JWT.Secret, cfg.JWT.ExpirationTime)
	captchaUC := usecase.NewCaptchaUseCase(cfg.Captcha.Width, cfg.Captcha.Height, cfg.Captcha.NoiseCount, cfg.Captcha.SessionLifetime)
	userUC := usecase.NewUserUseCase(userRepo, authUC)
	reportUC := usecase.NewReportUseCase(reportRepo, userUC, cfg.Business.ReportCostCents)

	// Инициализация handlers
	authHandler := handlers.NewAuthHandler(userUC, captchaUC)
	reportHandler := handlers.NewReportHandler(reportUC, userUC, cfg.Server.DefaultPageSize)
	captchaHandler := handlers.NewCaptchaHandler(captchaUC)
	userHandler := handlers.NewUserHandler(userUC)

	// Инициализация middleware
	authMiddleware := middleware.AuthMiddleware(authUC)

	// Настройка роутера
	r := chi.NewRouter()

	// Добавляем стандартные middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// API роуты
	r.Route("/api", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		// User routes (protected)
		r.Route("/user", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/link-anonymous", reportHandler.LinkAnonymous)
			r.Get("/reports", reportHandler.GetUserReports)
			r.Post("/topup", userHandler.TopUpBalance)
			r.Get("/balance", userHandler.GetBalance)
		})

		// Reports routes (protected)
		r.Route("/reports", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/{report_id}/purchase", reportHandler.PurchaseReport)
		})

		// Mock routes
		r.Route("/mock", func(r chi.Router) {
			r.Post("/create-report", reportHandler.CreateMockReport)
		})

		// Captcha routes
		r.Route("/captcha", func(r chi.Router) {
			r.Get("/generate", captchaHandler.GenerateCaptcha)
		})
	})

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
