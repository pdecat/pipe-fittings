package utils

import (
	"errors"
	"net"
	"syscall"
	"testing"
)

// TestIsPortBindable tests the IsPortBindable function - assumes that the 9783 port is not in use
//
// Updated test to use random port 9783 rather than very commonly used 8080 that often cause
// this test to fail because there will be something running in 8080
func TestIsPortBindable(t *testing.T) {
	// Test case 1: Port is bindable
	err := IsPortBindable("localhost", 9783)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "9783"))
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	defer ln.Close()
	// Test case 2: Port is already in use
	err = IsPortBindable("localhost", 9783)
	if err == nil {
		t.Error("Expected an error, but got nil")
	} else {
		expectedErr := syscall.EADDRINUSE
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected error: %v, but got: %v", expectedErr, err)
		}
	}
}
