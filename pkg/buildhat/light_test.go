package buildhat

import (
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

	// Verify EXACT command: "port 0 ; on ; set 1.00\r"
	// On() calls SetBrightness(100) which converts to 1.00
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	expectedCmd := "port 0 ; on ; set 1.00\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
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

	// Verify EXACT command: "port 1 ; coast\r"
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	expectedCmd := "port 1 ; coast\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
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

	// Verify EXACT command: "port 3 ; on ; set 0.75\r"
	// 75/100 = 0.75, formatted as %.2f = 0.75
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	expectedCmd := "port 3 ; on ; set 0.75\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
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
		if light.port != Port(tc.expected) {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, light.port)
		}
	}
}
