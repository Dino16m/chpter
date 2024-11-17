package main

import (
	"github.com/dino16m/chpter/order/repositories"
	"github.com/dino16m/chpter/order/rpcservices"
	"github.com/dino16m/chpter/order/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type RPCServices struct {
	OrderRPCService rpcservices.OrderRPCService
}

var retryPolicy = `{
            "methodConfig": [{
                "retryPolicy": {
                    "MaxAttempts": 4,
                    "InitialBackoff": ".01s",
                    "MaxBackoff": ".02s",
                    "BackoffMultiplier": 1.1,
                    "RetryableStatusCodes": [ "UNAVAILABLE" ]
                }
            }]
        }`

func getServices(db *gorm.DB, config Config) (RPCServices, error) {
	orderRepo := repositories.NewOrderRepository(db)
	client, err := grpc.NewClient(config.UserServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		return RPCServices{}, err
	}
	orderRPCService := rpcservices.NewOrderRPCService(orderRepo, services.NewUserService(client))
	return RPCServices{
		OrderRPCService: orderRPCService,
	}, nil
}
