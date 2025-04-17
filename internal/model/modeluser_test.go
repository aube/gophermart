package model

import (
	"testing"
)

func TestUserValidate(t *testing.T) {
	t.Run("valid login and password", func(t *testing.T) {
		u := &User{
			Login:    "user@example.com",
			Password: "password123",
		}

		err := u.Validate()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		u := &User{
			Login:    "user@example.com",
			Password: "", // empty password
		}

		err := u.Validate()

		if err == nil {
			t.Error("expected error, got none")
		}
	})
}
