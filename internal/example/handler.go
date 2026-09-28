package example

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler atbild par HTTP pieprasījumu saņemšanu, validāciju un atbilžu sūtīšanu.
// Tas izmanto Service, lai izpildītu biznesa loģiku.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes ir galvenā vieta, kur Handler reģistrē savus maršrutus pie padotās Gin grupas.
// Šādā veidā maršrutētājs (Router) paliek tīrs, un Handleris pats nosaka savus URL ceļus.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/examples")
	{
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.POST("", h.Create)
	}
}

// List atgriež visus ierakstus: GET /examples
func (h *Handler) List(c *gin.Context) {
	items, err := h.service.ListItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neizdevās iegūt ierakstus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// Get atgriež konkrētu ierakstu pēc ID: GET /examples/:id
func (h *Handler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nederīgs ID parametrs"})
		return
	}

	item, err := h.service.GetItem(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ieraksts netika atrasts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Create izveido jaunu ierakstu: POST /examples
func (h *Handler) Create(c *gin.Context) {
	var dto CreateExampleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateItem(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ieraksts veiksmīgi izveidots",
		"data":    created,
	})
}
