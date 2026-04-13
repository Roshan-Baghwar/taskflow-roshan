package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/config"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/handler"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/middleware"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/repository"
	"github.com/Roshan-Baghwar/taskflow-roshan/backend/internal/service"
)

func main() {
	cfg := config.Load()
	middleware.SetJWTSecret(cfg.JWTSecret)

	// DB connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("DB connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations on startup
	if err := runMigrations(db); err != nil {
		slog.Warn("Migration warning", "error", err)
	}

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	r := gin.Default()
	r.Use(gin.Recovery())

	// Auth routes
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// Protected routes
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())

	// Projects
	projects := api.Group("/projects")
	projects.GET("", h.GetProjects)
	projects.POST("", h.CreateProject)
	projects.GET("/:id", h.GetProject)
	projects.PATCH("/:id", h.UpdateProject)
	projects.DELETE("/:id", h.DeleteProject)
	projects.GET("/:id/stats", h.GetProjectStats) // bonus

	// Tasks
	projects.GET("/:id/tasks", h.GetProjectTasks)
	projects.POST("/:id/tasks", h.CreateTask)
	api.PATCH("/tasks/:id", h.UpdateTask)
	api.DELETE("/tasks/:id", h.DeleteTask)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	go func() {
		slog.Info("Server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")
}

func runMigrations(db *sqlx.DB) error {
	// sqlx.DB has a DB() method that returns the underlying *sql.DB
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}
	defer m.Close()
	return m.Up()
}