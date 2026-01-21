package postgres

import (
	"context"
	"database/sql"
	"supermarket/internal/adapter/repository/postgres/entity"
	"supermarket/internal/core/model"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PostgresProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *PostgresProductRepo {
	return &PostgresProductRepo{
		db: db,
	}
}

const (
	productGetQuery              = "SELECT id, create_date, name, category, price FROM prices"
	productGetPriceSumQuery      = "SELECT COALESCE(SUM(price), 0) FROM prices"
	productGetCategoryCountQuery = "SELECT COUNT(DISTINCT category) FROM prices"
	productInsertBatchQuery      = "INSERT INTO prices (create_date, name, category, price) VALUES ($1, $2, $3, $4)"
)

func (repo *PostgresProductRepo) Get(ctx context.Context) ([]model.Product, error) {
	var pgProducts []entity.PgProduct
	err := repo.db.SelectContext(ctx, &pgProducts, productGetQuery)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(model.ErrPersistenceFailed, err.Error())
	}

	products := make([]model.Product, len(pgProducts))
	for i, p := range pgProducts {
		products[i] = p.ToModel()
	}
	return products, nil
}

func (repo *PostgresProductRepo) CreateBatch(ctx context.Context, products []model.Product) (model.Total, error) {

	total := model.Total{}

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return total, err
	}
	defer tx.Rollback()

	for _, product := range products {
		pgProduct := entity.NewPgProduct(product)
		query := productInsertBatchQuery
		_, err := tx.ExecContext(ctx, query,
			pgProduct.Create_date,
			pgProduct.Name,
			pgProduct.Category,
			float64(pgProduct.PriceCents)/100.0)
		if err != nil {
			return total, err
		}
	}

	total.Items_count = int64(len(products))
	tx.GetContext(ctx, &total.Categories_count, productGetCategoryCountQuery)
	tx.GetContext(ctx, &total.Price_sum, productGetPriceSumQuery)

	tx.Commit()
	return total, nil
}
