package buildhat

import (
	"testing"
)

func TestDistanceSensor_GetDistance(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("A", 0, "150")

	sensor := brick.DistanceSensor(PortA)

	distance, err := sensor.GetDistance()
	if err != nil {
		t.Fatalf("GetDistance failed: %v", err)
	}

	if distance != 150 {
		t.Errorf("Expected distance 150mm, got %d", distance)
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

func TestDistanceSensor_AllPorts(t *testing.T) {
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
		sensor := brick.DistanceSensor(tc.port)
		if sensor.port != Port(tc.expected) {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, sensor.port)
		}
	}
}
