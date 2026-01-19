package v0

import (
	"supermarket/internal/core/port"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port string
}

type Handler struct {
	config         *Config
	productService port.IProductService
	parserService  port.IParserService
}

type HandlerParams struct {
	Config         *Config
	ProductService port.IProductService
	ParserService  port.IParserService
}

func NewHandler(params HandlerParams, router *gin.Engine) *Handler {
	handler := &Handler{
		config:         params.Config,
		productService: params.ProductService,
		parserService:  params.ParserService,
	}

	api := router.Group("/api")
	v0 := api.Group("/v0")
	v0.Use(LoggerMiddleware())
	{
		handler.initPricesRoutes(v0)
	}

	return handler
}
