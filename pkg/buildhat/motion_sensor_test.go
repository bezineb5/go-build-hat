package buildhat

import (
	"testing"
)

func TestMotionSensor_GetDistance(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 0, "100")

	sensor := brick.MotionSensor(PortA)
	distance, err := sensor.GetDistance()
	if err != nil {
		t.Fatalf("GetDistance failed: %v", err)
	}

	if distance != 100 {
		t.Errorf("Expected distance 100, got %d", distance)
	}
}

func TestMotionSensor_GetMovementCount(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 1, "42")

	sensor := brick.MotionSensor(PortA)
	count, err := sensor.GetMovementCount()
	if err != nil {
		t.Fatalf("GetMovementCount failed: %v", err)
	}

	if count != 42 {
		t.Errorf("Expected movement count 42, got %d", count)
	}
}

func TestMotionSensor_AllPorts(t *testing.T) {
	ports := []Port{PortA, PortB, PortC, PortD}

	for _, port := range ports {
		t.Run("Port"+port.String(), func(t *testing.T) {
			brick := TestBrick(t)
			defer CleanupTestBrick(brick)

			mockPort := brick.GetMockPort()
			portNum := string(rune('0' + port.Int()))
			mockPort.SimulateSensorResponse(portNum, 0, "50")

			sensor := brick.MotionSensor(port)
			distance, err := sensor.GetDistance()
			if err != nil {
				t.Fatalf("GetDistance failed for port %s: %v", port, err)
			}

			if distance != 50 {
				t.Errorf("Expected distance 50, got %d", distance)
			}
		})
	}
}
