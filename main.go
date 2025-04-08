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
	"os"
	"runtime"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

var (
	HTTP_PORT string
	GRPC_PORT string
)

func serveHTTP(logger *util.Logger) {
	// Register routes.
	http.HandleFunc(routes.TEST, controller.Test)
	http.HandleFunc(routes.REGISTER, controller.Register)
	http.HandleFunc(routes.VERIFY, controller.Verify)
	http.HandleFunc(routes.LOGIN, controller.Login)
	http.HandleFunc(routes.TEAMS, controller.Teams)

	logger.Info(fmt.Sprintf("Test http server on http://localhost%v/api/test", HTTP_PORT))
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")}, // Allowing frontend to access the server.
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)
	if err := http.ListenAndServe("0.0.0.0" + HTTP_PORT, corsHandler); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return
	}
}

func serveGRPC(logger *util.Logger) {
	lis, err := net.Listen("tcp", "0.0.0.0" + GRPC_PORT)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to listen TCP in GRPC PORT%v : %v", GRPC_PORT, err))
		return
	}

	s := grpc.NewServer()
	pb.RegisterAuthenticateServer(s, &grpcserver.GRPCServer{})
	logger.Info(fmt.Sprintf("Test grpc server on http://localhost%v", GRPC_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("failed to serve: %v", err))
		return
	}
}

func main() {

	HTTP_PORT = fmt.Sprintf(":%v", os.Getenv("HTTP_PORT"))
	GRPC_PORT = fmt.Sprintf(":%v", os.Getenv("GRPC_PORT"))

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
