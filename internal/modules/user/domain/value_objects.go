package domain

import (
	"net/mail"
	"strings"
)

type Email string

func NewEmail(value string) (Email, error) {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return "", ErrInvalidEmail
	}

	parsed, err := mail.ParseAddress(trimmed)
	if err != nil || parsed.Address != trimmed || parsed.Name != "" {
		return "", ErrInvalidEmail
	}

	return Email(trimmed), nil
}

func (e Email) String() string {
	return string(e)
}
