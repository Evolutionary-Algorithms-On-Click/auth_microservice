package controller

import (
    "net/http"
    
    "evolve/modules/resetpassword"
    "evolve/util"
)

func ResetPasswordVerify(w http.ResponseWriter, r *http.Request) {
    logger := util.SharedLogger
    logger.InfoCtx(r, "Reset password verify API called")
    
    // Parse request body
    data, err := util.Body(r)
    if err != nil {
        util.JSONResponse(w, http.StatusBadRequest, err.Error(), nil)
        return
    }
    
    // Extract fields
    email, _ := data["email"].(string)
    otp, _ := data["otp"].(string)
    newPassword, _ := data["new_password"].(string)
    
    if email == "" || otp == "" || newPassword == "" {
        util.JSONResponse(w, http.StatusBadRequest, "Email, OTP, and new password are required", nil)
        return
    }
    
    // Call business logic
    err = resetpassword.VerifyAndResetPassword(r.Context(), email, otp, newPassword)
    if err != nil {
        util.JSONResponse(w, http.StatusBadRequest, err.Error(), nil)
        return
    }
    
    util.JSONResponse(w, http.StatusOK, "Password reset successful", nil)
}
