package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/dino16m/chpter/order/models"
	"github.com/dino16m/chpter/order/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host        string
	Port        uint
	DatabaseURL string

	UserServiceURL string
}

func getConfig() Config {
	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
	if err != nil {
		panic("Failed to parse PORT")
	}
	return Config{
		Host:           os.Getenv("HOST"),
		Port:           uint(port),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		UserServiceURL: os.Getenv("USER_SERVICE_URL"),
	}
}

func initDB(databaseURL string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(databaseURL))
	if err != nil {
		panic("Failed to open database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Second)

	return db
}

func runServer(grpcServer *grpc.Server, host string, port uint) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	logrus.Infof("Starting gRPC server on %s:%d", host, port)
	return grpcServer.Serve(listener)
}

func registerGRPCServices(server *grpc.Server, services RPCServices) {
	rpc.RegisterOrderRPCServiceServer(server, services.OrderRPCService)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	config := getConfig()
	db := initDB(config.DatabaseURL)
	models.Connect(db)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	services, err := getServices(db, config)
	if err != nil {
		panic(err)
	}
	registerGRPCServices(grpcServer, services)

	err = runServer(grpcServer, config.Host, config.Port)

	logrus.Fatal(err.Error())
}
