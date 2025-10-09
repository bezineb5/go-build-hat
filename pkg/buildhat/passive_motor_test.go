package buildhat

import (
	"fmt"
	"strings"
	"testing"
)

func TestPassiveMotor_Start(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortA)

	err := motor.Start(75)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; set 75\r"
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	expectedCmd := "port 0 ; set 75\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestPassiveMotor_Start_DefaultSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortB)

	// Speed 0 should use default of 50
	err := motor.Start(0)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Verify EXACT command: "port 1 ; set 50\r"
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	expectedCmd := "port 1 ; set 50\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestPassiveMotor_Stop(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortC)

	err := motor.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Verify EXACT command: "port 2 ; set 0\r"
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	expectedCmd := "port 2 ; set 0\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestPassiveMotor_SetSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortD)

	speeds := []int{25, 50, 75, 100}
	for _, speed := range speeds {
		err := motor.SetSpeed(speed)
		if err != nil {
			t.Fatalf("SetSpeed(%d) failed: %v", speed, err)
		}

		// Verify EXACT command: "port 3 ; set <speed>\r"
		mockPort := brick.GetMockPort()
		commands := mockPort.GetWriteHistory()
		expectedCmd := fmt.Sprintf("port 3 ; set %d\r", speed)
		lastCmd := commands[len(commands)-1]
		if lastCmd != expectedCmd {
			t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
		}
	}
}

func TestPassiveMotor_AllPorts(t *testing.T) {
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
		motor := brick.PassiveMotor(tc.port)
		if motor.port != tc.expected {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, motor.port)
		}
	}
}

func TestPassiveMotor_Start_InvalidSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortA)

	testCases := []int{-101, -150, 101, 150, 200}

	for _, speed := range testCases {
		err := motor.Start(speed)
		if err == nil {
			t.Errorf("Start(%d) should have failed but didn't", speed)
		}
		if err != nil && !strings.Contains(err.Error(), "must be between -100 and 100") {
			t.Errorf("Start(%d) error message should mention valid range, got: %v", speed, err)
		}
	}
}

func TestPassiveMotor_Start_ValidBounds(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortA)

	// Test boundary values
	validSpeeds := []int{-100, -50, 50, 100}

	for _, speed := range validSpeeds {
		err := motor.Start(speed)
		if err != nil {
			t.Errorf("Start(%d) should have succeeded but failed: %v", speed, err)
		}
	}
}

func TestPassiveMotor_SetSpeed_InvalidSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortB)

	testCases := []int{-101, -200, 101, 150}

	for _, speed := range testCases {
		err := motor.SetSpeed(speed)
		if err == nil {
			t.Errorf("SetSpeed(%d) should have failed but didn't", speed)
		}
		if err != nil && !strings.Contains(err.Error(), "must be between -100 and 100") {
			t.Errorf("SetSpeed(%d) error message should mention valid range, got: %v", speed, err)
		}
	}
}

func TestPassiveMotor_SetSpeed_ValidBounds(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.PassiveMotor(PortB)

	// Test boundary values and negative speeds
	validSpeeds := []int{-100, -75, -50, -25, 0, 25, 50, 75, 100}

	for _, speed := range validSpeeds {
		err := motor.SetSpeed(speed)
		if err != nil {
			t.Errorf("SetSpeed(%d) should have succeeded but failed: %v", speed, err)
		}
	}
}
