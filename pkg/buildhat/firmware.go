package buildhat

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"time"
)

//go:embed data/*
var embeddedData embed.FS

// FirmwareManager handles firmware updates for the BuildHat
type FirmwareManager struct {
	brick *Brick
}

// NewFirmwareManager creates a new firmware manager
func NewFirmwareManager(brick *Brick) *FirmwareManager {
	return &FirmwareManager{
		brick: brick,
	}
}

// CheckAndUpdateFirmware checks if firmware update is needed and performs it
func (fm *FirmwareManager) CheckAndUpdateFirmware() error {
	fm.brick.logger.Info("Checking firmware status")

	// Clear the port first
	fm.brick.scanner.Scan()

	// Check for bootloader signature
	if fm.isInBootloaderMode() {
		fm.brick.logger.Info("Bootloader detected, updating firmware")
		return fm.updateFirmware()
	}

	fm.brick.logger.Info("Firmware is up to date")
	return nil
}

// isInBootloaderMode checks if the BuildHat is in bootloader mode
func (fm *FirmwareManager) isInBootloaderMode() bool {
	const bootloaderSignature = "BuildHAT bootloader version"

	// Use GetHardwareVersion which properly uses futures
	version, err := fm.brick.GetHardwareVersion()
	if err != nil {
		return false
	}

	// Check if the version string contains the bootloader signature
	return strings.Contains(version, bootloaderSignature)
}

// updateFirmware performs the firmware update process
func (fm *FirmwareManager) updateFirmware() error {
	fm.brick.logger.Info("Loading embedded firmware files")

	// Load firmware and signature from embedded files
	firmware, err := fm.loadEmbeddedFile("data/firmware.bin")
	if err != nil {
		return fmt.Errorf("failed to load firmware: %w", err)
	}

	signature, err := fm.loadEmbeddedFile("data/signature.bin")
	if err != nil {
		return fmt.Errorf("failed to load signature: %w", err)
	}

	fm.brick.logger.Info("Firmware loaded", "size", len(firmware), "signature_size", len(signature))

	// Step 1: Clear and get the prompt
	if err := fm.brick.writeCommand("clear"); err != nil {
		return fmt.Errorf("failed to clear: %w", err)
	}
	time.Sleep(100 * time.Millisecond)
	fm.brick.scanner.Scan()

	// Step 2: Load the firmware
	checksum := fm.calculateChecksum(firmware)
	loadCmd := fmt.Sprintf("load %d %d", len(firmware), checksum)
	if err := fm.brick.writeCommand(loadCmd); err != nil {
		return fmt.Errorf("failed to load firmware: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Send firmware data with STX/ETX markers
	_, _ = fm.brick.writer.Write([]byte{0x02}) // STX
	_, _ = fm.brick.writer.Write(firmware)
	_, _ = fm.brick.writer.Write([]byte{0x03}) // ETX
	_, _ = fm.brick.writer.Write([]byte("\r"))
	time.Sleep(10 * time.Millisecond)

	// Step 3: Load the signature
	fm.brick.scanner.Scan() // Clear any pending data

	sigCmd := fmt.Sprintf("signature %d", len(signature))
	if err := fm.brick.writeCommand(sigCmd); err != nil {
		return fmt.Errorf("failed to load signature: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	// Send signature data with STX/ETX markers
	_, _ = fm.brick.writer.Write([]byte{0x02}) // STX
	_, _ = fm.brick.writer.Write(signature)
	_, _ = fm.brick.writer.Write([]byte{0x03}) // ETX
	_, _ = fm.brick.writer.Write([]byte("\r"))
	time.Sleep(10 * time.Millisecond)

	// Step 4: Reboot
	if err := fm.brick.writeCommand("reboot"); err != nil {
		return fmt.Errorf("failed to reboot: %w", err)
	}
	time.Sleep(100 * time.Millisecond)
	fm.brick.scanner.Scan()

	// Wait for boot to complete
	time.Sleep(1500 * time.Millisecond)

	fm.brick.logger.Info("Firmware update completed successfully")
	return nil
}

// loadEmbeddedFile loads a file from the embedded filesystem
func (fm *FirmwareManager) loadEmbeddedFile(path string) ([]byte, error) {
	file, err := embeddedData.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded file %s: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded file %s: %w", path, err)
	}

	return data, nil
}

// calculateChecksum calculates the firmware checksum using the same algorithm as the C# implementation
func (fm *FirmwareManager) calculateChecksum(firmware []byte) uint32 {
	var check uint32 = 1

	for _, b := range firmware {
		if (check & 0x80000000) != 0 {
			check = (check << 1) ^ 0x1d872b41
		} else {
			check <<= 1
		}

		check = (check ^ uint32(b)) & 0xFFFFFFFF
	}

	return check
}

// GetEmbeddedFirmwareVersion returns the version of the embedded firmware
func (fm *FirmwareManager) GetEmbeddedFirmwareVersion() (string, error) {
	versionData, err := fm.loadEmbeddedFile("data/version")
	if err != nil {
		return "", fmt.Errorf("failed to load embedded version: %w", err)
	}

	version := strings.TrimSpace(string(versionData))
	return version, nil
}

// CheckFirmwareVersion compares current firmware version with embedded version
func (fm *FirmwareManager) CheckFirmwareVersion() (bool, error) {
	embeddedVersion, err := fm.GetEmbeddedFirmwareVersion()
	if err != nil {
		return false, fmt.Errorf("failed to get embedded version: %w", err)
	}

	currentVersion, err := fm.brick.GetHardwareVersion()
	if err != nil {
		return false, fmt.Errorf("failed to get current version: %w", err)
	}

	// Extract just the numeric part from current version (before any timestamp)
	currentNumeric := strings.Fields(currentVersion)[0]

	fm.brick.logger.Info("Version comparison",
		"embedded", embeddedVersion,
		"current_full", currentVersion,
		"current_numeric", currentNumeric)

	return embeddedVersion == currentNumeric, nil
}
