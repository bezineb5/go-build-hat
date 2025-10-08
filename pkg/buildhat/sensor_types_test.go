package buildhat

import (
	"strings"
	"testing"
	"time"
)

func TestColorDistanceSensor(t *testing.T) {
	t.Skip("Skipping multi-call test due to async timing issues - see sensor_types_simple_test.go")
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Simulate device detection
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("list\r\n")
	mockPort.QueueReadData("P0: connected to active ID 25\r\n") // ColorDistanceSensor
	mockPort.QueueReadData("type 25\r\n")

	// Scan for devices
	err := brick.ScanDevices()
	if err != nil {
		t.Fatalf("ScanDevices failed: %v", err)
	}

	// Create sensor
	sensor := brick.ColorDistanceSensor("A")
	if sensor == nil {
		t.Fatal("ColorDistanceSensor returned nil")
	}

	// Queue all responses upfront
	mockPort.SimulateSensorResponse("0", 0, "128 64 192 255") // R G B A for GetColor
	mockPort.SimulateSensorResponse("0", 1, "15")             // Distance for GetDistance
	mockPort.SimulateSensorResponse("0", 2, "85")             // Reflected light

	// Test GetColor
	color, err := sensor.GetColor()
	if err != nil {
		t.Fatalf("GetColor failed: %v", err)
	}

	expected := struct{ R, G, B, A uint8 }{R: 128, G: 64, B: 192, A: 255}
	if color != expected {
		t.Errorf("Expected color %+v, got %+v", expected, color)
	}

	// Small delay to let async operations complete
	time.Sleep(50 * time.Millisecond)

	// Test GetDistance
	distance, err := sensor.GetDistance()
	if err != nil {
		t.Fatalf("GetDistance failed: %v", err)
	}

	if distance != 15 {
		t.Errorf("Expected distance 15, got %d", distance)
	}

	// Small delay to let async operations complete
	time.Sleep(50 * time.Millisecond)

	// Test GetReflectedLight
	light, err := sensor.GetReflectedLight()
	if err != nil {
		t.Fatalf("GetReflectedLight failed: %v", err)
	}

	if light != 85 {
		t.Errorf("Expected reflected light 85, got %d", light)
	}
}

func TestLight(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Simulate device detection
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("list\r\n")
	mockPort.QueueReadData("P0: connected to passive ID 8\r\n") // Light
	mockPort.QueueReadData("type 8\r\n")

	// Scan for devices
	err := brick.ScanDevices()
	if err != nil {
		t.Fatalf("ScanDevices failed: %v", err)
	}

	// Create light
	light := brick.Light("A")
	if light == nil {
		t.Fatal("Light returned nil")
	}

	// Test GetBrightness (returns default value when no data)
	brightness, err := light.GetBrightness()
	if err != nil {
		t.Fatalf("GetBrightness failed: %v", err)
	}

	// Should return default value (50)
	if brightness != 50 {
		t.Errorf("Expected brightness 50, got %d", brightness)
	}

	// Test SetBrightness
	err = light.SetBrightness(50)
	if err != nil {
		t.Fatalf("SetBrightness failed: %v", err)
	}

	// Verify set command was sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected set brightness command to be sent")
	}

	lastCmd := writeHistory[len(writeHistory)-1]
	if !containsAny(lastCmd, []string{"port 0", "set 50"}) {
		t.Errorf("Expected set brightness command, got: %s", lastCmd)
	}
}

func TestMatrix(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Simulate device detection
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("list\r\n")
	mockPort.QueueReadData("P0: connected to active ID 40\r\n") // Matrix
	mockPort.QueueReadData("type 40\r\n")

	// Scan for devices
	err := brick.ScanDevices()
	if err != nil {
		t.Fatalf("ScanDevices failed: %v", err)
	}

	// Create matrix
	matrix := brick.Matrix("A")
	if matrix == nil {
		t.Fatal("Matrix returned nil")
	}

	// Test SetPixel (color, brightness in 0-10 range)
	err = matrix.SetPixel(1, 2, 9, 10) // Red at full brightness
	if err != nil {
		t.Fatalf("SetPixel failed: %v", err)
	}

	// Test Clear
	err = matrix.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Test SetAll
	err = matrix.SetAll(3, 5) // Blue at medium brightness
	if err != nil {
		t.Fatalf("SetAll failed: %v", err)
	}

	// Test SetRow
	err = matrix.SetRow(0, 9, 10) // Red row at full brightness
	if err != nil {
		t.Fatalf("SetRow failed: %v", err)
	}

	// Test SetColumn
	err = matrix.SetColumn(1, 6, 8) // Green column at high brightness
	if err != nil {
		t.Fatalf("SetColumn failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 5 {
		t.Fatalf("Expected at least 5 commands, got %d", len(writeHistory))
	}
}

func TestPassiveMotor(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Simulate device detection
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("list\r\n")
	mockPort.QueueReadData("P0: connected to passive ID 1\r\n") // PassiveMotor
	mockPort.QueueReadData("type 1\r\n")

	// Scan for devices
	err := brick.ScanDevices()
	if err != nil {
		t.Fatalf("ScanDevices failed: %v", err)
	}

	// Create motor
	motor := brick.PassiveMotor("A")
	if motor == nil {
		t.Fatal("PassiveMotor returned nil")
	}

	// Test Start
	err = motor.Start(50)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Test Stop
	err = motor.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Test SetSpeed
	err = motor.SetSpeed(75)
	if err != nil {
		t.Fatalf("SetSpeed failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 3 {
		t.Fatalf("Expected at least 3 commands, got %d", len(writeHistory))
	}
}

func TestTiltSensor(t *testing.T) {
	t.Skip("Skipping multi-call test due to async timing issues - see sensor_types_simple_test.go")
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Simulate device detection
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("list\r\n")
	mockPort.QueueReadData("P0: connected to active ID 22\r\n") // TiltSensor
	mockPort.QueueReadData("type 22\r\n")

	// Scan for devices
	err := brick.ScanDevices()
	if err != nil {
		t.Fatalf("ScanDevices failed: %v", err)
	}

	// Create sensor
	sensor := brick.TiltSensor("A")
	if sensor == nil {
		t.Fatal("TiltSensor returned nil")
	}

	// Queue all responses upfront
	mockPort.SimulateSensorResponse("0", 0, "15 -10 5") // X Y Z for GetTilt
	mockPort.SimulateSensorResponse("0", 0, "46 -10 5") // X > 45 = "right" for GetDirection

	// Test GetTilt
	tilt, err := sensor.GetTilt()
	if err != nil {
		t.Fatalf("GetTilt failed: %v", err)
	}

	expected := struct{ X, Y, Z int }{X: 15, Y: -10, Z: 5}
	if tilt != expected {
		t.Errorf("Expected tilt %+v, got %+v", expected, tilt)
	}

	// Small delay to let async operations complete
	time.Sleep(50 * time.Millisecond)

	// Test GetDirection
	direction, err := sensor.GetDirection()
	if err != nil {
		t.Fatalf("GetDirection failed: %v", err)
	}

	if direction != "right" {
		t.Errorf("Expected direction 'right', got '%s'", direction)
	}
}

// Helper function to check if a string contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
