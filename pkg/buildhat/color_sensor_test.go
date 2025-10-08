package buildhat

import (
	"strings"
	"testing"
)

func TestColorSensor_GetColor(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Setup mock scanner

	// Simulate sensor responses (mode 5 returns RGBI in 0-1024 range)
	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("D", 5, "512 256 768 1024") // R G B I

	// Create color sensor
	sensor := brick.ColorSensor(PortD)

	// Test getting color
	color, err := sensor.GetColor()
	if err != nil {
		t.Fatalf("GetColor failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 1 {
		t.Fatalf("Expected at least 1 command, got %d", len(writeHistory))
	}

	// Check select mode command (mode 5 for RGBI)
	selectCmd := writeHistory[0]
	if !strings.Contains(selectCmd, "port 3") || !strings.Contains(selectCmd, "select 5") {
		t.Errorf("Select mode command incorrect: %s", selectCmd)
	}

	// Verify color values (converted from 0-1024 to 0-255)
	// 512/1024*255 = 127, 256/1024*255 = 63, 768/1024*255 = 191, 1024/1024*255 = 255
	expected := Color{R: 127, G: 63, B: 191, A: 255}
	if color != expected {
		t.Errorf("Expected color %+v, got %+v", expected, color)
	}
}

func TestColorSensor_GetReflectedLight(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Setup mock scanner

	// Simulate sensor responses (mode 1 for reflected light)
	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("D", 1, "75")

	// Create color sensor
	sensor := brick.ColorSensor(PortD)

	// Test getting reflected light
	light, err := sensor.GetReflectedLight()
	if err != nil {
		t.Fatalf("GetReflectedLight failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 1 {
		t.Fatalf("Expected at least 1 command, got %d", len(writeHistory))
	}

	// Check select mode command
	selectCmd := writeHistory[0]
	if !strings.Contains(selectCmd, "port 3") || !strings.Contains(selectCmd, "select 1") {
		t.Errorf("Select mode command incorrect: %s", selectCmd)
	}

	// Verify light value
	if light != 75 {
		t.Errorf("Expected reflected light 75, got %d", light)
	}
}

func TestColorSensor_GetAmbientLight(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Setup mock scanner

	// Simulate sensor responses (mode 2 for ambient light)
	mockPort := brick.GetMockPort()
	mockPort.SimulateSensorResponse("D", 2, "45")

	// Create color sensor
	sensor := brick.ColorSensor(PortD)

	// Test getting ambient light
	light, err := sensor.GetAmbientLight()
	if err != nil {
		t.Fatalf("GetAmbientLight failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 1 {
		t.Fatalf("Expected at least 1 command, got %d", len(writeHistory))
	}

	// Check select mode command
	selectCmd := writeHistory[0]
	if !strings.Contains(selectCmd, "port 3") || !strings.Contains(selectCmd, "select 2") {
		t.Errorf("Select mode command incorrect: %s", selectCmd)
	}

	// Verify light value
	if light != 45 {
		t.Errorf("Expected ambient light 45, got %d", light)
	}
}
