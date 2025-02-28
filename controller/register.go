package controller

import (
	"context"
	"evolve/modules"
	"evolve/util"
	"net/http"
)

func Register(res http.ResponseWriter, req *http.Request) {
	var logger = util.NewLogger()
	logger.Info("Register API called.")

	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	regReq, err := modules.RegisterReqFromJSON(data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := regReq.Register(context.Background())
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// TODO: Properly set the token in cookie and test it.
	// Set the token in cookie.
	http.SetCookie(res, &http.Cookie{
		Name:  "token",
		Value: token,
	})

	util.JSONResponse(res, http.StatusOK, "Success", nil)
}
