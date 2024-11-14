package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
	//"github.com/stretchr/testify/assert"
	//"fmt"	
)

func TestJWTWithDebug(t *testing.T) {
	// Create a fixed test UUID for consistency
	userID := uuid.New()
	secret := "your-test-secret-key"
	
	// Generate token
	t.Log("Generating token...")
	token, err := MakeJWT(userID, secret, 24*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Log token details
	t.Logf("Generated token: %s", token)
	t.Logf("Token length: %d", len(token))
	t.Logf("Token bytes: %v", []byte(token))
	
	// Validate token ONCE
	t.Log("Validating token...")
	parsedUserID, err := ValidateJWT(token, secret)
	
	// Log validation results
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	} else {
		t.Logf("Validation successful. Parsed UUID: %s", parsedUserID)
		if parsedUserID != userID {
			t.Errorf("UUID mismatch. Expected: %s, Got: %s", userID, parsedUserID)
		}
	}
}

func TestInvalidTokenWithDebug(t *testing.T) {
	secret := "your-test-secret-key"
	
	invalidToken := "invalid.token.string"
	t.Logf("Testing invalid token: %s", invalidToken)
	
	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

