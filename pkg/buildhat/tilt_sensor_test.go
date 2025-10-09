package buildhat

import (
	"testing"
	"time"
)

// Multi-call integration test
func TestTiltSensor_AllMethods(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Create sensor
	sensor := brick.TiltSensor(PortA)
	if sensor == nil {
		t.Fatal("TiltSensor returned nil")
	}

	// Test GetTilt - queue response just before calling
	mockPort.SimulateSensorResponse("0", 0, "15 -10 5") // X Y Z for GetTilt
	time.Sleep(10 * time.Millisecond)                   // Let reader process it
	mockPort.ClearWriteHistory()                        // Clear any previous commands
	tilt, err := sensor.GetTilt()
	if err != nil {
		t.Fatalf("GetTilt failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; select 0\r" (mode 0 for tilt)
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) > 0 {
		expectedCmd := "port 0 ; select 0\r"
		if writeHistory[0] != expectedCmd {
			t.Errorf("Expected exact command '%s', got: %s", expectedCmd, writeHistory[0])
		}
	}

	expected := struct{ X, Y, Z int }{X: 15, Y: -10, Z: 5}
	if tilt != expected {
		t.Errorf("Expected tilt %+v, got %+v", expected, tilt)
	}

	// Test GetDirection - queue response just before calling
	mockPort.SimulateSensorResponse("0", 0, "46 -10 5") // X > 45 = TiltRight for GetDirection
	time.Sleep(10 * time.Millisecond)                   // Let reader process it
	mockPort.ClearWriteHistory()                        // Clear previous commands
	direction, err := sensor.GetDirection()
	if err != nil {
		t.Fatalf("GetDirection failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; select 0\r" (GetDirection calls GetTilt internally)
	writeHistory = mockPort.GetWriteHistory()
	if len(writeHistory) > 0 {
		expectedCmd := "port 0 ; select 0\r"
		if writeHistory[0] != expectedCmd {
			t.Errorf("Expected exact command '%s', got: %s", expectedCmd, writeHistory[0])
		}
	}

	if direction != TiltRight {
		t.Errorf("Expected direction TiltRight, got %v (%s)", direction, direction.String())
	}
}

// Single-call test for GetTilt
func TestTiltSensor_GetTilt(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 0, "15 -10 5") // X Y Z

	sensor := brick.TiltSensor(PortA)
	tilt, err := sensor.GetTilt()
	if err != nil {
		t.Fatalf("GetTilt failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; select 0\r" (mode 0 for tilt)
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected command to be sent")
	}
	expectedCmd := "port 0 ; select 0\r"
	if writeHistory[0] != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, writeHistory[0])
	}

	expected := struct{ X, Y, Z int }{X: 15, Y: -10, Z: 5}
	if tilt != expected {
		t.Errorf("Expected tilt %+v, got %+v", expected, tilt)
	}
}

// Test all tilt directions
func TestTiltSensor_GetDirection_AllDirections(t *testing.T) {
	tests := []struct {
		name     string
		tiltData string
		expected TiltDirection
	}{
		{"Right", "50 0 0", TiltRight},
		{"Left", "-50 0 0", TiltLeft},
		{"Forward", "0 50 0", TiltForward},
		{"Backward", "0 -50 0", TiltBackward},
		{"Level", "0 0 0", TiltLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brick := TestBrick(t)
			defer CleanupTestBrick(brick)

			mockPort := brick.GetMockPort()
			mockPort.SimulateSensorResponse("0", 0, tt.tiltData)
			time.Sleep(10 * time.Millisecond)

			sensor := brick.TiltSensor(PortA)
			direction, err := sensor.GetDirection()
			if err != nil {
				t.Fatalf("GetDirection failed: %v", err)
			}

			// Verify EXACT command: "port 0 ; select 0\r" (GetDirection calls GetTilt)
			writeHistory := mockPort.GetWriteHistory()
			if len(writeHistory) == 0 {
				t.Fatal("Expected command to be sent")
			}
			expectedCmd := "port 0 ; select 0\r"
			if writeHistory[0] != expectedCmd {
				t.Errorf("Expected exact command '%s', got: %s", expectedCmd, writeHistory[0])
			}

			if direction != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, direction)
			}
		})
	}
}
