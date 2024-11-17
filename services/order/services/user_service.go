package services

import (
	"context"
	"time"

	"github.com/dino16m/chpter/order/apperrs"
	"github.com/dino16m/chpter/order/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	grpcConnection grpc.ClientConnInterface
	client         rpc.UserRPCServiceClient
}

func NewUserService(conn grpc.ClientConnInterface) *UserService {
	return &UserService{
		grpcConnection: conn,
		client:         rpc.NewUserRPCServiceClient(conn),
	}
}

func (s *UserService) GetUserById(id uint, ctx context.Context) (*rpc.UserMessage, error) {
	const MAX_TRIES = 5
	var customer *rpc.UserMessage
	var err error
	for tries := 0; tries < MAX_TRIES; tries++ {
		customer, err = s.client.GetUser(ctx, &rpc.IdMessage{Id: int64(id)})

		statusCode, _ := status.FromError(err)

		if statusCode.Code() == codes.NotFound {
			err = apperrs.NewNotFoundError("Customer not found", err)
			break
		}
		if statusCode.Code() == codes.Unavailable {
			time.Sleep(100 * time.Nanosecond * time.Duration(tries))
			continue
		}
		break
	}
	return customer, err
}
