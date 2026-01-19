package dto

import (
	"supermarket/internal/core/model"
)

type TotalDTO struct {
	Items_count      int64   `json:"total_items"`
	Categories_count int64   `json:"total_categories"`
	Price_sum        float64 `json:"total_price"`
}

func NewTotalDTO(total model.Total) *TotalDTO {
	return &TotalDTO{
		Items_count:      total.Items_count,
		Categories_count: total.Categories_count,
		Price_sum:        total.Price_sum,
	}
}
