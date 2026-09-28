package monitor

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	m, err := h.svc.Create(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	monitors, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, monitors)
}

func (h *Handler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	m, err := h.svc.Get(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		writeMonitorError(c, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	m, err := h.svc.Update(c.Request.Context(), userID, c.Param("id"), input)
	if err != nil {
		writeMonitorError(c, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *Handler) SetEnabled(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	var input EnabledRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	m, err := h.svc.SetEnabled(c.Request.Context(), userID, c.Param("id"), input.Enabled)
	if err != nil {
		writeMonitorError(c, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	err := h.svc.Delete(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		writeMonitorError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func writeMonitorError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
