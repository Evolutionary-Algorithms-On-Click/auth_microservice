package main

import (
	"aidanwoods.dev/go-paseto"
	"context"
	"evolve/config"
	"evolve/controller"
	grpcserver "evolve/controller/grpc"
	"evolve/db"
	pb "evolve/proto"
	"evolve/routes"
	"evolve/util"
	"fmt"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"runtime"
)

var (
	HTTP_PORT string
	GRPC_PORT string
)

func serveHTTP(logger *util.LoggerService) {
	// Register routes.
	http.HandleFunc(routes.TEST, controller.Test)
	http.HandleFunc(routes.REGISTER, controller.Register)
	http.HandleFunc(routes.VERIFY, controller.Verify)
	http.HandleFunc(routes.LOGIN, controller.Login)

	
	logger.Info(fmt.Sprintf("Test http server on http://localhost%v/api/test", HTTP_PORT))
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")}, // Allowing frontend to access the server.
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	handler := util.Log.LogMiddleware(corsHandler)

	if err := http.ListenAndServe(HTTP_PORT, handler); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err), err)
		return
	}
}

func serveGRPC(logger *util.LoggerService) {
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to listen TCP in GRPC PORT%v : %v", GRPC_PORT, err), err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterAuthenticateServer(s, &grpcserver.GRPCServer{})
	logger.Info(fmt.Sprintf("Test grpc server on http://localhost%v", GRPC_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Error(fmt.Sprintf("failed to serve: %v", err), err)
		return
	}
}

func main() {

	HTTP_PORT = fmt.Sprintf(":%v", os.Getenv("HTTP_PORT"))
	GRPC_PORT = fmt.Sprintf(":%v", os.Getenv("GRPC_PORT"))
	
	logger, err := util.InitLogger(os.Getenv("ENV")) // "DEVELOPMENT" or "PRODUCTION"
	if err != nil {
		fmt.Println("failed to init logger:", err)
		return
	}
	util.Log = logger

	

	// Initialize db with schema.
	if err := db.InitDb(context.Background()); err != nil {
		logger.Error("failed to init db", err)
		logger.Error(err.Error(), err)
		return
	}

	// Initialize key.
	key := paseto.NewV4AsymmetricSecretKey()
	config.PrivateKey, config.PublicKey = key, key.Public()

	go serveHTTP(util.Log)
	go serveGRPC(util.Log)

	runtime.Goexit()
}
