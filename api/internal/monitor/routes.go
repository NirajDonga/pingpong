package monitor

import "github.com/gin-gonic/gin"

func RegisterRoutes(protected *gin.RouterGroup, h *Handler) {
	monitors := protected.Group("/monitors")

	monitors.POST("", h.Create)
	monitors.GET("", h.List)
	monitors.GET("/:id", h.Get)
	monitors.PUT("/:id", h.Update)
	monitors.PATCH("/:id/enabled", h.SetEnabled)
	monitors.DELETE("/:id", h.Delete)
}
