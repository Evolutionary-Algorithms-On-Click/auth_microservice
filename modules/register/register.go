package register

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"evolve/util/auth"
	dbutil "evolve/util/db/user"
	mailer "evolve/util/mail"
	"fmt"
)

type RegisterReq struct {
	Email    string `json:"email"`
	UserName string `json:"userName"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
}

func RegisterReqFromJSON(data map[string]any) (*RegisterReq, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var regReq *RegisterReq
	if err := json.Unmarshal(jsonData, &regReq); err != nil {
		return nil, err
	}

	return regReq, nil
}

func (r *RegisterReq) validate() error {
	if len(r.Email) == 0 {
		return fmt.Errorf("invalid email: %v", r.Email)
	}

	if len(r.UserName) == 0 {
		return fmt.Errorf("invalid username: %v", r.UserName)
	}

	if len(r.FullName) == 0 {
		return fmt.Errorf("invalid full name: %v", r.FullName)
	}

	if len(r.Password) == 0 {
		return fmt.Errorf("invalid password: %v", r.Password)
	}

	return nil
}

func (r *RegisterReq) Register(ctx context.Context) (string, error) {
	var logger = util.NewLogger()

	if err := r.validate(); err != nil {
		return "", err
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Register: failed to get pool connection: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	// Check if user already registered.
	if isNewUser := dbutil.IsNewUser(ctx, r.Email, r.UserName, db); !isNewUser {
		return "", fmt.Errorf("already registered")
	}

	// Delete user if already exists in registerOtp table.
	if _, err := db.Exec(ctx, "DELETE FROM registerOtp WHERE email = $1", r.Email); err != nil {
		logger.Error(fmt.Sprintf("Register: failed to delete from registerOtp: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	otp := auth.GenerateOTP()
	ct, err := db.Exec(ctx, "INSERT INTO registerOtp(email, otp) VALUES($1, $2)", r.Email, otp)
	if err != nil {
		logger.Error(fmt.Sprintf("Register: failed to insert into registerOtp: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	if ct.RowsAffected() == 0 {
		logger.Error("Register: failed to insert into registerOtp")
		return "", fmt.Errorf("something went wrong")
	}

	// logger.Info(fmt.Sprintf("OTP: %v", otp))
	if err := mailer.OTPVerifyEmail(r.Email, otp); err != nil {
		logger.Error(fmt.Sprintf("Register: failed to send OTP email: %v", err))
		return "", fmt.Errorf("something went wrong - check your email again")
	}

	// Generate Token.
	token, err := auth.Token(map[string]string{
		"email":    r.Email,
		"userName": r.UserName,
		"fullName": r.FullName,
		"password": r.Password,
		"purpose":  "register",
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Register: failed to generate token: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	return token, nil
}
