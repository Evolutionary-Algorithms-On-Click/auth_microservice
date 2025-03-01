package main

import (
	"context"
	"evolve/config"
	"evolve/controller"
	"evolve/db"
	"evolve/routes"
	"evolve/util"
	"fmt"
	"net/http"

	"aidanwoods.dev/go-paseto"
)

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

	// Register routes.
	http.HandleFunc(routes.TEST, controller.Test)
	http.HandleFunc(routes.REGISTER, controller.Register)
	http.HandleFunc(routes.VERIFY, controller.Verify)
	http.HandleFunc(routes.LOGIN, controller.Login)

	logger.Info(fmt.Sprintf("Test http server on http://localhost%v/api/test", config.PORT))
	if err := http.ListenAndServe(config.PORT, nil); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return
	}
}
