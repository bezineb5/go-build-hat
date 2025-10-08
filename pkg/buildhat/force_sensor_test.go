package buildhat

import (
	"strings"
	"testing"
)

func TestForceSensor_GetForce(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("A", 0, "75")

	sensor := brick.ForceSensor(PortA)

	force, err := sensor.GetForce()
	if err != nil {
		t.Fatalf("GetForce failed: %v", err)
	}

	if force != 75 {
		t.Errorf("Expected force 75, got %d", force)
	}

	// Verify command
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	lastCmd := commands[len(commands)-1]
	if !strings.Contains(lastCmd, "port 0") || !strings.Contains(lastCmd, "select 0") {
		t.Errorf("Expected select mode command, got: %s", lastCmd)
	}
}

func TestForceSensor_AllPorts(t *testing.T) {
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
		sensor := brick.ForceSensor(tc.port)
		if sensor.port != tc.expected {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, sensor.port)
		}
	}
}
