package entity

import (
	"strconv"
	"supermarket/internal/core/model"
	"time"
)

type PgProduct struct {
	Id          int64     `db:"id"`
	Create_date time.Time `db:"create_date"`
	Name        string    `db:"name"`
	Category    string    `db:"category"`
	PriceCents  int64     `db:"price"`
}

func (p *PgProduct) ToModel() model.Product {
	return model.Product{
		Id:          p.Id,
		Create_date: p.Create_date,
		Name:        p.Name,
		Category:    p.Category,
		PriceCents:  p.PriceCents,
		PriceStr:    strconv.FormatFloat(float64(p.PriceCents)/100.0, 'f', 2, 64),
	}
}

func NewPgProduct(p model.Product) PgProduct {
	return PgProduct{
		Id:          p.Id,
		Create_date: p.Create_date,
		Name:        p.Name,
		Category:    p.Category,
		PriceCents:  p.PriceCents,
	}
}
