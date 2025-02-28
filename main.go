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
)

func main() {
	var logger = util.NewLogger()

	// Initialize db with schema.
	if err := db.InitDb(context.Background()); err != nil {
		logger.Error("failed to init db")
		logger.Error(err.Error())
		return
	}

	// Register routes.
	http.HandleFunc(routes.TEST, controller.Test)

	if err := http.ListenAndServe(config.PORT, nil); err != nil {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		return
	}
	logger.Info(fmt.Sprintf("Test http server on http://localhost%v/api/test", config.PORT))
}
