package buildhat

import (
	"strings"
	"testing"
)

func TestLight_On(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	light := brick.Light(PortA)

	err := light.On()
	if err != nil {
		t.Fatalf("On failed: %v", err)
	}

	// Verify command was sent
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	lastCmd := commands[len(commands)-1]
	if !strings.Contains(lastCmd, "port 0") {
		t.Errorf("Expected port 0 in command, got: %s", lastCmd)
	}
}

func TestLight_Off(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	light := brick.Light(PortB)

	err := light.Off()
	if err != nil {
		t.Fatalf("Off failed: %v", err)
	}

	// Verify coast command was sent
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	lastCmd := commands[len(commands)-1]
	if !strings.Contains(lastCmd, "port 1") || !strings.Contains(lastCmd, "coast") {
		t.Errorf("Expected 'port 1 ; coast' command, got: %s", lastCmd)
	}
}

func TestLight_SetBrightness(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	light := brick.Light(PortC)

	testCases := []struct {
		brightness int
		shouldFail bool
	}{
		{0, false},   // Should call Off()
		{50, false},  // Valid
		{100, false}, // Valid
		{-1, true},   // Invalid - too low
		{101, true},  // Invalid - too high
	}

	for _, tc := range testCases {
		err := light.SetBrightness(tc.brightness)
		if tc.shouldFail && err == nil {
			t.Errorf("SetBrightness(%d) should have failed but didn't", tc.brightness)
		}
		if !tc.shouldFail && err != nil {
			t.Errorf("SetBrightness(%d) failed: %v", tc.brightness, err)
		}
	}
}

func TestLight_SetBrightness_Values(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	light := brick.Light(PortD)

	// Test mid-range brightness
	err := light.SetBrightness(75)
	if err != nil {
		t.Fatalf("SetBrightness(75) failed: %v", err)
	}

	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	lastCmd := commands[len(commands)-1]

	// Should contain port 3 and set command with value ~0.75
	if !strings.Contains(lastCmd, "port 3") || !strings.Contains(lastCmd, "set") {
		t.Errorf("Expected brightness command, got: %s", lastCmd)
	}
}

func TestLight_GetBrightness(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	light := brick.Light(PortA)

	// GetBrightness should return default value when no data available
	brightness, err := light.GetBrightness()
	if err != nil {
		t.Fatalf("GetBrightness failed: %v", err)
	}

	// Should return default of 50 when no data
	if brightness != 50 {
		t.Errorf("Expected default brightness 50, got %d", brightness)
	}
}

func TestLight_AllPorts(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	ports := []struct {
		port     Port
		expected int
	}{
		{PortA, 0},
		{PortB, 1},
		{PortC, 2},
		{PortD, 3},
	}

	for _, tc := range ports {
		light := brick.Light(tc.port)
		if light.port != tc.expected {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, light.port)
		}
	}
}
