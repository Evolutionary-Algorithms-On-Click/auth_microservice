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

func serveHTTP() {
	logger := util.SharedLogger
	// Register routes.
	http.HandleFunc(routes.TEST, controller.Test)
	http.HandleFunc(routes.REGISTER, controller.Register)
	http.HandleFunc(routes.VERIFY, controller.Verify)
	http.HandleFunc(routes.LOGIN, controller.Login)
	http.HandleFunc(routes.PASSWORD_RESET, controller.ResetPasswordRequest)
	http.HandleFunc(routes.PASSWORD_VERIFY, controller.ResetPasswordVerify)
	http.HandleFunc(routes.CREATETEAM, controller.CreateTeam)
	http.HandleFunc(routes.GETTEAMS, controller.GetTeams)
	http.HandleFunc(routes.GETMEMBERS, controller.GetTeamMembers)
	http.HandleFunc(routes.ADDMEMBERS, controller.AddTeamMembers)
	http.HandleFunc(routes.DELETEMEMBERS, controller.DeleteTeamMembers)

	logger.Info(fmt.Sprintf("Test http server on http://localhost%v/api/test", HTTP_PORT))
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")}, // Allowing frontend to access the server.
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"X-CSRF-Token"},
	}).Handler(http.DefaultServeMux)

	handler := util.SharedLogger.LogMiddleware(corsHandler)

	finalHandler := util.CSRFMiddleware(handler)

	if err := http.ListenAndServe(HTTP_PORT, finalHandler); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err), err)
		return
	}
}

func serveGRPC() {
	logger := util.SharedLogger
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

	logger, err := util.InitLogger(os.Getenv("ENV"))
	if err != nil {
		fmt.Println("failed to init logger:", err)
		return
	}
	util.SharedLogger = logger

	// Initialize db with schema.
	if err := db.InitDb(context.Background()); err != nil {
		logger.Error("failed to init db", err)
		return
	}

	// Initialize key.
	key := paseto.NewV4AsymmetricSecretKey()
	config.PrivateKey, config.PublicKey = key, key.Public()

	go serveHTTP()
	go serveGRPC()

	runtime.Goexit()
}
