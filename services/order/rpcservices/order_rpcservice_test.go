package rpcservices

import (
	"context"
	"testing"
	"time"

	"github.com/dino16m/chpter/order/apperrs"
	"github.com/dino16m/chpter/order/models"
	"github.com/dino16m/chpter/order/rpc"
	"github.com/dino16m/chpter/order/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserService struct {
	err   error
	user  *rpc.UserMessage
	sleep time.Duration
}

func (s *MockUserService) GetUserById(id uint, ctx context.Context) (*rpc.UserMessage, error) {
	time.Sleep(s.sleep)
	return s.user, s.err
}

type MockOrderRepository struct {
	err              error
	orderId          uint
	sleep            time.Duration
	transactionCount int
}

func (s *MockOrderRepository) Create(order *models.Order) error {
	time.Sleep(s.sleep)
	order.ID = s.orderId
	return s.err
}

func (s *MockOrderRepository) WithTransaction(cb func(types.OrderRepository) error) error {
	s.transactionCount += 1
	return cb(s)
}

type OrderRPCTest struct {
	suite.Suite

	orderRPCService OrderRPCService
	userService     *MockUserService
	orderRepository *MockOrderRepository
}

func (t *OrderRPCTest) SetupTest() {
	t.orderRepository = &MockOrderRepository{}
	t.userService = &MockUserService{}

	t.orderRPCService = NewOrderRPCService(t.orderRepository, t.userService)
}

func (t *OrderRPCTest) TestCreateOrder_OnUserServiceError_ReturnsError() {
	request := &rpc.NewOrderMessage{
		CustomerId: 100,
	}
	t.userService.err = status.New(codes.Aborted, "Aborted").Err()
	response, err := t.orderRPCService.CreateOrder(context.TODO(), request)

	t.Equal(t.orderRepository.transactionCount, 1, "Transaction not used")
	t.Nil(response)
	t.NotNil(err)
}

func (t *OrderRPCTest) TestCreateOrder_OnOrderRepositoryError_ReturnsError() {
	request := &rpc.NewOrderMessage{
		CustomerId: 100,
	}
	t.orderRepository.err = apperrs.ErrServerError
	response, err := t.orderRPCService.CreateOrder(context.TODO(), request)

	t.Equal(t.orderRepository.transactionCount, 1, "Transaction not used")
	t.Nil(response)
	t.NotNil(err)
}

func (t *OrderRPCTest) TestCreateOrder_OnSuccess_ReturnsOrder() {
	request := &rpc.NewOrderMessage{
		CustomerId: 100,
	}
	t.orderRepository.orderId = 100
	t.userService.user = &rpc.UserMessage{Id: 200}
	response, err := t.orderRPCService.CreateOrder(context.TODO(), request)

	t.Equal(t.orderRepository.transactionCount, 1, "Transaction not used")
	t.Nil(err)
	t.Equal(response.Customer.Id, int64(200))
	t.Equal(response.Id, int64(100))
}

func (t *OrderRPCTest) TestCreateOrder_OnRequest_HandlesSubRequestsConcurrently() {
	request := &rpc.NewOrderMessage{
		CustomerId: 100,
	}
	sleepDuration := time.Second * 3

	minSequentialRequestDuration := sleepDuration * 2

	t.orderRepository.sleep = sleepDuration
	t.orderRepository.orderId = 100
	t.userService.user = &rpc.UserMessage{Id: 200}
	t.userService.sleep = sleepDuration

	start := time.Now()
	response, err := t.orderRPCService.CreateOrder(context.TODO(), request)
	end := time.Now()

	elapsed := end.Sub(start)

	t.Less(elapsed, minSequentialRequestDuration)
	t.Equal(t.orderRepository.transactionCount, 1, "Transaction not used")
	t.Nil(err)
	t.Equal(response.Customer.Id, int64(200))
	t.Equal(response.Id, int64(100))
}

func TestOrderRPC(t *testing.T) {
	suite.Run(t, new(OrderRPCTest))
}
