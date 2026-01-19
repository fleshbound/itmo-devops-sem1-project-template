package app

import (
	v0 "supermarket/internal/adapter/http/v0"
	postgresRepo "supermarket/internal/adapter/repository/postgres"
	"supermarket/internal/app/config"
	"supermarket/internal/app/server"
	"supermarket/internal/core/service"
	"supermarket/pkg/database/postgres"
	"supermarket/pkg/logrus"

	log "github.com/sirupsen/logrus"
)

func RunWeb() {
	cfg := config.GetConfig()
	logrus.InitLogrus(&cfg.Logger)
	log.Info("application startup...")

	db, err := postgres.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		log.WithFields(log.Fields{
			"from":    "main()",
			"problem": "NewPostgresDB",
		}).Fatal(err.Error())
	}

	productRepo := postgresRepo.NewProductRepo(db)
	productService := service.NewProductService(productRepo)
	parserService := service.NewParserService()

	handlerParams := v0.HandlerParams{
		Config:         &cfg.Web,
		ProductService: productService,
		ParserService:  parserService,
	}
	gin := server.NewGinRouter()
	v0 := v0.NewHandler(handlerParams, gin)

	serverParams := server.ServerParams{
		Cfg:     &cfg.Web,
		Handler: v0,
		Router:  gin,
	}
	server := server.NewServer(serverParams)

	log.Info("app started")
	log.Fatal("server shutdown", server.ListenAndServe())
}
