package util

import "net/mail"

func ValidMail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
