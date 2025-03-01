package main

import (
	"context"
	"evolve/config"
	"evolve/controller"
	grpcserver "evolve/controller/grpc"
	"evolve/db"
	pb "evolve/proto"
	"evolve/routes"
	"evolve/util"
	"fmt"
	"net"
	"net/http"

	"runtime"

	"aidanwoods.dev/go-paseto"
	"google.golang.org/grpc"
)

func serveHTTP(logger *util.Logger) {
	// Register routes.
	http.HandleFunc(routes.TEST, controller.Test)
	http.HandleFunc(routes.REGISTER, controller.Register)
	http.HandleFunc(routes.VERIFY, controller.Verify)
	http.HandleFunc(routes.LOGIN, controller.Login)

	logger.Info(fmt.Sprintf("Test http server on http://localhost%v/api/test", config.HTTP_PORT))
	if err := http.ListenAndServe(config.HTTP_PORT, nil); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return
	}
}

func serveGRPC(logger *util.Logger) {
	lis, err := net.Listen("tcp", config.GRPC_PORT)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to listen TCP in GRPC PORT%v : %v", config.GRPC_PORT, err))
		return
	}

	s := grpc.NewServer()
	pb.RegisterAuthenticateServer(s, &grpcserver.GRPCServer{})
	logger.Info(fmt.Sprintf("Test grpc server on http://localhost%v", config.GRPC_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("failed to serve: %v", err))
		return
	}
}

func main() {
	var logger = util.NewLogger()

	// Initialize db with schema.
	if err := db.InitDb(context.Background()); err != nil {
		logger.Error("failed to init db")
		logger.Error(err.Error())
		return
	}

	// Initialize key.
	key := paseto.NewV4AsymmetricSecretKey()
	config.PrivateKey, config.PublicKey = key, key.Public()

	go serveHTTP(logger)
	go serveGRPC(logger)

	runtime.Goexit()
}
