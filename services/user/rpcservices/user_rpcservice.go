package rpcservices

import (
	"context"

	"github.com/dino16m/chpter/user/models"
	"github.com/dino16m/chpter/user/repositories"
	"github.com/dino16m/chpter/user/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRepository interface {
	Create(user *models.User) error
	FindById(id uint) (models.User, error)
}

type UserRPCService struct {
	userRepo UserRepository
	rpc.UnimplementedUserRPCServiceServer
}

func NewUserRPCService(userRepo *repositories.UserRepository) UserRPCService {
	return UserRPCService{userRepo: userRepo}
}

func toUserMessage(user models.User) *rpc.UserMessage {
	return &rpc.UserMessage{
		Id:        int64(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: user.CreatedAt.Unix(),
			Nanos:   int32(user.CreatedAt.Nanosecond()),
		},
	}
}
func (s UserRPCService) CreateUser(ctx context.Context, in *rpc.NewUserMessage) (*rpc.UserMessage, error) {
	user := models.User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
	}
	err := s.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}
	return toUserMessage(user), nil
}

func (s UserRPCService) GetUser(ctx context.Context, in *rpc.IdMessage) (*rpc.UserMessage, error) {
	user, err := s.userRepo.FindById(uint(in.Id))
	if err != nil {
		return nil, err
	}
	return toUserMessage(user), nil
}
