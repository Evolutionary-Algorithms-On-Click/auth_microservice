package controller

import (
	"evolve/util"
	"log"
	"net/http"
)

func Test(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		log.Println("[INFO]: Test API called from", req.RemoteAddr)
		util.JSONResponse(res, http.StatusOK, "It works! 👍🏻", nil)
	default:
		util.JSONResponse(res, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}
