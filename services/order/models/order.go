package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LineItem struct {
	gorm.Model
	ProductId uint
	Quantity  int64
	UnitPrice decimal.Decimal
	OrderId   uint
}

func (l *LineItem) Price() decimal.Decimal {
	return l.UnitPrice.Mul(decimal.NewFromInt(l.Quantity))
}

type Order struct {
	gorm.Model
	CustomerId uint
	LineItems  []LineItem `gorm:"foreignKey:OrderId"`
}

func (order Order) TotalPrice() decimal.Decimal {
	total := decimal.NewFromInt(0)
	for _, item := range order.LineItems {
		total = total.Add(item.Price())
	}

	return total
}
