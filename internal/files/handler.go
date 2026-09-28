package files

import (
	"fmt"
	"lanshare/config"
	"net/http"

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
	fileParam := c.Param("file")
	// get file
	// delete file from db and storage
	// return success
	file, err := h.service.GetById(fileParam)
	if err != nil {

		return
	}
	h.ser

}
