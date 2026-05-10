package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, protected *gin.RouterGroup, h *Handler) {
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	protected.GET("/me", h.Me)
}
