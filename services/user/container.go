package main

import (
	"github.com/dino16m/chpter/user/repositories"
	"github.com/dino16m/chpter/user/rpcservices"
	"gorm.io/gorm"
)

type RPCServices struct {
	UserRPCService rpcservices.UserRPCService
}

func getServices(db *gorm.DB) RPCServices {
	userRepo := repositories.NewUserRepository(db)
	userRPCService := rpcservices.NewUserRPCService(userRepo)

	return RPCServices{
		UserRPCService: userRPCService,
	}
}
