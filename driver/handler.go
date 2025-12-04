package driver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/drivers", h.CreateDriver)
	r.GET("/drivers", h.ListDrivers)
	r.PUT("/drivers/:id", h.UpdateDriver)
	r.GET("/drivers/nearby", h.NearbyDrivers)
}

// CreateDriver godoc
// @Summary      Driver oluştur
// @Description  Yeni bir driver kaydı oluşturur
// @Tags         Drivers
// @Accept       json
// @Produce      json
// @Param        driver  body      CreateDriverRequest  true  "Driver info"
// @Success      201     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /drivers [post]
func (h *Handler) CreateDriver(c *gin.Context) {
	var req CreateDriverRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	id, err := h.Service.CreateDriver(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

// ListDrivers godoc
// @Summary      Driver listesi
// @Description  Sistemdeki driver’ları listeler
// @Tags         Drivers
// @Produce      json
// @Param        page      query  int  false  "Page number"
// @Param        pageSize  query  int  false  "Page size"
// @Success      200  {array}  Driver
// @Failure      500  {object} map[string]string
// @Router       /drivers [get]
func (h *Handler) ListDrivers(c *gin.Context) {
	// page & pageSize şimdilik sadece okunuyor, service tarafında tüm liste dönüyoruz
	_, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	_, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	list, err := h.Service.ListDrivers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list drivers"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// UpdateDriver godoc
// @Summary      Driver güncelle
// @Description  ID ile driver bilgisini günceller
// @Tags         Drivers
// @Accept       json
// @Produce      json
// @Param        id      path   string               true  "Driver ID"
// @Param        driver  body   UpdateDriverRequest  true  "Driver info"
// @Success      200  {object}  map[string]bool
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /drivers/{id} [put]
func (h *Handler) UpdateDriver(c *gin.Context) {
	id := c.Param("id")

	var req UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.Service.UpdateDriver(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated": true})
}

// NearbyDrivers godoc
// @Summary      Yakındaki taksiler
// @Description  Girilen konuma 6 km yarıçapındaki taksileri listeler
// @Tags         Drivers
// @Produce      json
// @Param        lat       query  number  true   "Latitude"
// @Param        lon       query  number  true   "Longitude"
// @Param        taksiType query  string  false  "Taxi type"
// @Success      200  {array}  NearbyDriver
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /drivers/nearby [get]
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

	list, err := h.Service.FindNearbyDrivers(c.Request.Context(), lat, lon, taxiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "nearby driver aranırken hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, list)
}
