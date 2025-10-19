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

	logger := util.SharedLogger
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

// CheckEmail checks if an email exists in users table
func CheckEmail(ctx context.Context, email string, db *pgxpool.Pool) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    
    var exists bool
    err := db.QueryRow(ctx, query, email).Scan(&exists)
    if err != nil {
        return false, err
    }
    
    return exists, nil
}

// GetUserIDByEmail returns user ID for given email
func GetUserIDByEmail(ctx context.Context, email string, db *pgxpool.Pool) (string, error) {
    query := `SELECT id FROM users WHERE email = $1`
    
    var userID string
    err := db.QueryRow(ctx, query, email).Scan(&userID)
    if err != nil {
        return "", err
    }
    
    return userID, nil
}

// UpdatePassword updates user's password hash
func UpdatePassword(ctx context.Context, userID string, hashedPassword string, db *pgxpool.Pool) error {
    query := `UPDATE users SET password = $1, updatedAt = now() WHERE id = $2`
    
    _, err := db.Exec(ctx, query, hashedPassword, userID)
    return err
}
