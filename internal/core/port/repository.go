package port

import (
	"context"
	"supermarket/internal/core/model"
)

type IProductRepository interface {
	Get(ctx context.Context) ([]model.Product, error)
	CreateBatch(ctx context.Context, products []model.Product) (model.Total, error)
}
