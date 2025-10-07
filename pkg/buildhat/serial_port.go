package buildhat

import (
	"fmt"
	"log/slog"
	"time"

	"go.bug.st/serial"
)

// NewSerialPort creates a new serial port connection
func NewSerialPort(portPath string) (serial.Port, error) {
	// Configure serial port settings for BuildHat
	mode := &serial.Mode{
		BaudRate: 115200,            // BuildHat uses 115200 baud
		DataBits: 8,                 // 8 data bits
		Parity:   serial.NoParity,   // No parity
		StopBits: serial.OneStopBit, // 1 stop bit
	}

	// Open the serial port
	port, err := serial.Open(portPath, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %w", portPath, err)
	}

	// Configure port settings
	if err := port.SetReadTimeout(100 * time.Millisecond); err != nil {
		return nil, fmt.Errorf("failed to set read timeout: %w", err)
	}

	return port, nil
}

// GetAvailablePorts returns a list of available serial ports
func GetAvailablePorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("failed to get serial ports list: %w", err)
	}
	return ports, nil
}

// DetectBuildHatPort attempts to detect the BuildHat serial port
func DetectBuildHatPort(logger *slog.Logger) (string, error) {
	ports, err := GetAvailablePorts()
	if err != nil {
		return "", err
	}

	logger.Info("Scanning for BuildHat on available ports", "ports", ports)

	// Common BuildHat port names on Raspberry Pi
	buildhatPorts := []string{
		"/dev/serial0", // Primary UART on Raspberry Pi
		"/dev/ttyAMA0", // Alternative UART name
	}

	// Check common BuildHat ports first
	for _, portName := range buildhatPorts {
		for _, availablePort := range ports {
			if availablePort == portName {
				logger.Info("Found potential BuildHat port", "port", portName)
				return portName, nil
			}
		}
	}

	// If no common ports found, return the first available port
	if len(ports) > 0 {
		logger.Info("Using first available port", "port", ports[0])
		return ports[0], nil
	}

	return "", fmt.Errorf("no serial ports found")
}
