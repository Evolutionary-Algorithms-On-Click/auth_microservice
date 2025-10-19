package resetpassword

import (
    "context"
    "fmt"
    
    "golang.org/x/crypto/bcrypt"
    
    "evolve/db/connection"
    "evolve/util/db/resetpassword"
    dbutil "evolve/util/db/user"
)

// VerifyAndResetPassword verifies OTP and updates password
func VerifyAndResetPassword(ctx context.Context, email, otpCode, newPassword string) error {
    // Validate inputs
    if email == "" || otpCode == "" || newPassword == "" {
        return fmt.Errorf("email, OTP, and new password are required")
    }
    
    if len(newPassword) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
    }
    
    // Get database connection
    db, err := connection.PoolConn(ctx)
    if err != nil {
        return fmt.Errorf("database connection error: %w", err)
    }
    
    // Check if email exists
    exists, err := dbutil.CheckEmail(ctx, email, db)
    if err != nil {
        return fmt.Errorf("database error: %w", err)
    }
    if !exists {
        return fmt.Errorf("invalid email or OTP")
    }
    
    // Get user ID
    userID, err := dbutil.GetUserIDByEmail(ctx, email, db)
    if err != nil {
        return fmt.Errorf("failed to get user ID: %w", err)
    }
    
    // Verify OTP
    isValid, err := resetpassword.VerifyOTP(ctx, userID, otpCode, db)
    if err != nil {
        return fmt.Errorf("failed to verify OTP: %w", err)
    }
    if !isValid {
        return fmt.Errorf("invalid or expired OTP")
    }
    
    // Hash new password (same as registration)
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(newPassword),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Update password
    err = dbutil.UpdatePassword(ctx, userID, string(hashedPassword), db)
    if err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }
    
    // Mark OTP as used
    err = resetpassword.MarkOTPAsUsed(ctx, userID, otpCode, db)
    if err != nil {
        // Logs but doesnt't fail - password will have beeen already updated
        fmt.Printf("Warning: failed to mark OTP as used: %v\n", err)
    }
    
    return nil
}
