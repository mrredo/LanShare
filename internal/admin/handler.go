package admin

import (
	"lanshare/pkg/httputil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const SessionId = "session_id"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.Login)
	rg.POST("/logout", h.Logout)

}
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Nav autorizēts",
			})
			return
		}

		if !h.service.IsSessionValid(sessionID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Sesija nav derīga",
			})
			return
		}

		c.Next()
	}
}

func (h *Handler) CreateSession(c *gin.Context) {
	session, err := h.service.CreateSession()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.SetCookie(
		SessionId,
		session.ID,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		true,
		true,
	)
	c.JSON(200, gin.H{"message": "Veiksmīga pieslēgšanās"})
}
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		SessionId,
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
	c.JSON(http.StatusOK, gin.H{
		"message": "Atslēgts",
	})
}
func (h *Handler) Login(c *gin.Context) {
	password := httputil.GetInputDefault(c, "password", "")
	if !h.service.settingsService.IsPasswordValid(password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parole nav pareiza."})
		return
	}
	h.CreateSession(c)
}
