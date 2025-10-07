package buildhat

import (
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// MockSerialPort for testing
type MockSerialPort struct {
	responses []string
	index     int
	commands  []string
}

func (m *MockSerialPort) Read(p []byte) (n int, err error) {
	if m.index >= len(m.responses) {
		time.Sleep(10 * time.Millisecond)
		return 0, nil
	}

	response := m.responses[m.index] + "\n"
	m.index++

	copy(p, response)
	return len(response), nil
}

func (m *MockSerialPort) Write(p []byte) (n int, err error) {
	m.commands = append(m.commands, strings.TrimSpace(string(p)))
	return len(p), nil
}

func TestBrickCreation(t *testing.T) {
	mockPort := &MockSerialPort{
		responses: []string{
			"Firmware version: 1636109636 2021-11-05T10:53:56+00:00",
			"12.5",
			"P0: connected to active ID 31", // SpikePrimeLargeMotor
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Reduce noise in tests
	}))

	brick := NewBrick(mockPort, mockPort, logger)
	defer brick.Close()

	// Give it time to initialize and process device connection
	time.Sleep(200 * time.Millisecond)

	// Test basic functionality
	if err := brick.SetLedMode(models.Green); err != nil {
		t.Errorf("Failed to set LED mode: %v", err)
	}

	// Manually set sensor type for testing (since the mock doesn't trigger the connection handler)
	brick.SetSensorType(models.PortA, models.SpikePrimeLargeMotor)

	// Test motor power (should work now that we have a motor connected)
	if err := brick.SetPowerLevel(models.PortA, 50); err != nil {
		t.Errorf("Failed to set motor power: %v", err)
	}

	if err := brick.FloatMotor(models.PortA); err != nil {
		t.Errorf("Failed to float motor: %v", err)
	}

	// Check that commands were sent
	if len(mockPort.commands) == 0 {
		t.Error("No commands were sent")
	}
}

func TestSensorTypeMethods(t *testing.T) {
	// Test IsMotor
	if !models.SpikePrimeLargeMotor.IsMotor() {
		t.Error("SpikePrimeLargeMotor should be identified as a motor")
	}

	if models.SpikePrimeColorSensor.IsMotor() {
		t.Error("SpikePrimeColorSensor should not be identified as a motor")
	}

	// Test IsActiveSensor
	if !models.SpikePrimeLargeMotor.IsActiveSensor() {
		t.Error("SpikePrimeLargeMotor should be identified as an active sensor")
	}

	if models.SystemMediumMotor.IsActiveSensor() {
		t.Error("SystemMediumMotor should not be identified as an active sensor")
	}
}

func TestSensorPortMethods(t *testing.T) {
	// Test Byte conversion
	if models.PortA.Byte() != 0 {
		t.Error("PortA should have byte value 0")
	}

	if models.PortB.Byte() != 1 {
		t.Error("PortB should have byte value 1")
	}

	// Test String conversion
	if models.PortA.String() != "Port A" {
		t.Error("PortA string representation should be 'Port A'")
	}
}

func TestLedModeMethods(t *testing.T) {
	// Test String conversion
	if models.Green.String() != "Green" {
		t.Error("Green LED mode string should be 'Green'")
	}

	if models.VoltageDependant.String() != "Voltage Dependent" {
		t.Error("VoltageDependant LED mode string should be 'Voltage Dependent'")
	}
}

func TestPositionWayMethods(t *testing.T) {
	// Test String conversion
	if models.Shortest.String() != "Shortest" {
		t.Error("Shortest position way string should be 'Shortest'")
	}

	if models.Clockwise.String() != "Clockwise" {
		t.Error("Clockwise position way string should be 'Clockwise'")
	}
}

func TestLedColorMethods(t *testing.T) {
	// Test String conversion
	if models.LedOff.String() != "Off" {
		t.Error("LedOff string should be 'Off'")
	}

	if models.LedRed.String() != "Red" {
		t.Error("LedRed string should be 'Red'")
	}

	if models.LedGreen.String() != "Green" {
		t.Error("LedGreen string should be 'Green'")
	}

	if models.LedBlue.String() != "Blue" {
		t.Error("LedBlue string should be 'Blue'")
	}

	if models.LedPaleGreen.String() != "Pale Green" {
		t.Error("LedPaleGreen string should be 'Pale Green'")
	}
}
