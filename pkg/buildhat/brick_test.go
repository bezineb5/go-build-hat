package buildhat

import (
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestNewBrick(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	brick := NewBrick(mockPort, mockPort, logger)

	if brick == nil {
		t.Fatal("NewBrick returned nil")
	}

	if brick.input != mockPort {
		t.Error("Expected input to be set to mockPort")
	}

	if brick.writer != mockPort {
		t.Error("Expected writer to be set to mockPort")
	}

	if brick.logger != logger {
		t.Error("Expected logger to be set")
	}

	if brick.scanner == nil {
		t.Error("Expected scanner to be initialized")
	}

	// Check that connections are initialized
	for i := range NumPorts {
		if brick.connections[i] == nil {
			t.Errorf("Expected connection %d to be initialized", i)
		}
		if brick.connections[i].TypeID != -1 {
			t.Errorf("Expected connection %d TypeID to be -1, got %d", i, brick.connections[i].TypeID)
		}
		if brick.connections[i].Connected {
			t.Errorf("Expected connection %d to be disconnected", i)
		}
	}

	// Cleanup
	brick.Close()
}

func TestNewBrick_WithNilLogger(t *testing.T) {
	mockPort := NewMockSerialPort(nil)
	brick := NewBrick(mockPort, mockPort, nil)

	if brick.logger == nil {
		t.Error("Expected logger to be set to default logger")
	}

	brick.Close()
}

func TestBrick_Initialize(t *testing.T) {
	t.Skip("Skipping slow test - Initialize has 5.5s of hardcoded sleeps for real hardware timing")
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Queue responses for initialization in the format expected by parseLine
	mockPort := brick.GetMockPort()
	// Initialize calls: isInBootloaderMode (1 version), GetHardwareVersion (1 version), list
	mockPort.QueueReadData("Firmware version: 1737564117 2025-01-22T16:41:57+00:00\r\n") // For isInBootloaderMode
	mockPort.QueueReadData("Firmware version: 1737564117 2025-01-22T16:41:57+00:00\r\n") // For GetHardwareVersion
	mockPort.QueueReadData("P0: connected to active ID 4B\r\n")
	mockPort.QueueReadData("type 4B\r\n")
	mockPort.QueueReadData("P1: no device detected\r\n")
	mockPort.QueueReadData("P2: no device detected\r\n")
	mockPort.QueueReadData("P3: no device detected\r\n")

	err := brick.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 3 {
		t.Fatalf("Expected at least 3 commands (version for bootloader check, version, list), got %d", len(writeHistory))
	}

	// The firmware manager checks for bootloader mode by sending version
	if !strings.Contains(writeHistory[0], "version") {
		t.Errorf("Expected first version command (bootloader check): %s", writeHistory[0])
	}

	// Then Initialize sends version and list
	if !strings.Contains(writeHistory[1], "version") {
		t.Errorf("Expected second version command: %s", writeHistory[1])
	}
	if !strings.Contains(writeHistory[2], "list") {
		t.Errorf("Expected list command: %s", writeHistory[2])
	}
}

func TestBrick_Close(t *testing.T) {
	brick := TestBrick(t)

	// Close should not panic
	err := brick.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}

	// Calling close again should not panic
	err = brick.Close()
	if err != nil {
		t.Errorf("Second close returned error: %v", err)
	}
}

func TestBrick_GetHardwareVersion(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Queue version response in the format expected by parseLine
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("Firmware version: 1737564117 2025-01-22T16:41:57+00:00\r\n")

	version, err := brick.GetHardwareVersion()
	if err != nil {
		t.Fatalf("GetHardwareVersion failed: %v", err)
	}

	if version != "1737564117 2025-01-22T16:41:57+00:00" {
		t.Errorf("Expected version '1737564117 2025-01-22T16:41:57+00:00', got '%s'", version)
	}

	// Verify version command was sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected version command to be sent")
	}

	if !strings.Contains(writeHistory[0], "version") {
		t.Errorf("Expected version command: %s", writeHistory[0])
	}
}

