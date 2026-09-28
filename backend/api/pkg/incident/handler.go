package incident

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number"})
		return
	}

	incidents, err := h.repo.List(c.Request.Context(), uid, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func (h *Handler) ListForMonitor(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	mid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monitor id"})
		return
	}

	limit, err := parseLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a number"})
		return
	}

	incidents, err := h.repo.ListForMonitor(c.Request.Context(), uid, mid, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func parseLimit(c *gin.Context) (int, error) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil {
		return 0, err
	}
	if limit <= 0 || limit > 500 {
		return 100, nil
	}
	return limit, nil
}
