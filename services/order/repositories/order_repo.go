package repositories

import (
	"github.com/dino16m/chpter/order/apperrs"
	"github.com/dino16m/chpter/order/models"
	"github.com/dino16m/chpter/order/types"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) WithTransaction(cb func(types.OrderRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		repo := &OrderRepository{db: tx}
		return cb(repo)
	})
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindById(id uint) (models.Order, error) {
	var order models.Order
	tx := r.db.Joins("LineItems").First(&order, id)
	return apperrs.WrapNotFound(order, tx.Error, "Order not found")
}
