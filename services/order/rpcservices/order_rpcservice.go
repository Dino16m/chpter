package rpcservices

import (
	"context"
	"errors"
	"sync"

	"github.com/dino16m/chpter/order/apperrs"
	"github.com/dino16m/chpter/order/models"
	"github.com/dino16m/chpter/order/rpc"
	"github.com/dino16m/chpter/order/types"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderRPCService struct {
	orderRepository types.OrderRepository
	rpc.UnimplementedOrderRPCServiceServer

	userService types.UserService
}

type Result struct {
	Err   error
	Value any
}

func (r Result) Ok() bool {
	return r.Err == nil
}

func NewOrderRPCService(orderRepo types.OrderRepository, userService types.UserService) OrderRPCService {
	return OrderRPCService{orderRepository: orderRepo, userService: userService}
}

func toOrderMessage(order models.Order, user *rpc.UserMessage) *rpc.OrderMessage {
	var lineItems []*rpc.LineItemMessage

	for _, item := range order.LineItems {
		lineItem := &rpc.LineItemMessage{
			ProductId: int64(item.ProductId),
			Quantity:  item.Quantity,
			UnitPrice: &rpc.Decimal{Value: item.UnitPrice.String()},
			Price:     &rpc.Decimal{Value: item.Price().String()},
		}
		lineItems = append(lineItems, lineItem)
	}
	return &rpc.OrderMessage{
		Id:         int64(order.ID),
		CustomerId: int64(order.CustomerId),
		Customer:   user,
		Total:      &rpc.Decimal{Value: order.TotalPrice().String()},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: order.CreatedAt.Unix(),
			Nanos:   int32(order.CreatedAt.Nanosecond()),
		},
		Items: lineItems,
	}
}

func fromOrderMessage(message *rpc.NewOrderMessage) models.Order {
	var lineItems []models.LineItem

	for _, newItem := range message.Items {
		unitPrice, _ := decimal.NewFromString(newItem.UnitPrice.Value)
		item := models.LineItem{
			ProductId: uint(newItem.ProductId),
			Quantity:  newItem.Quantity,
			UnitPrice: unitPrice,
		}
		lineItems = append(lineItems, item)
	}
	return models.Order{
		CustomerId: uint(message.CustomerId),
		LineItems:  lineItems,
	}
}

func (s OrderRPCService) CreateOrder(ctx context.Context, in *rpc.NewOrderMessage) (*rpc.OrderMessage, error) {
	order, err := s.createOrderSync(ctx, in)

	if err != nil && status.Code(err) == codes.Unknown {
		return order, apperrs.NewServerError("An error occured while creating a new order", err)
	}

	return order, err
}

func (s OrderRPCService) createOrderSync(ctx context.Context, in *rpc.NewOrderMessage) (*rpc.OrderMessage, error) {
	var response *rpc.OrderMessage
	err := s.orderRepository.WithTransaction(func(repository types.OrderRepository) error {
		order := fromOrderMessage(in)
		err := repository.Create(&order)

		if err != nil {
			return err
		}
		user, err := s.userService.GetUserById(uint(in.CustomerId), ctx)
		if err != nil {
			return err
		}
		response = toOrderMessage(order, user)
		return nil
	})

	return response, err
}

func (s OrderRPCService) createOrderAsync(ctx context.Context, in *rpc.NewOrderMessage) (*rpc.OrderMessage, error) {
	results := make(chan Result, 2)
	channelContext, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	var response *rpc.OrderMessage
	err := s.orderRepository.WithTransaction(func(repository types.OrderRepository) error {
		wg.Add(1)
		go waiter(
			wg,
			channelContext,
			cancelFunc,
			results,
			func(res chan Result, ctx context.Context) {
				order := fromOrderMessage(in)
				err := repository.Create(&order)
				res <- Result{Value: order, Err: err}
				close(res)
			},
		)

		wg.Add(1)
		go waiter(
			wg,
			channelContext,
			cancelFunc,
			results,
			func(res chan Result, ctx context.Context) {
				user, err := s.userService.GetUserById(uint(in.CustomerId), ctx)
				res <- Result{Value: user, Err: err}
				close(res)
			},
		)

		wg.Wait()
		close(results)

		orderMessage, err := resultToOrderMessage(results)
		response = orderMessage
		return err
	})

	return response, err
}

func waiter(
	wg *sync.WaitGroup,
	ctx context.Context,
	cancelFunc func(),
	results chan<- Result,
	cb func(res chan Result, ctx context.Context),
) {
	defer wg.Done()

	output := make(chan Result)
	go cb(output, ctx)
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			return
		case result := <-output:
			results <- result
			if result.Err != nil {
				cancelFunc()
			}

			return
		}
	}

}

func resultToOrderMessage(results chan Result) (*rpc.OrderMessage, error) {
	var err error = nil
	var order models.Order
	var user *rpc.UserMessage
	for result := range results {
		if !result.Ok() {
			err = errors.Join(err, result.Err)
			continue
		}
		switch result.Value.(type) {
		case models.Order:
			order = result.Value.(models.Order)
		case *rpc.UserMessage:
			user = result.Value.(*rpc.UserMessage)
		}
	}
	if err != nil {
		return nil, err
	}
	return toOrderMessage(order, user), nil
}
