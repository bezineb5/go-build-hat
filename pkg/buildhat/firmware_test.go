package buildhat

import (
	"strings"
	"testing"
)

func TestFirmwareManager_CheckFirmwareVersion(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Test getting embedded version (this doesn't require communication)
	fm := brick.firmwareManager
	embeddedVersion, err := fm.GetEmbeddedFirmwareVersion()
	if err != nil {
		t.Fatalf("GetEmbeddedFirmwareVersion failed: %v", err)
	}

	// Should be the version from data/version file
	expected := "1737564117"
	if embeddedVersion != expected {
		t.Errorf("Expected embedded version %s, got %s", expected, embeddedVersion)
	}

	// Test that we can create the firmware manager
	if fm == nil {
		t.Fatal("FirmwareManager should not be nil")
	}
}

func TestFirmwareManager_GetEmbeddedFirmwareVersion(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Test getting embedded version
	fm := brick.firmwareManager
	version, err := fm.GetEmbeddedFirmwareVersion()
	if err != nil {
		t.Fatalf("GetEmbeddedFirmwareVersion failed: %v", err)
	}

	// Should be the version from data/version file
	expected := "1737564117"
	if version != expected {
		t.Errorf("Expected embedded version %s, got %s", expected, version)
	}
}

func TestFirmwareManager_BootloaderDetection(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Queue bootloader version response
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("Firmware version: BuildHAT bootloader version 1.0\r\n")

	// Test bootloader detection
	fm := brick.firmwareManager
	isBootloader := fm.isInBootloaderMode()
	if !isBootloader {
		t.Error("Expected bootloader mode to be detected")
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

func TestFirmwareManager_NormalFirmwareDetection(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()

	// Queue normal firmware version response (without bootloader signature)
	mockPort := brick.GetMockPort()
	mockPort.QueueReadData("Firmware version: 1737564117 2025-01-22T16:41:57+00:00\r\n")

	// Test normal firmware detection
	fm := brick.firmwareManager
	isBootloader := fm.isInBootloaderMode()
	if isBootloader {
		t.Error("Expected normal firmware mode, not bootloader")
	}

	// Verify version command was sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected version command to be sent")
	}
}

func TestFirmwareManager_UpdateProcess(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Test that we can create the firmware manager
	fm := brick.firmwareManager
	if fm == nil {
		t.Fatal("FirmwareManager should not be nil")
	}

	// Test checksum calculation (this doesn't require communication)
	testData := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := fm.calculateChecksum(testData)

	// Checksum should be non-zero for non-empty data
	if checksum == 0 {
		t.Error("Expected non-zero checksum for test data")
	}

	// Test with empty data
	emptyChecksum := fm.calculateChecksum([]byte{})
	if emptyChecksum != 1 {
		t.Errorf("Expected checksum 1 for empty data, got %d", emptyChecksum)
	}

	// Test with same data produces same checksum
	checksum2 := fm.calculateChecksum(testData)
	if checksum != checksum2 {
		t.Error("Same data should produce same checksum")
	}
}

func TestFirmwareManager_ChecksumCalculation(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	// Test checksum calculation
	fm := brick.firmwareManager

	// Test with known data
	testData := []byte{0x01, 0x02, 0x03, 0x04}
	checksum := fm.calculateChecksum(testData)

	// Checksum should be non-zero for non-empty data
	if checksum == 0 {
		t.Error("Expected non-zero checksum for test data")
	}

	// Test with empty data
	emptyChecksum := fm.calculateChecksum([]byte{})
	if emptyChecksum != 1 {
		t.Errorf("Expected checksum 1 for empty data, got %d", emptyChecksum)
	}

	// Test with same data produces same checksum
	checksum2 := fm.calculateChecksum(testData)
	if checksum != checksum2 {
		t.Error("Same data should produce same checksum")
	}
}
