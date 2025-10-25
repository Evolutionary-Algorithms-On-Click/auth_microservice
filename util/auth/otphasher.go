package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashOTP creates a SHA-256 hash of the OTP for secure storage
func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(hash[:])
}