func TestBrick_GetVoltage(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Queue voltage response
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("vin\r\n")
	mockPort.QueueReadData("8.2 V\r\n")

	voltage, err := brick.GetVoltage()
	if err != nil {
		t.Fatalf("GetVoltage failed: %v", err)
	}

	if voltage != 8.2 {
		t.Errorf("Expected voltage 8.2, got %f", voltage)
	}

	// Verify vin command was sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected vin command to be sent")
	}

	if !strings.Contains(writeHistory[0], "vin") {
		t.Errorf("Expected vin command: %s", writeHistory[0])
	}
}

func TestBrick_ScanDevices(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	err := brick.ScanDevices()
	if err != nil {
		t.Fatalf("ScanDevices failed: %v", err)
	}

	// Verify list command was sent
	mockPort := brick.GetMockPort()
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected list command to be sent")
	}

	if !strings.Contains(writeHistory[0], "list") {
		t.Errorf("Expected list command: %s", writeHistory[0])
	}
}

func TestBrick_GetDeviceInfo(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	devices := brick.GetDeviceInfo()

	if len(devices) != NumPorts {
		t.Errorf("Expected %d devices, got %d", NumPorts, len(devices))
	}

	// Check that all ports are present
	expectedPorts := []Port{PortA, PortB, PortC, PortD}
	for _, port := range expectedPorts {
		if _, exists := devices[port]; !exists {
			t.Errorf("Expected port %s to be present", port)
		}
	}

	// All should be disconnected initially
	for port, info := range devices {
		if info.Connected {
			t.Errorf("Expected port %s to be disconnected", port)
		}
		if info.TypeID != -1 {
			t.Errorf("Expected port %s TypeID to be -1, got %d", port, info.TypeID)
		}
	}
}

func TestBrick_GetDeviceInfo_WithConnectedDevices(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Simulate connected devices
	brick.mu.Lock()
	brick.connections[0].TypeID = 75 // Medium Angular Motor
	brick.connections[0].Connected = true
	brick.connections[1].TypeID = 61 // Color Sensor
	brick.connections[1].Connected = true
	brick.mu.Unlock()

	devices := brick.GetDeviceInfo()

	// Check port A (motor)
	portA := devices[PortA]
	if !portA.Connected {
		t.Error("Expected port A to be connected")
	}
	if portA.TypeID != 75 {
		t.Errorf("Expected port A TypeID to be 75, got %d", portA.TypeID)
	}
	if portA.Category != DeviceCategoryMotor {
		t.Errorf("Expected port A Category to be Motor, got %s", portA.Category)
	}

	// Check port B (sensor)
	portB := devices[PortB]
	if !portB.Connected {
		t.Error("Expected port B to be connected")
	}
	if portB.TypeID != 61 {
		t.Errorf("Expected port B TypeID to be 61, got %d", portB.TypeID)
	}
	if portB.Category != DeviceCategorySensor {
		t.Errorf("Expected port B Category to be Sensor, got %s", portB.Category)
	}
}

func TestBrick_GetEmbeddedFirmwareVersion(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	version, err := brick.GetEmbeddedFirmwareVersion()
	if err != nil {
		t.Fatalf("GetEmbeddedFirmwareVersion failed: %v", err)
	}

	if version == "" {
		t.Error("Expected non-empty firmware version")
	}
}

func TestBrick_CheckFirmwareVersion(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Queue version response in the format expected by parseLine
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("Firmware version: 1737564117 2025-01-22T16:41:57+00:00\r\n")

	match, err := brick.CheckFirmwareVersion()
	if err != nil {
		t.Fatalf("CheckFirmwareVersion failed: %v", err)
	}

	// Should match since we're simulating the same version
	if !match {
		t.Error("Expected firmware versions to match")
	}
}
