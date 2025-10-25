package controller

import (
	"net/http"

	"evolve/modules/resetpassword"
	"evolve/util"
)

// ResetPasswordRequest handles password reset OTP request
func ResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	logger := util.SharedLogger
	logger.InfoCtx(r, "Reset password request API called")

	data, err := util.Body(r)
	if err != nil {
		logger.ErrorCtx(r, "failed to parse request body", err)
		util.JSONResponse(w, http.StatusBadRequest, "Something went wrong. Please try again.", nil)
		return
	}

	email, ok := data["email"].(string)
	if !ok || email == "" {
		util.JSONResponse(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	err = resetpassword.RequestPasswordReset(r.Context(), email)
	if err != nil {
		logger.ErrorCtx(r, "Reset password request failed", err)
		util.JSONResponse(w, http.StatusInternalServerError, "Failed to process request", nil)
		return
	}

	util.JSONResponse(w, http.StatusOK, "If the email exists, an OTP has been sent", nil)
}

// ResetPasswordVerify handles OTP verification and password reset
func ResetPasswordVerify(w http.ResponseWriter, r *http.Request) {
	logger := util.SharedLogger
	logger.InfoCtx(r, "Reset password verify API called")

	data, err := util.Body(r)
	if err != nil {
		logger.ErrorCtx(r, "failed to parse request body", err)
		util.JSONResponse(w, http.StatusBadRequest, "Something went wrong. Please try again.", nil)
		return
	}

	email, _ := data["email"].(string)
	otp, _ := data["otp"].(string)
	newPassword, _ := data["new_password"].(string)

	if email == "" || otp == "" || newPassword == "" {
		logger.ErrorCtx(r, "missing required fields", nil)
		util.JSONResponse(w, http.StatusBadRequest, "Email, OTP, and new password are required", nil)
		return
	}

	err = resetpassword.VerifyAndResetPassword(r.Context(), email, otp, newPassword)
	if err != nil {
		logger.ErrorCtx(r, "Reset password verify failed", err)
		util.JSONResponse(w, http.StatusBadRequest, "Failed to verify OTP. Please try again.", nil)
		return
	}

	util.JSONResponse(w, http.StatusOK, "Password reset successful", nil)
}
