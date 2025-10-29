package controller

import (
	"crypto/rand"
	"encoding/base64"
	"evolve/modules"
	"evolve/util"
	"net/http"
	"time"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func Login(res http.ResponseWriter, req *http.Request) {

	logger := util.SharedLogger
	logger.InfoCtx(req, "Login API called.")

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

	csrfToken := generateCSRFToken()
	http.SetCookie(res, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})

	res.Header().Set("X-CSRF-Token", csrfToken)

	util.JSONResponse(res, http.StatusOK, "Success", user)
}
