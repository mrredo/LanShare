package files

import (
	"errors"
	"fmt"
	"lanshare/config"
	"log"
	"net/http"
	"os"

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
	group := rg.Group("/files")
	group.POST("/upload", h.UploadFile)
	group.GET("/:file", h.DownloadFile)
	group.DELETE("/:file", h.DeleteFile)
}
func (h *Handler) UploadFile(c *gin.Context) {
	var dto UploadDTO

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ownerCookie, err := c.Cookie("lanshare_owner")
	if err != nil {
		c.JSON(500, gin.H{
			"error": "owner cookie not found",
		})
		return
	}
	file, err := h.service.CreateFile(&dto, ownerCookie)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success":   true,
		"message":   "Dati tika veiksmīgi augšupielādēti",
		"share_url": fmt.Sprintf("http://%s/share/%s", config.DomainName, file.ID),
	})
}

func (h *Handler) DownloadFile(c *gin.Context) {
	fileID := c.Param("file")

	file, err := h.service.GetById(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Fails netika atrasts",
		})
		return
	}

	if file.StoragePath == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Fails nav pieejams",
		})
		return
	}

	c.Header(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, file.Filename),
	)

	c.File(*file.StoragePath)
}

func (h *Handler) DeleteFile(c *gin.Context) {
	ownerCookie, err := c.Cookie(config.SessionCookie)
	if err != nil {
	}
	adminCookie, err := c.Cookie(config.AdminSessionCookie)
	if err != nil {
	}
	fileParam := c.Param("file")
	file, err := h.service.GetById(fileParam)
	if err != nil {
		// TODO:

		return
	}

	// pārbaudām, vai pieder fails (Ja ir admin cookie, šis neskaitās)
	if ownerCookie != file.OwnerCookie && adminCookie != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Tev nepieder šis fails!",
		})
		return
	}
	// pārbaudām, vai ir admin
	if !h.service.authService.IsAdminSessionValid(adminCookie) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Tu neesi administrators!",
		})
		return
	}

	if err := h.service.DeleteById(fileParam); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Faila dzēšana neizdevās.",
		})
		return

	}
	if file.StoragePath != nil {
		if err := os.Remove(*file.StoragePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Faila dzēšana neizdevās!",
			})
			return
		}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
		},
	)

}
