package domain

import "testing"

func TestNewEmailNormalizesValidEmail(t *testing.T) {
	email, err := NewEmail("  ADA@EXAMPLE.COM  ")
	if err != nil {
		t.Fatalf("expected valid email, got %v", err)
	}
	if email.String() != "ada@example.com" {
		t.Fatalf("expected normalized email, got %q", email.String())
	}
}

func TestNewEmailRejectsDisplayNameAddress(t *testing.T) {
	if _, err := NewEmail("Ada <ada@example.com>"); err == nil {
		t.Fatal("expected display-name email to be rejected")
	}
}
