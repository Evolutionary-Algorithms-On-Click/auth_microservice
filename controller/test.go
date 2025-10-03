package controller

import (
	"evolve/util"
	"net/http"
)

func Test(res http.ResponseWriter, req *http.Request) {

	logger := util.SharedLogger
	switch req.Method {
	case "GET":
		logger.InfoCtx(req, "GET /api/test called.")
		util.JSONResponse(res, http.StatusOK, "It works! 👍🏻", nil)
	default:
		util.JSONResponse(res, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}
