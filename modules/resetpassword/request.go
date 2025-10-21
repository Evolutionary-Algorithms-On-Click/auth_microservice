package resetpassword

import (
	"context"
	"fmt"

	"evolve/db/connection"
	"evolve/util/auth"
	"evolve/util/db/resetpassword"
	dbutil "evolve/util/db/user"
	mailer "evolve/util/mail"
)

// RequestPasswordReset generates and sends OTP for password reset
func RequestPasswordReset(ctx context.Context, email string) error {
	// Get database connection
	db, err := connection.PoolConn(ctx)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Check if email exists
	exists := !dbutil.IsNewUser(ctx, email, "", db)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// not revealing if email exists or not for security reasons
	if !exists {
		return nil
	}

	// Get user ID
	userID, err := dbutil.UserIDFromEmail(ctx, email, db)
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	// Generate OTP
	otpCode := auth.GenerateOTP()
	hashedOTP := auth.HashOTP(otpCode)
	
	// Save OTP to database
	err = resetpassword.SaveOTP(ctx, userID, hashedOTP, db)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send OTP email (uses your existing mailer function)
	err = mailer.OTPVerifyEmail(email, otpCode)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
