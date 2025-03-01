package controller

import (
	"evolve/modules"
	"evolve/util"
	"net/http"
)

func Login(res http.ResponseWriter, req *http.Request) {
	var logger = util.NewLogger()
	logger.Info("Login API called.")

	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	loginReq, err := modules.LoginReqFromJSON(data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := loginReq.Login(req.Context())
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Set the token in cookie.
	http.SetCookie(res, &http.Cookie{
		Name:     "t",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	util.JSONResponse(res, http.StatusOK, "Success", nil)
}
