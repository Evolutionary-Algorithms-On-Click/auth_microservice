package modules

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"evolve/util/auth"
	dbutil "evolve/util/db/user"
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
		logger.Error("Register: failed to get pool connection")
		return "", err
	}

	// Check if any user same email/username is unique.
	if isNewUser := dbutil.IsNewUser(ctx, r.Email, r.UserName, db); !isNewUser {
		return "", fmt.Errorf("user already exists")
	}

	// Delete user if already exists in registerOtp table.
	if _, err := db.Exec(ctx, "DELETE FROM registerOtp WHERE email = $1", r.Email); err != nil {
		logger.Error("Register: failed to delete from registerOtp")
		return "", err
	}

	if _, err := db.Exec(ctx, "INSERT INTO registerOtp(email, otp, purpose) VALUES($1, $2, $3)", r.Email, auth.GenerateOTP(), "register"); err != nil {
		logger.Error("Register: failed to insert into registerOtp")
		return "", err
	}

	// TODO: Send OTP to user's email.

	// Generate Token.
	token, err := auth.Token(map[string]any{
		"email":    r.Email,
		"userName": r.UserName,
		"fullName": r.FullName,
		"password": r.Password,
		"purpose":  "register",
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
