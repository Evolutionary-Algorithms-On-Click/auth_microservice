package controller

import (
	"net/http"

	"evolve/modules/resetpassword"
	"evolve/util"
)

func ResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	logger := util.SharedLogger
	logger.InfoCtx(r, "Reset password request API called")

	// Parse request body
	data, err := util.Body(r)
	if err != nil {
		util.JSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Extract email
	email, ok := data["email"].(string)
	if !ok || email == "" {
		util.JSONResponse(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	// Call business logic
	err = resetpassword.RequestPasswordReset(r.Context(), email)
	if err != nil {
		logger.ErrorCtx(r, "Reset password request failed", err)
		util.JSONResponse(w, http.StatusInternalServerError, "Failed to process request", nil)
		return
	}

	util.JSONResponse(w, http.StatusOK, "If the email exists, an OTP has been sent", nil)
}
