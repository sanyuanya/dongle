package data

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	if err := CheckPasswordHash(password, hashedPassword); !err {
		t.Fatal(err)
	}

	wrongPassword := "wrong-password"
	if err := CheckPasswordHash(wrongPassword, hashedPassword); !err {
		t.Fatal("Expected error but got nil")
	}
	t.Log("All checks passed")
}
