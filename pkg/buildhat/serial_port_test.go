package buildhat

import (
	"log/slog"
	"os"
	"testing"
)

func TestGetAvailablePorts(t *testing.T) {
	ports, err := GetAvailablePorts()
	if err != nil {
		t.Logf("GetAvailablePorts failed (expected on systems without serial ports): %v", err)
		return
	}

	t.Logf("Available ports: %v", ports)
}

func TestDetectBuildHatPort(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	port, err := DetectBuildHatPort(logger)
	if err != nil {
		t.Logf("DetectBuildHatPort failed (expected on systems without BuildHat): %v", err)
		return
	}

	t.Logf("Detected BuildHat port: %s", port)
}

func TestRealSerialPortCreation(t *testing.T) {

	// Try to create a serial port (this will likely fail on most systems)
	_, err := NewSerialPort("/dev/ttyAMA0")
	if err != nil {
		t.Logf("NewRealSerialPort failed (expected on systems without /dev/ttyAMA0): %v", err)
		return
	}

	t.Log("RealSerialPort created successfully")
}
