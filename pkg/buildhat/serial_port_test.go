package buildhat

import (
	"strings"
	"testing"
)

// TestNewSerialPort_ErrorHandling tests that NewSerialPort returns appropriate errors
// for non-existent ports without trying to access real hardware
func TestNewSerialPort_ErrorHandling(t *testing.T) {
	// Test with a clearly non-existent port
	_, err := NewSerialPort("/dev/nonexistent_serial_port_12345")
	if err == nil {
		t.Error("Expected error when opening nonexistent port")
	}

	// Verify error message format
	errStr := err.Error()
	if !strings.Contains(errStr, "failed to open serial port") {
		t.Errorf("Expected error message to contain 'failed to open serial port', got: %s", errStr)
	}
}

// Note: We intentionally do NOT test GetAvailablePorts or DetectBuildHatPort
// because they would scan real hardware ports on the system.
// These functions are designed to interact with real hardware and should only
// be tested in integration tests with actual BuildHat hardware present.
