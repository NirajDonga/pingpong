package result

import (
	"log"
	"net/http"
	"strconv"

	"github.com/NirajDonga/pingpong/backend/api/internal/monitor"
	ws "github.com/NirajDonga/pingpong/backend/api/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	monitorSvc monitor.Service
	resultRepo *TinybirdRepository
}

func NewHandler(monitorSvc monitor.Service, resultRepo *TinybirdRepository) *Handler {
	return &Handler{
		monitorSvc: monitorSvc,
		resultRepo: resultRepo,
	}
}

func (h *Handler) History(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	monitorID := c.Param("id")
	if _, err := h.monitorSvc.Get(c.Request.Context(), userID, monitorID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "monitor not found"})
		return
	}

	mid, err := uuid.Parse(monitorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monitor id"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number"})
		return
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	results, err := h.resultRepo.History(c.Request.Context(), mid, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) StreamHistory(c *gin.Context, wsManager *ws.Manager) {
	monitorID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws: upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	wsManager.AddClient(monitorID, conn)
	defer wsManager.RemoveClient(monitorID, conn)

	// Block until the client disconnects.
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
