package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestHashPassword(t *testing.T){
	authService := NewAuthService("test-secret")

	password := "password123"

	hash , err := authService.HashPassword(password)

		if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == password {
		t.Error("password should not be stored as plain text")
	}

	if hash == "" {
		t.Error("password hash should not be empty")
	}
}

func TestCheckPassword(t *testing.T) {
	authService := NewAuthService("test-secret")

	password := "password123"

	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !authService.CheckPassword(password, hash) {
		t.Error("correct password should be accepted")
	}

	if authService.CheckPassword("wrong-password", hash) {
		t.Error("wrong password should be rejected")
	}
}




func TestGenerateAccessToken(t *testing.T) {
	authService := NewAuthService("test-secret")

	token, err := authService.GenerateAcessTokens(
		"user-123",
		"ritik@example.com",
		"org-123",
	)

	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if token == "" {
		t.Error("access token should not be empty")
	}
}

func TestGenerateAccessTokenClaims(t *testing.T) {

	authService := NewAuthService("test-secret")

		tokenString, err := authService.GenerateAcessTokens(
		"user-123",
		"ritik@example.com",
		"org-123",
	)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		},
	)

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

		claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatal("failed to extract claims")
	}
if claims.UserID != "user-123" {
		t.Errorf("expected user ID user-123, got %s", claims.UserID)
	}

	if claims.Email != "ritik@example.com" {
		t.Errorf("expected email ritik@example.com, got %s", claims.Email)
	}

	if claims.OrganizationID != "org-123" {
		t.Errorf("expected organization ID org-123, got %s", claims.OrganizationID)
	}

	if !token.Valid {
		t.Error("token should be valid")
	}
}

func TestGenerateAccessTokenExpiry(t *testing.T) {
	authService := NewAuthService("test-secret")

	tokenString, err := authService.GenerateAcessTokens(
		"user-123",
		"ritik@example.com",
		"org-123",
	)

	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		},
	)

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatal("failed to extract claims")
	}

	if claims.ExpiresAt == nil {
		t.Fatal("token expiry should be set")
	}

	if !claims.ExpiresAt.Time.After(time.Now()) {
		t.Error("token should expire in the future")
	}
}

func TestAccessTokenWrongSecret(t *testing.T) {
	authService := NewAuthService("test-secret")

	tokenString, err := authService.GenerateAcessTokens(
		"user-123",
		"ritik@example.com",
		"org-123",
	)

	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("wrong-secret"), nil
		},
	)

	if err == nil {
		t.Fatal("token should fail with wrong secret")
	}

	if token != nil && token.Valid {
		t.Error("token should not be valid with wrong secret")
	}
}
func TestAccessTokenAlgorithm(t *testing.T) {
	authService := NewAuthService("test-secret")

	tokenString, err := authService.GenerateAcessTokens(
		"user-123",
		"ritik@example.com",
		"org-123",
	)

	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte("test-secret"), nil
		},
	)

	if err != nil {
		t.Fatalf("token should be valid with HS256: %v", err)
	}

	if !token.Valid {
		t.Error("token should be valid")
	}

	if token.Method != jwt.SigningMethodHS256 {
		t.Error("token should use HS256")
	}
}