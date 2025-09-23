package controller

import (
	"evolve/modules"
	"evolve/util"
	"net/http"
	"time"
)

func Login(res http.ResponseWriter, req *http.Request) {
	
	logger := util.Log_var
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

	user, err := loginReq.Login(req.Context())
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Set the token in cookie.
	http.SetCookie(res, &http.Cookie{
		Name:     "t",
		Value:    user["token"],
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(48 * time.Hour),
	})

	delete(user, "token")

	util.JSONResponse(res, http.StatusOK, "Success", user)
}
