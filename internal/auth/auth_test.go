package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// First, we need to create some tokens for testing
	user, _ := uuid.NewRandom()
	tokensecret1 := "correctPassword123!"
	tokensecret2 := "anotherPassword456!"
	token1, _ := MakeJWT(user, tokensecret1, time.Hour) // Token valid for 1 hour
	token2, _ := MakeJWT(user, tokensecret2, 0)         // Token valid for 0 hours

	tests := []struct {
		name        string
		tokensecret string
		token       string
		wantErr     bool
		matchUser   uuid.UUID
	}{
		{
			name:        "Correct token",
			tokensecret: tokensecret1,
			token:       token1,
			wantErr:     false,
			matchUser:   user,
		},
		{
			name:        "Incorrect token",
			tokensecret: tokensecret2,
			token:       token1,
			wantErr:     true,
			matchUser:   user,
		},
		{
			name:        "Token doesn't match different hash",
			tokensecret: tokensecret1,
			token:       token2,
			wantErr:     true,
			matchUser:   user,
		},
		{
			name:        "Empty token",
			tokensecret: "",
			token:       token1,
			wantErr:     true,
			matchUser:   user,
		},
		{
			name:        "Invalid date",
			tokensecret: tokensecret2,
			token:       token2,
			wantErr:     true,
			matchUser:   user,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := ValidateJWT(tt.token, tt.tokensecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchUser {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.matchUser, match)
			}
		})
	}
}
