package driver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/drivers", h.CreateDriver)
	r.GET("/drivers", h.ListDrivers)
	r.PUT("/drivers/:id", h.UpdateDriver)
	r.DELETE("/drivers/:id", h.DeleteDriver)
	r.GET("/drivers/nearby", h.FindNearbyDrivers)
}

func (h *Handler) CreateDriver(c *gin.Context) {
	var req CreateDriverRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.Service.CreateDriver(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) UpdateDriver(c *gin.Context) {
	id := c.Param("id")

	var req UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.Service.UpdateDriver(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "driver updated successfully",
	})
}

func (h *Handler) DeleteDriver(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.DeleteDriver(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListDrivers(c *gin.Context) {
	drivers, err := h.Service.ListDrivers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, drivers)
}

func (h *Handler) FindNearbyDrivers(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	taxiType := c.Query("taksiType")

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat/lon value"})
		return
	}

	result, err := h.Service.FindNearbyDrivers(c.Request.Context(), lat, lon, taxiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not search nearby drivers"})
		return
	}

	c.JSON(http.StatusOK, result)
}
