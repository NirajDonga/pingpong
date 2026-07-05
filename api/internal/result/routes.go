package result

import (
	ws "github.com/NirajDonga/pingpong/api/internal/websocket"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *Handler, wsManager *ws.Manager) {
	protected.GET("/monitors/:id/checks", h.History)
	protected.GET("/monitors/:id/ws", func(c *gin.Context) {
		h.StreamHistory(c, wsManager)
	})
}
