package modules

import (
	"context"
	"encoding/json"
	"evolve/db/connection"
	"evolve/util"
	"evolve/util/auth"
	dbutil "evolve/util/db/user"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") ||
		strings.HasPrefix(s, "$2b$") ||
		strings.HasPrefix(s, "$2y$")
}

func (l *LoginReq) Login(ctx context.Context) (map[string]string, error) {

	logger := util.SharedLogger
	if err := l.validate(); err != nil {
		return nil, err
	}

	db, err := connection.PoolConn(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Login: failed to get pool connection: %v", err), err)
		return nil, fmt.Errorf("something went wrong")
	}

	// id is UUID(16 bytes) Google's UUID.
	// Query user and get password hash
	var id uuid.UUID
	var role string
	var storedPasswordHash string
	err = db.QueryRow(ctx, "SELECT id, role, password FROM users WHERE username = $1 OR email = $2", l.UserName, l.Email).Scan(&id, &role, &storedPasswordHash)

	if err != nil {
		logger.Error(fmt.Sprintf("Login: failed to query user: %v", err), err)
		return nil, fmt.Errorf("invalid username/email or password")
	}

	if isBcryptHash(storedPasswordHash) {
		err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(l.Password))
		if err != nil {
			logger.Info("Login: invalid password attempt")
			return nil, fmt.Errorf("invalid username/email or password")
		}
	} else {
		// Compare plain text for legacy users
		if storedPasswordHash != l.Password {
			logger.Info("Login: invalid password attempt (legacy)")
			return nil, fmt.Errorf("invalid username/email or password")
		}

		// Rehash and update DB for this user
		newHash, err := bcrypt.GenerateFromPassword([]byte(l.Password), bcrypt.DefaultCost)
		if err == nil {
			_, _ = db.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", string(newHash), id)
			logger.Info(fmt.Sprintf("Upgraded password hash for user %v", id))
		}
	}

	// Password verified, get user details
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
