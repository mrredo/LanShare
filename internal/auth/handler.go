package auth

import (
	"lanshare/config"
	"lanshare/pkg/httputil"
	"lanshare/pkg/str"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
		sessionID, err := c.Cookie(config.AdminSessionCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Nav autorizēts",
			})
			return
		}

		if !h.service.IsAdminSessionValid(sessionID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Sesija nav derīga",
			})
			return
		}

		c.Next()
	}
}

// UserSessionMiddleware apply to all routes, so the user no matter what has a session
func (h *Handler) UserSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(config.SessionCookie)
		if err != nil || sessionID == "" {
			sessionValue, err := str.GenerateRandomString(32)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.SetCookie(config.SessionCookie, sessionValue, config.SessionCookieAge, "/", "", config.CookiesOnHTTPSOnly, true)
		}
		c.Next()
	}
}

func (h *Handler) CreateSession(c *gin.Context) {
	session, err := h.service.CreateSession()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie(
		config.AdminSessionCookie,
		session.ID,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		config.CookiesOnHTTPSOnly,
		true,
	)
	c.JSON(200, gin.H{"message": "Veiksmīga pieslēgšanās"})
}
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		config.AdminSessionCookie,
		"",
		-1,
		"/",
		"",
		config.CookiesOnHTTPSOnly,
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
