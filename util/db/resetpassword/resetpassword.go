package resetpassword

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveOTP stores an OTP code in the database with 15-minute expiration
func SaveOTP(ctx context.Context, userID string, otpCode string, db *pgxpool.Pool) error {
	expiresAt := time.Now().Add(15 * time.Minute)

	query := `
        INSERT INTO password_reset_otps (user_id, otp_code, expires_at)
        VALUES ($1, $2, $3)
    `

	_, err := db.Exec(ctx, query, userID, otpCode, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	return nil
}

// VerifyOTP checks if an OTP is valid, not expired, and not used
func VerifyOTP(ctx context.Context, userID string, otpCode string, db *pgxpool.Pool) (bool, error) {
	query := `
        SELECT id, expires_at, is_used
        FROM password_reset_otps
        WHERE user_id = $1 AND otp_code = $2
        ORDER BY created_at DESC
        LIMIT 1
    `

	var otpID string
	var expiresAt time.Time
	var isUsed bool

	err := db.QueryRow(ctx, query, userID, otpCode).Scan(&otpID, &expiresAt, &isUsed)

	if err == sql.ErrNoRows {
		return false, nil // OTP not found
	}
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}

	// Check if expired
	if time.Now().After(expiresAt) {
		return false, nil
	}

	// Check if already used
	if isUsed {
		return false, nil
	}

	return true, nil
}

// MarkOTPAsUsed marks an OTP as used to prevent reuse
func MarkOTPAsUsed(ctx context.Context, userID string, otpCode string, db *pgxpool.Pool) error {
	query := `
        UPDATE password_reset_otps
        SET is_used = TRUE
        WHERE user_id = $1 AND otp_code = $2
    `

	_, err := db.Exec(ctx, query, userID, otpCode)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	return nil
}
