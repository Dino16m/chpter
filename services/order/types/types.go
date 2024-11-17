package types

import (
	"context"

	"github.com/dino16m/chpter/order/models"
	"github.com/dino16m/chpter/order/rpc"
)

type OrderRepository interface {
	Create(order *models.Order) error
	WithTransaction(cb func(OrderRepository) error) error
}

type UserService interface {
	GetUserById(id uint, ctx context.Context) (*rpc.UserMessage, error)
}
