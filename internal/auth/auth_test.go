package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {

	userID := uuid.New()
	tokenSecret := "my-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	returnedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	if returnedUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, returnedUserID)
	}

	// Test with an expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT for expired token returned an error: %v", err)
	}

	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}

	// Test with an invalid token
	_, err = ValidateJWT("invalid-token", tokenSecret)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// Test with a token signed with a different secret
	otherSecret := "other-secret-key"
	_, err = ValidateJWT(token, otherSecret)
	if err == nil {
		t.Error("Expected error for token signed with different secret, got nil")
	}

}
func TestGetBearerToken(t *testing.T) {

	headers := http.Header{}
	headers.Set("Authorization", "Bearer some-token")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}

	if token != "some-token" {
		t.Errorf("Expected token 'some-token', got '%v'", token)
	}

	// Test with missing Authorization header
	headers = http.Header{}
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for missing Authorization header, got nil")
	}

	// Test with invalid Authorization header format
	headers.Set("Authorization", "InvalidFormat some-token")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for invalid Authorization header format, got nil")
	}

	// test with wrong format of Authorization header
	headers.Set("Authorization", "Bearersome-token")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for wrong format of Authorization header, got nil")
	}

	// Test with empty token in Authorization header
	headers.Set("Authorization", "Bearer ")
	_, err = GetBearerToken(headers)
	if err == nil {
		t.Error("Expected error for empty token in Authorization header, got nil")
	}

	// Test with extra spaces in Authorization header
	headers.Set("Authorization", "Bearer     some-token   ")
	token, err = GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error with extra spaces: %v", err)
	}
	if token != "some-token" {
		t.Errorf("Expected token 'some-token' with extra spaces, got '%v'", token)
	}

}
