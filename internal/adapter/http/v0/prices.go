package v0

import (
	"io"
	"net/http"
	"supermarket/internal/adapter/http/v0/dto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) initPricesRoutes(api *gin.RouterGroup) {
	pricesGroup := api.Group("/prices")
	{
		pricesGroup.GET("", h.getPrices)
		pricesGroup.POST("", h.createPrices)
	}
}

// @Summary getPrices
// @Tags prices
// @Description get prices
// @Success 200 {object} zip-archive
// @Router /api/v0/prices  [get]
func (h *Handler) getPrices(context *gin.Context) {
	log.Info("Handler::getPrices start")
	products, err := h.productService.GetAllProducts(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	log.Info("Handler::getPrices start generation")
	csvData, err := h.parserService.GenerateCsv(products)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSV"})
		return
	}

	log.Info("Handler::getPrices start zip creation")
	zipData, err := h.parserService.CreateZipFile(csvData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file"})
		return
	}

	log.Info("Handler::getPrices success")
	context.Header("Content-Type", "application/zip")
	context.Header("Content-Disposition", "attachment; filename=data.zip")
	context.Data(http.StatusOK, "application/zip", zipData)
}

// @Summary createPrices
// @Tags prices
// @Description post prices
// @Param input body zip-archive true "user info"
// @Success 200 {object} dto.TotalDTO
// @Router /api/v0/prices  [post]
func (h *Handler) createPrices(context *gin.Context) {
	log.Info("Handler::createPrices start")
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	log.Info("Handler::createPrices open")
	uploadedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer uploadedFile.Close()

	log.Info("Handler::createPrices read")
	fileBytes, err := io.ReadAll(uploadedFile)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	log.Info("Handler::createPrices parse zip")
	products, err := h.parserService.ParseZip(context.Request.Context(), file.Filename, fileBytes)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse zip file"})
		return
	}

	log.Info("Handler::createPrices create batch")
	total, err := h.productService.CreateBatch(context.Request.Context(), products)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create products"})
		return
	}

	log.Info("Handler::createPrices create DTO")
	totalDTO := dto.NewTotalDTO(total)
	log.Info("Handler::createPrices success")
	context.JSON(http.StatusOK, totalDTO)
}
