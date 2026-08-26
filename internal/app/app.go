package app

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-boilerplate/internal/auth"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/database"
	"go-boilerplate/internal/example"
	apphttp "go-boilerplate/internal/http"
	"go-boilerplate/internal/user"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

type Deps struct {
	Config config.Config
	Logger *slog.Logger
}

func New(deps Deps) (*App, error) {
	if deps.Config.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	db, err := database.NewPostgres(deps.Config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	jwtManager, err := auth.NewJWTManager(deps.Config.JWTSecret, deps.Config.JWTIssuer, deps.Config.JWTAccessTokenTTL)
	if err != nil {
		_ = database.Close(db)
		return nil, err
	}

	userRepository := user.NewRepository(db)
	authService := auth.NewService(userRepository, jwtManager)
	authHandler := auth.NewHandler(authService)

	exampleService := example.NewService()
	exampleHandler := example.NewHandler(exampleService)

	router := apphttp.NewRouter(apphttp.RouterDeps{
		Logger:               deps.Logger,
		CORSAllowedOrigins:   deps.Config.CORSAllowedOrigins,
		CORSAllowCredentials: deps.Config.CORSAllowCredentials,
		JWTManager:           jwtManager,
		ExampleHandler:       exampleHandler,
		AuthHandler:          authHandler,
	})

	return &App{Router: router, DB: db}, nil
}

func (a *App) Close() error {
	return database.Close(a.DB)
}
