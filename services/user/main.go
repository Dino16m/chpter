package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/dino16m/chpter/user/models"
	"github.com/dino16m/chpter/user/rpc"
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
}

func getConfig() Config {
	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
	if err != nil {
		panic("Failed to parse PORT")
	}
	return Config{
		Host:        os.Getenv("HOST"),
		Port:        uint(port),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func initDB(databaseURL string) *gorm.DB {
	const MAX_TRIES = 5
	sleepDuration := time.Second * 5
	var err error
	for i := 0; i < MAX_TRIES; i++ {
		db, err := gorm.Open(mysql.Open(databaseURL))

		if err != nil {
			time.Sleep(sleepDuration * time.Duration(i))
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			time.Sleep(sleepDuration * time.Duration(i))
			continue
		}
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Second)
		return db
	}

	panic(err)

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
	rpc.RegisterUserRPCServiceServer(server, services.UserRPCService)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	config := getConfig()
	db := initDB(config.DatabaseURL)
	models.Connect(db)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	services := getServices(db)
	registerGRPCServices(grpcServer, services)

	err := runServer(grpcServer, config.Host, config.Port)

	logrus.Fatal(err.Error())
}
