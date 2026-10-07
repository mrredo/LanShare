package files

import (
	"errors"
	"fmt"
	"lanshare/config"
	"lanshare/internal/auth"
	"lanshare/pkg/httputil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	authHandler *auth.Handler
}

func NewHandler(service *Service, authHandler *auth.Handler) *Handler {
	return &Handler{
		service:     service,
		authHandler: authHandler,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, fe *gin.RouterGroup) {
	group := rg.Group("/files")
	group.Use()
	{
		group.POST("/upload", h.UploadFile)
		group.GET("/", h.FileList)
		group.GET("/:file/download", h.DownloadFile)
		group.GET("/:file", h.FileInfo)
		group.DELETE("/:file", h.DeleteFile)
	}
	fe.GET("/share")
	fe.GET("/files/:file")
	fe.GET("/")
}
func (h *Handler) FileInfo(c *gin.Context) {
	fileID := c.Param("file")

	file, err := h.service.GetById(fileID)
	response404 := gin.H{
		"error": "Fails netika atrasts",
	}
	if err != nil || file.ID == "" {
		c.JSON(http.StatusNotFound, response404)
		return
	}
	if file.IsExpired() {
		h.service.DeleteById(fileID)

		c.JSON(http.StatusNotFound, response404)
		return
	}

	c.JSON(200, gin.H{
		"file": file,
	})
}

func (h *Handler) FileList(c *gin.Context) {
	getUserFiles := httputil.GetInputDefault(c, "my_files", "true") == "true"
	getPublicFiles := httputil.GetInputDefault(c, "public_files", "true") == "true"

	limitFilesStr, lfOk := httputil.GetInput(c, "limit")
	limitUserFilesStr, lufOk := httputil.GetInput(c, "limit_my_files")
	limitFiles, err := strconv.Atoi(limitFilesStr)
	if err != nil || !lfOk {
		limitFiles = -1
	}
	limitUserFiles, err := strconv.Atoi(limitUserFilesStr)
	if err != nil || !lufOk {
		limitFiles = -1
	}

	publicFiles := make([]File, 0)

	response := gin.H{}

	if getPublicFiles {
		publicFiles, err = h.service.FindFileList(limitFiles)
		if err != nil {
			response["public_files_error"] = "Nevarēja atrast publiskos failus."
		} else {
			response["public_files"] = FilesToFileResponsefunc(publicFiles)
		}
	}

	userFiles := make([]File, 0)

	if getUserFiles {
		cookie, err := c.Cookie(config.SessionCookie)
		if err != nil {
			response["user_files_error"] = "Lietotājs nav reģistrēts sistēmā."
		} else {
			userFiles, err = h.service.FindFileListForUser(cookie, limitUserFiles)
			if err != nil {
				response["user_files_error"] = "Nevarēja atrast lietotāja failus."
			} else {
				response["user_files"] = FilesToFileResponsefunc(userFiles)
			}
		}
	}

	c.JSON(200, response)
}

func (h *Handler) UploadFile(c *gin.Context) {
	var dto UploadDTO

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ownerCookie, err := c.Cookie(config.SessionCookie)
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
		fmt.Sprintf(`attachment; filename="%s"`, *file.Filename),
	)

	c.File(*file.StoragePath)
}

func (h *Handler) DeleteFile(c *gin.Context) {
	ownerCookie, _ := c.Cookie(config.SessionCookie)

	adminCookie, _ := c.Cookie(config.AdminSessionCookie)

	fileParam := c.Param("file")
	file, err := h.service.GetById(fileParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Fails netika atrasts",
		})
		return
	}

	// pārbaudām, vai pieder fails (Ja ir admin cookie, šis neskaitās)
	if ownerCookie != file.OwnerCookie && adminCookie != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Tev nepieder šis fails!",
		})
		return
	}
	// pārbaudām, vai nepieder fails un vai ir admin
	if ownerCookie != file.OwnerCookie && !h.service.authService.IsAdminSessionValid(adminCookie) {
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
