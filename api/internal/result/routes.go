package result

import "github.com/gin-gonic/gin"

func RegisterRoutes(protected *gin.RouterGroup, h *Handler) {
	protected.GET("/monitors/:id/checks", h.History)
}
