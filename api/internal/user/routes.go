package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, protected *gin.RouterGroup, h *Handler) {
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.POST("/logout", h.Logout)

	protected.GET("/me", h.Me)
}
