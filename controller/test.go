package controller

import (
	"evolve/util"
	"net/http"
)

func Test(res http.ResponseWriter, req *http.Request) {
	
	logger := util.Log_var
	switch req.Method {
	case "GET":
		logger.Info("GET /api/test called.")
		util.JSONResponse(res, http.StatusOK, "It works! 👍🏻", nil)
	default:
		util.JSONResponse(res, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}
