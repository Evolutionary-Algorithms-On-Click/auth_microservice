package modules

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"evolve/util/auth"
	dbutil "evolve/util/db/user"
	"fmt"

	"github.com/google/uuid"
)

type LoginReq struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginReqFromJSON(data map[string]any) (*LoginReq, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var loginReq *LoginReq
	if err := json.Unmarshal(jsonData, &loginReq); err != nil {
		return nil, err
	}

	return loginReq, nil
}

func (l *LoginReq) validate() error {
	if len(l.Password) == 0 {
		return fmt.Errorf("invalid password: %v", l.Password)
	}

	// Either username or email should be provided.
	if len(l.UserName) == 0 && len(l.Email) == 0 {
		return fmt.Errorf("invalid username/email: %v %v", l.UserName, l.Email)
	}

	return nil
}

func (l *LoginReq) Login(ctx context.Context) (map[string]string, error) {
	
	logger := util.LogVar
	if err := l.validate(); err != nil {
		return nil, err
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Login: failed to get pool connection: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	// id is UUID(16 bytes) Google's UUID.
	var id uuid.UUID
	var role string

	err = db.QueryRow(ctx, "SELECT id, role FROM users WHERE (username = $1 OR email = $2) AND password = $3", l.UserName, l.Email, l.Password).Scan(&id, &role)
	if err != nil {
		logger.Error(fmt.Sprintf("Login: failed to query user: %v", err), err)
		return nil, fmt.Errorf("invalid username/email or password")
	}

	user, err := dbutil.UserById(ctx, id.String(), db)
	if err != nil {
		logger.Error(fmt.Sprintf("Login: failed to get user by id: %v", err), err)
		return nil, fmt.Errorf("user not found")
	}

	// Generate token.
	token, err := auth.Token(map[string]string{
		"id":      id.String(),
		"role":    role,
		"purpose": "login",
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Login: failed to generate token: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	user["token"] = token

	return user, nil
}
