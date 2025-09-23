package mailer

import (
	"evolve/util"
	"fmt"
	"net/smtp"
	"os"
)

var (
	from     = os.Getenv("MAILER_EMAIL")
	hostname = "smtp.gmail.com"
	port     = "587"
)

func email(to string, subject string, body string) error {
	
	logger := util.Log_var
	if err := smtp.SendMail(
		fmt.Sprintf("%s:%s", hostname, port),
		smtp.PlainAuth("", from, os.Getenv("MAILER_PASSWORD"), hostname),
		from,
		[]string{to},
		[]byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)),
	); err != nil {
		logger.Error("failed to send email")
		return err
	}

	logger.Info(fmt.Sprintf("email sent to %s", to))
	return nil
}

func OTPVerifyEmail(to string, otp string) error {
	return email(to, "[EvOC] OTP Verification", fmt.Sprintf("Your OTP is %s", otp))
}
