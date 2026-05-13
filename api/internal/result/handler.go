package result

import (
	"net/http"
	"strconv"

	"github.com/NirajDonga/pingpong/api/internal/monitor"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	monitorSvc monitor.Service
	resultSvc  *Service
}

func NewHandler(monitorSvc monitor.Service, resultSvc *Service) *Handler {
	return &Handler{
		monitorSvc: monitorSvc,
		resultSvc:  resultSvc,
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

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number"})
		return
	}

	results, err := h.resultSvc.History(c.Request.Context(), monitorID, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
