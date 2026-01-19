package model

import (
	"time"
)

type Product struct {
	Id          int64
	Create_date time.Time
	Name        string
	Category    string
	PriceCents  int64
	PriceStr    string
}
