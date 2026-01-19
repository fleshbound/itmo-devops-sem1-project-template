package service

import (
	"context"
	"supermarket/internal/core/model"
	"supermarket/internal/core/port"

	log "github.com/sirupsen/logrus"
)

type ProductService struct {
	productRepo port.IProductRepository
}

func NewProductService(productRepo port.IProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (p *ProductService) CreateBatch(ctx context.Context, params []port.CreateProductParam) (model.Total, error) {
	productsDTO := make([]model.Product, len(params))
	for i, param := range params {
		productsDTO[i] = model.Product{
			Id:          param.Id,
			Create_date: param.Create_date,
			Name:        param.Name,
			Category:    param.Category,
			PriceCents:  param.PriceCents,
			PriceStr:    param.PriceStr,
		}
	}

	total, err := p.productRepo.CreateBatch(ctx, productsDTO)
	if err != nil {
		log.WithFields(log.Fields{
			"from": "ProductServiceCreateBatch",
		}).Error(err.Error())
		return model.Total{}, err
	}

	log.WithFields(log.Fields{
		"count": len(params),
	}).Info("ProductServiceCreateBatch OK")
	return total, nil
}

func (p *ProductService) GetAllProducts(ctx context.Context) ([]port.CreateProductParam, error) {
	log.Info("ProductService::GetAllProducts start")

	products, err := p.productRepo.Get(ctx)
	if err != nil {
		log.Error("ProductService::GetAllProducts error repo get")
		return nil, err
	}

	var params []port.CreateProductParam
	for _, product := range products {
		params = append(params, port.CreateProductParam{
			Id:          product.Id,
			Create_date: product.Create_date,
			Name:        product.Name,
			Category:    product.Category,
			PriceCents:  product.PriceCents,
			PriceStr:    product.PriceStr,
		})
	}

	log.Error("ProductService::GetAllProducts success")

	return params, nil
}
