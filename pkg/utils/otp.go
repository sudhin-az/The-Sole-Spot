package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/smtp"
	"net/textproto"
	"os"
)

func GenerateOTP() string {
	otp := make([]byte, 6)
	rand.Read(otp)
	for i := 0; i < len(otp); i++ {
		otp[i] = otp[i]%10 + '0'
	}
	return string(otp)
}

func SendOTPEmail(email, otp string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration is not set correctly")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	from := smtpUser
	to := []string{email}
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is %s", otp)

	fmt.Println("UserOTP", otp)

	header := textproto.MIMEHeader{}
	header.Set("From", from)
	header.Set("To", email)
	header.Set("Subject", subject)
	header.Set("MIME-Version", "1.0")
	header.Set("Content-Type", "text/plain; charset=UTF-8")

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v[0])
	}
	message += "\r\n" + body

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// func GenerateReferralCode() string {
// 	return fmt.Sprintf("%s%d", utils.RandomString(5), time.Now().UnixNano()/1e6)
// }

func ParseJSON(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}
