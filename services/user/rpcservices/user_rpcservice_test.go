package rpcservices

import (
	"context"
	"testing"

	"github.com/dino16m/chpter/user/apperrs"
	"github.com/dino16m/chpter/user/models"
	"github.com/dino16m/chpter/user/rpc"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserRepository struct {
	userId uint
	user   models.User
	err    error
}

func (m *MockUserRepository) Create(user *models.User) error {
	user.ID = m.userId
	return m.err
}
func (m *MockUserRepository) FindById(id uint) (models.User, error) {

	return m.user, m.err
}

type UserRPCServiceTest struct {
	suite.Suite

	userRPCService *UserRPCService
	userRepository *MockUserRepository
}

func (t *UserRPCServiceTest) SetupTest() {
	t.userRepository = &MockUserRepository{}
	t.userRPCService = &UserRPCService{userRepo: t.userRepository}
}

func (t *UserRPCServiceTest) TestCreateUser_OnSuccess_ReturnsUserMessage() {
	request := &rpc.NewUserMessage{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	t.userRepository.userId = 100
	response, err := t.userRPCService.CreateUser(context.Background(), request)

	t.Nil(err)
	t.Equal(response.Id, int64(100))
}

func (t *UserRPCServiceTest) TestCreateUser_OnError_ReturnsServerErrorCode() {
	request := &rpc.NewUserMessage{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	t.userRepository.err = status.Error(codes.Internal, "error creating user")

	response, err := t.userRPCService.CreateUser(context.Background(), request)

	t.Nil(response)
	t.Equal(codes.Internal, status.Code(err))
	t.NotNil(err)
}

func (t *UserRPCServiceTest) TestGetUser_OnSuccess_ReturnsUserMessage() {
	request := &rpc.IdMessage{
		Id: 100,
	}
	t.userRepository.user = models.User{FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	t.userRepository.user.ID = uint(request.Id)

	response, err := t.userRPCService.GetUser(context.Background(), request)

	t.Nil(err)
	t.Equal(response.Id, request.Id)
	t.Equal(response.FirstName, "John")
	t.Equal(response.LastName, "Doe")
	t.Equal(response.Email, "john@example.com")
}

func (t *UserRPCServiceTest) TestGetUser_OnNotFound_ReturnsNotFoundCode() {
	request := &rpc.IdMessage{
		Id: 100,
	}
	t.userRepository.err = apperrs.NewNotFoundError("User not found")

	response, err := t.userRPCService.GetUser(context.Background(), request)

	t.Nil(response)
	t.Equal(codes.NotFound, status.Code(err))
	t.NotNil(err)
}

func TestUserRPC(t *testing.T) {
	suite.Run(t, new(UserRPCServiceTest))
}
