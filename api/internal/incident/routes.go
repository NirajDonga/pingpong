package incident

import "github.com/gin-gonic/gin"

func RegisterRoutes(protected *gin.RouterGroup, h *Handler) {
	protected.GET("/incidents", h.List)
	protected.GET("/monitors/:id/incidents", h.ListForMonitor)
}
