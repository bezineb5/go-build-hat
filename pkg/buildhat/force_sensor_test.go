package buildhat

import (
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

func TestForceSensor_AllPorts(t *testing.T) {
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
		sensor := brick.ForceSensor(tc.port)
		if sensor.port != Port(tc.expected) {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, sensor.port)
		}
	}
}
