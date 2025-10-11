package buildhat

import (
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

	// Verify EXACT command: "port 0 ; select 0\r"
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	expectedCmd := "port 0 ; select 0\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
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
		port     Port
		expected int
	}{
		{PortA, 0},
		{PortB, 1},
		{PortC, 2},
		{PortD, 3},
	}

	for _, tc := range ports {
		sensor := brick.ButtonSensor(tc.port)
		if sensor.port != Port(tc.expected) {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, sensor.port)
		}
	}
}
