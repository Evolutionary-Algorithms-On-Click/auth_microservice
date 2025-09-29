package dbutil

import (
	"context"
	"evolve/util"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func IsNewUser(ctx context.Context, email, userName string, db *pgxpool.Pool) bool {
	if err := db.QueryRow(ctx, "SELECT email, userName FROM users WHERE email = $1 OR userName = $2", email, userName).Scan(&email, &userName); err != nil {
		return true
	}
	return false
}

func UserById(ctx context.Context, id string, db *pgxpool.Pool) (map[string]string, error) {
	
	logger := util.LogVar
	var role, email, userName, fullName string
	if err := db.QueryRow(ctx, "SELECT role, email, userName, fullName FROM users WHERE id = $1", id).Scan(&role, &email, &userName, &fullName); err != nil {
		logger.Error(fmt.Sprintf("Error getting user by id %v: %v", id, err), err)
		return nil, fmt.Errorf("user not found")
	}

	return map[string]string{
		"id":       id,
		"role":     role,
		"email":    email,
		"userName": userName,
		"fullName": fullName,
	}, nil
}
