package register

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	dbutil "evolve/util/db/user"
	"fmt"
)

type VerifyReq struct {
	OTP string `json:"otp"`
}

func VerifyReqFromJSON(data map[string]any) (*VerifyReq, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var verifyReq *VerifyReq
	if err := json.Unmarshal(jsonData, &verifyReq); err != nil {
		return nil, err
	}

	return verifyReq, nil
}

func (r *VerifyReq) validate() error {
	if len(r.OTP) != 6 {
		return fmt.Errorf("invalid otp: %v", r.OTP)
	}

	return nil
}

func (r *VerifyReq) Verify(ctx context.Context, user map[string]string) error {
	
	logger := util.Log_var
	if err := r.validate(); err != nil {
		return err
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("RegisterVerify: failed to get pool connection: %v", err))
		return fmt.Errorf("something went wrong")
	}

	// Check if user already registered.
	if isNewUser := dbutil.IsNewUser(ctx, user["email"], "", db); !isNewUser {
		return fmt.Errorf("already registered")
	}

	ct, err := db.Exec(ctx, "DELETE FROM registerOtp WHERE email = $1 AND otp = $2", user["email"], r.OTP)
	if err != nil {
		logger.Error(fmt.Sprintf("RegisterVerify: failed to delete from registerOtp: %v", err))
		return fmt.Errorf("something went wrong")
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("invalid otp")
	}

	// Insert user into user table.
	if _, err := db.Exec(ctx, "INSERT INTO users (email, userName, fullName, password) VALUES($1, $2, $3, $4)", user["email"], user["userName"], user["fullName"], user["password"]); err != nil {
		logger.Error(fmt.Sprintf("RegisterVerify: failed to insert into users: %v", err))
		return fmt.Errorf("something went wrong")
	}

	return nil

}
