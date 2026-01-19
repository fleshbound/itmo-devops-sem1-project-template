package server

import (
	"fmt"
	"net/http"
	v0 "supermarket/internal/adapter/http/v0"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(cors.Default())
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}

type ServerParams struct {
	Cfg     *v0.Config
	Handler *v0.Handler
	Router  *gin.Engine
}

func NewServer(params ServerParams) *http.Server {
	server := &http.Server{
		Handler:      params.Router,
		Addr:         fmt.Sprintf(":%s", params.Cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return server
}
