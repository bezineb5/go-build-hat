package buildhat

import (
	"testing"
)

// Simplified sensor tests that avoid race conditions
// These test one method call at a time to avoid async timing issues

func TestColorDistanceSensor_GetColor_Simple(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 0, "128 64 192 255") // R G B A

	sensor := brick.ColorDistanceSensor("A")
	color, err := sensor.GetColor()
	if err != nil {
		t.Fatalf("GetColor failed: %v", err)
	}

	expected := struct{ R, G, B, A uint8 }{R: 128, G: 64, B: 192, A: 255}
	if color != expected {
		t.Errorf("Expected color %+v, got %+v", expected, color)
	}
}

func TestMotionSensor_GetDistance_Simple(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 0, "100")

	sensor := brick.MotionSensor("A")
	distance, err := sensor.GetDistance()
	if err != nil {
		t.Fatalf("GetDistance failed: %v", err)
	}

	if distance != 100 {
		t.Errorf("Expected distance 100, got %d", distance)
	}
}

func TestTiltSensor_GetTilt_Simple(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 0, "15 -10 5") // X Y Z

	sensor := brick.TiltSensor("A")
	tilt, err := sensor.GetTilt()
	if err != nil {
		t.Fatalf("GetTilt failed: %v", err)
	}

	expected := struct{ X, Y, Z int }{X: 15, Y: -10, Z: 5}
	if tilt != expected {
		t.Errorf("Expected tilt %+v, got %+v", expected, tilt)
	}
}
