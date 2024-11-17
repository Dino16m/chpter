package services

import (
	"context"
	"testing"

	"github.com/dino16m/chpter/order/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockUserRPCServiceClient struct {
	mock.Mock

	responseCB func(int64) (*rpc.UserMessage, error)
	tries      int64
}

func (m *MockUserRPCServiceClient) CreateUser(ctx context.Context, in *rpc.NewUserMessage, opts ...grpc.CallOption) (*rpc.UserMessage, error) {
	m.tries++
	if m.responseCB != nil {
		return m.responseCB(m.tries)
	}
	return nil, nil
}
func (m *MockUserRPCServiceClient) GetUser(ctx context.Context, in *rpc.IdMessage, opts ...grpc.CallOption) (*rpc.UserMessage, error) {
	m.tries++
	if m.responseCB != nil {
		return m.responseCB(m.tries)
	}
	return nil, nil
}

type UserServiceTest struct {
	suite.Suite

	userClient  *MockUserRPCServiceClient
	userService *UserService
}

func (t *UserServiceTest) SetupTest() {
	t.userClient = &MockUserRPCServiceClient{}
	t.userService = &UserService{
		client: t.userClient,
	}
}

func (t *UserServiceTest) TestGetUser_OnSuccess_ReturnsUser() {
	t.userClient.responseCB = func(tries int64) (*rpc.UserMessage, error) {
		return &rpc.UserMessage{Id: 100, FirstName: "John", LastName: "Doe"}, nil
	}

	user, err := t.userService.GetUserById(1, context.Background())

	t.Nil(err)
	t.Equal(user.FirstName, "John")
	t.Equal(user.LastName, "Doe")
	t.Equal(int64(1), t.userClient.tries)
}

func (t *UserServiceTest) TestGetUser_OnUpstreamUnavailable_Retries() {
	var expectedTries int64 = 3
	t.userClient.responseCB = func(tries int64) (*rpc.UserMessage, error) {
		if tries == expectedTries {
			return &rpc.UserMessage{Id: 100, FirstName: "John", LastName: "Doe"}, nil
		}
		return nil, status.Error(codes.Unavailable, "upstream service unavailable")
	}

	user, err := t.userService.GetUserById(1, context.Background())

	t.Nil(err)
	t.Equal(user.FirstName, "John")
	t.Equal(user.LastName, "Doe")
	t.Equal(expectedTries, t.userClient.tries)
}

func (t *UserServiceTest) TestGetUser_OnUpstreamError_ReturnsError() {
	t.userClient.responseCB = func(tries int64) (*rpc.UserMessage, error) {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	_, err := t.userService.GetUserById(1, context.Background())

	t.NotNil(err)
	t.Equal(int64(1), t.userClient.tries)
}

func TestUserRPC(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}
