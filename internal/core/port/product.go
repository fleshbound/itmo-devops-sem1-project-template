package port

import (
	"context"
	"supermarket/internal/core/model"
	"time"
)

type CreateProductParam struct {
	Id          int64
	Create_date time.Time
	Name        string
	Category    string
	PriceCents  int64
	PriceStr    string
}

type IProductService interface {
	CreateBatch(ctx context.Context, param []CreateProductParam) (model.Total, error)
	GetAllProducts(ctx context.Context) ([]CreateProductParam, error)
}
