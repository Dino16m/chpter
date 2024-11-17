package e2e_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/dino16m/chpter/e2e/rpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConcurrencyTest struct {
	suite.Suite
	customerId          int64
	userServiceAddress  string
	orderServiceAddress string
}

type Result struct {
	order *rpc.OrderMessage
	error error
}

type Report struct {
	results []Result
}

func (r *Report) AddResult(order *rpc.OrderMessage, err error) {
	r.results = append(r.results, Result{order: order, error: err})
}

func (r *Report) SuccessRate() float64 {
	successCount := 0
	for _, result := range r.results {
		if result.error == nil {
			successCount++
		}
	}
	return (float64(successCount) / float64(len(r.results))) * 100
}

func (t *ConcurrencyTest) SetupTest() {
	t.userServiceAddress = os.Getenv("USER_SERVICE_ADDR")
	t.orderServiceAddress = os.Getenv("ORDER_SERVICE_ADDR")

	userClient := t.getUserRPCService()
	user, err := userClient.CreateUser(context.Background(), &rpc.NewUserMessage{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	})
	if err != nil {
		panic(err)
	}
	t.customerId = user.Id
}

func (t *ConcurrencyTest) getUserRPCService() rpc.UserRPCServiceClient {

	client, err := grpc.NewClient(t.userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return rpc.NewUserRPCServiceClient(client)
}

func (t *ConcurrencyTest) getOrderRPCService() rpc.OrderRPCServiceClient {
	client, err := grpc.NewClient(t.orderServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	rpcServiceClient := rpc.NewOrderRPCServiceClient(client)

	return rpcServiceClient
}

func (t *ConcurrencyTest) createOrder(index int, orderClient rpc.OrderRPCServiceClient, customerId int64) (*rpc.OrderMessage, error) {
	return orderClient.CreateOrder(context.Background(), &rpc.NewOrderMessage{
		CustomerId: customerId,
		Items: []*rpc.NewLineItemMessage{
			{
				ProductId: int64((index + 1) * 2),
				Quantity:  2,
				UnitPrice: &rpc.Decimal{Value: "10"},
			},
			{
				ProductId: int64((index + 1) * 3),
				Quantity:  3,
				UnitPrice: &rpc.Decimal{Value: "12"},
			},
		},
	})
}

func (t *ConcurrencyTest) runConcurrentOrderCreationWithSharedConnection(concurrency int) *Report {
	orderClient := t.getOrderRPCService()
	wg := &sync.WaitGroup{}

	results := make(chan Result, concurrency)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()

			order, err := t.createOrder(i, orderClient, t.customerId)

			results <- Result{order: order, error: err}
		}(i)
	}
	wg.Wait()
	close(results)
	return buildReport(results)
}
func buildReport(results chan Result) *Report {
	report := &Report{}
	for result := range results {
		report.AddResult(result.order, result.error)
	}
	return report
}

func (t *ConcurrencyTest) TestConcurrency_OnSharedConcurrentRequests_HasHighThroughput() {

	report := t.runConcurrentOrderCreationWithSharedConnection(10000)
	fmt.Println("Concureency is ", len(report.results))
	fmt.Println("Success rate is", report.SuccessRate())

	t.GreaterOrEqual(report.SuccessRate(), float64(95))
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ConcurrencyTest))
}
