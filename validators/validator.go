package validators

import (
	"errors"
	"net"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) (bool, error) {

	if len(email) < 3 && len(email) > 254 {
		return false, errors.New("invalid email")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, err
	}

	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false, errors.New("invalid email, no MX records")
	}
	return true, nil
}
