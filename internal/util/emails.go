package util

import "net/mail"

// note that this also accepts local addresses like "user@localhost"
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
