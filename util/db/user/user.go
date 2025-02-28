package dbutil

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func IsNewUser(ctx context.Context, email, userName string, db *pgxpool.Pool) bool {
	if err := db.QueryRow(ctx, "SELECT email, userName FROM users WHERE email = $1 OR userName = $2", email, userName).Scan(&email, &userName); err != nil {
		return true
	}
	return false
}
