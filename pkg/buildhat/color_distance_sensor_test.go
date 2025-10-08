package buildhat

import (
	"testing"
	"time"
)

// Multi-call integration test
func TestColorDistanceSensor_AllMethods(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Create sensor
	sensor := brick.ColorDistanceSensor(PortA)
	if sensor == nil {
		t.Fatal("ColorDistanceSensor returned nil")
	}

	// Test GetColor - queue response just before calling
	mockPort.SimulateSensorResponse("0", 0, "128 64 192 255") // R G B A for GetColor
	time.Sleep(10 * time.Millisecond)                         // Let reader process it
	color, err := sensor.GetColor()
	if err != nil {
		t.Fatalf("GetColor failed: %v", err)
	}

	expected := Color{R: 128, G: 64, B: 192, A: 255}
	if color != expected {
		t.Errorf("Expected color %+v, got %+v", expected, color)
	}

	// Test GetDistance - queue response just before calling
	mockPort.SimulateSensorResponse("0", 1, "15") // Distance for GetDistance
	time.Sleep(10 * time.Millisecond)             // Let reader process it
	distance, err := sensor.GetDistance()
	if err != nil {
		t.Fatalf("GetDistance failed: %v", err)
	}

	if distance != 15 {
		t.Errorf("Expected distance 15, got %d", distance)
	}

	// Test GetReflectedLight - queue response just before calling
	mockPort.SimulateSensorResponse("0", 2, "85") // Reflected light
	time.Sleep(10 * time.Millisecond)             // Let reader process it
	light, err := sensor.GetReflectedLight()
	if err != nil {
		t.Fatalf("GetReflectedLight failed: %v", err)
	}

	if light != 85 {
		t.Errorf("Expected reflected light 85, got %d", light)
	}
}

// Single-call test for GetColor
func TestColorDistanceSensor_GetColor(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("0", 0, "128 64 192 255") // R G B A

	sensor := brick.ColorDistanceSensor(PortA)
	color, err := sensor.GetColor()
	if err != nil {
		t.Fatalf("GetColor failed: %v", err)
	}

	expected := Color{R: 128, G: 64, B: 192, A: 255}
	if color != expected {
		t.Errorf("Expected color %+v, got %+v", expected, color)
	}
}
