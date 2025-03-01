package auth

import (
	"fmt"
	"time"
)

func GenerateOTP() string {
	return fmt.Sprint(time.Now().Nanosecond())[:6]
}
