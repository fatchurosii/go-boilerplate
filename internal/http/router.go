package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"go-boilerplate/internal/auth"
	"go-boilerplate/internal/example"
	"go-boilerplate/internal/http/middleware"
	"go-boilerplate/internal/http/response"
)

type RouterDeps struct {
	Logger               *slog.Logger
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool
	JWTManager           *auth.JWTManager
	ExampleHandler       *example.Handler
	AuthHandler          *auth.Handler
}

func NewRouter(deps RouterDeps) *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Logger(deps.Logger),
		middleware.Recovery(deps.Logger),
		middleware.CORS(deps.CORSAllowedOrigins, deps.CORSAllowCredentials),
	)

	api := router.Group("/api")
	api.GET("/health", health)
	deps.AuthHandler.RegisterRoutes(api, middleware.Auth(deps.JWTManager))
	deps.ExampleHandler.RegisterRoutes(api)
	router.NoRoute(notFound)

	return router
}

func health(c *gin.Context) {
	response.Success(c, nethttp.StatusOK, "success", gin.H{"status": "ok"})
}

func notFound(c *gin.Context) {
	response.Error(c, response.NotFound("route not found"))
}
