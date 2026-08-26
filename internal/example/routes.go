package example

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	routes := router.Group("/examples")
	routes.POST("", h.Create)
	routes.GET("", h.List)
	routes.GET("/page", h.Paginate)
}
