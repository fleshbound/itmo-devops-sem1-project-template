package v0

import (
	"io"
	"net/http"
	"supermarket/internal/adapter/http/v0/dto"

	"github.com/gin-gonic/gin"
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
	products, err := h.productService.GetAllProducts(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	csvData, err := h.parserService.GenerateCsv(products)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSV"})
		return
	}

	zipData, err := h.parserService.CreateZipFile(csvData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file"})
		return
	}

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
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer uploadedFile.Close()

	fileBytes, err := io.ReadAll(uploadedFile)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	products, err := h.parserService.ParseZip(context.Request.Context(), file.Filename, fileBytes)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse zip file"})
		return
	}

	total, err := h.productService.CreateBatch(context.Request.Context(), products)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create products"})
		return
	}

	totalDTO := dto.NewTotalDTO(total)
	context.JSON(http.StatusOK, totalDTO)
}
