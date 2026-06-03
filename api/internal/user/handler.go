package user

import (
	"net/http"
	"time"

	"github.com/NirajDonga/pingpong/api/internal/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc          Service
	cookieSecure bool
}

func NewHandler(svc Service, cookieSecure bool) *Handler {
	return &Handler{
		svc:          svc,
		cookieSecure: cookieSecure,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	result, err := h.svc.Register(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.setSessionCookie(c, result.Token)
	c.JSON(http.StatusCreated, result.User)
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	result, err := h.svc.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setSessionCookie(c, result.Token)
	c.JSON(http.StatusOK, result.User)
}

func (h *Handler) Logout(c *gin.Context) {
	h.setCookieSameSite(c)
	c.SetCookie(auth.SessionCookieName, "", -1, "/", "", h.cookieSecure, true)
	c.Status(http.StatusNoContent)
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": userID})
}

func (h *Handler) setSessionCookie(c *gin.Context, token string) {
	h.setCookieSameSite(c)
	c.SetCookie(
		auth.SessionCookieName,
		token,
		int((24 * time.Hour).Seconds()),
		"/",
		"",
		h.cookieSecure,
		true,
	)
}

func (h *Handler) setCookieSameSite(c *gin.Context) {
	if h.cookieSecure {
		c.SetSameSite(http.SameSiteNoneMode)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
}
