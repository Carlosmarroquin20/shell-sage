package service

import "testing"

// TestGetUserName_Success tests the successful case
func TestGetUserName_Success(t *testing.T) {
	// Create service instance
	service := &UserService{}

	// Call method with valid ID
	name, err := service.GetUserName(1)

	// Check that no error occurred
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Validate returned name
	expected := "John Doe"
	if name != expected {
		t.Errorf("expected %s, got %s", expected, name)
	}
}

// TestGetUserName_InvalidID tests invalid input case
func TestGetUserName_InvalidID(t *testing.T) {
	service := &UserService{}

	// Call method with invalid ID
	_, err := service.GetUserName(0)

	// Ensure error is returned
	if err == nil {
		t.Errorf("expected error for invalid ID, got nil")
	}
}

// TestGetUserName_NotFound tests user not found case
func TestGetUserName_NotFound(t *testing.T) {
	service := &UserService{}

	// Call method with non-existing user ID
	_, err := service.GetUserName(99)

	// Ensure error is returned
	if err == nil {
		t.Errorf("expected error for user not found, got nil")
	}
}