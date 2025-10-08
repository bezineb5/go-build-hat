package buildhat

import (
	"strings"
	"testing"
)

func TestButtonSensor_IsPressed_True(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("A", 0, "1")

	sensor := brick.ButtonSensor(PortA)

	pressed, err := sensor.IsPressed()
	if err != nil {
		t.Fatalf("IsPressed failed: %v", err)
	}

	if !pressed {
		t.Error("Expected button to be pressed (true), got false")
	}

	// Verify command was sent
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	lastCmd := commands[len(commands)-1]
	if !strings.Contains(lastCmd, "port 0") || !strings.Contains(lastCmd, "select 0") {
		t.Errorf("Expected select mode command, got: %s", lastCmd)
	}
}

func TestButtonSensor_IsPressed_False(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("B", 0, "0")

	sensor := brick.ButtonSensor(PortB)

	pressed, err := sensor.IsPressed()
	if err != nil {
		t.Fatalf("IsPressed failed: %v", err)
	}

	if pressed {
		t.Error("Expected button to be unpressed (false), got true")
	}
}

func TestButtonSensor_AllPorts(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	ports := []struct {
		port     BuildHatPort
		expected int
	}{
		{PortA, 0},
		{PortB, 1},
		{PortC, 2},
		{PortD, 3},
	}

	for _, tc := range ports {
		sensor := brick.ButtonSensor(tc.port)
		if sensor.port != tc.expected {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, sensor.port)
		}
	}
}
