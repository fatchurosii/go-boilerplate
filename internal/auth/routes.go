package auth

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(router gin.IRouter, authMiddleware gin.HandlerFunc) {
	routes := router.Group("/auth")
	routes.POST("/login", h.Login)
	routes.POST("/logout", authMiddleware, h.Logout)
	routes.GET("/me", authMiddleware, h.Me)
}
