package modules

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"evolve/util/auth"
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

func (l *LoginReq) Login(ctx context.Context) (string, error) {
	var logger = util.NewLogger()

	if err := l.validate(); err != nil {
		return "", err
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Register: failed to get pool connection: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	// id is UUID(16 bytes) Google's UUID.
	var id uuid.UUID
	var role string

	err = db.QueryRow(ctx, "SELECT id, role FROM users WHERE username = $1 OR email = $2 AND password = $3", l.UserName, l.Email, l.Password).Scan(&id, &role)
	if err != nil {
		logger.Error(fmt.Sprintf("Register: failed to query user: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	token, err := auth.Token(map[string]string{
		"id":      id.String(),
		"role":    role,
		"purpose": "login",
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Register: failed to generate token: %v", err))
		return "", fmt.Errorf("something went wrong")
	}

	return token, nil
}
