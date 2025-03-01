package controller

import (
	"evolve/modules/register"
	"evolve/util"
	"evolve/util/auth"
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

	regReq, err := register.RegisterReqFromJSON(data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := regReq.Register(req.Context())
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

func Verify(res http.ResponseWriter, req *http.Request) {
	var logger = util.NewLogger()
	logger.Info("Verify API called.")

	token, err := req.Cookie("t")
	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "You got to try way better than that.", nil)
		return
	}

	// Validate the token.
	payLoad, err := auth.ValidateToken(token.Value)
	if err != nil {
		util.JSONResponse(res, http.StatusUnauthorized, "Session Expired.", nil)
		return
	}

	// Check if purpose is register.
	if payLoad["purpose"] != "register" {
		util.JSONResponse(res, http.StatusUnauthorized, "Good try.", nil)
		return
	}

	data, err := util.Body(req)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	verifyReq, err := register.VerifyReqFromJSON(data)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = verifyReq.Verify(req.Context(), payLoad)
	if err != nil {
		util.JSONResponse(res, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.JSONResponse(res, http.StatusOK, "Registration Successful.", nil)
}
