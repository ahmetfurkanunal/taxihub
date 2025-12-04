package driver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/drivers", h.CreateDriver)
	r.GET("/drivers", h.ListDrivers)
	r.PUT("/drivers/:id", h.UpdateDriver)
	r.GET("/drivers/nearby", h.NearbyDrivers)
}

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateDriver(c *gin.Context) {
	var req CreateDriverRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	id, err := h.Service.CreateDriver(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

func (h *Handler) ListDrivers(c *gin.Context) {
	list, err := h.Service.ListDrivers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list drivers"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateDriver(c *gin.Context) {
	id := c.Param("id")

	var req UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.Service.UpdateDriver(c, id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated": true})
}

func (h *Handler) NearbyDrivers(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	taxiType := c.Query("taksiType")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat ve lon zorunludur"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz lat"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz lon"})
		return
	}

	list, err := h.Service.FindNearbyDrivers(c, lat, lon, taxiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "nearby driver aranırken hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, list)
}
