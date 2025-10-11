package buildhat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Brick represents a BuildHat device
type Brick struct {
	input   io.Reader
	writer  io.Writer
	scanner *bufio.Scanner
	logger  *slog.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.RWMutex

	// Connection state
	connections    [NumPorts]*Connection
	vinFutures     []chan float64
	versionFutures []chan string
	sensorFutures  [NumPorts][]chan []any // Sensor data futures per port
	rampFutures    [NumPorts][]chan bool  // Ramp completion futures per port
	pulseFutures   [NumPorts][]chan bool  // Pulse completion futures per port

	// Firmware management
	firmwareManager *FirmwareManager
}

// Connection represents a port connection
type Connection struct {
	TypeID     int
	Connected  bool
	SimpleMode int
	CombiMode  int
	Data       []any
}

// NewBrick creates a new BuildHat instance
func NewBrick(reader io.Reader, writer io.Writer, logger *slog.Logger) *Brick {
	if logger == nil {
		logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	brick := &Brick{
		input:          reader,
		writer:         writer,
		logger:         logger,
		ctx:            ctx,
		cancel:         cancel,
		vinFutures:     make([]chan float64, 0),
		versionFutures: make([]chan string, 0),
	}

	// Initialize sensor futures for each port
	for i := range NumPorts {
		brick.sensorFutures[i] = make([]chan []any, 0)
		brick.rampFutures[i] = make([]chan bool, 0)
		brick.pulseFutures[i] = make([]chan bool, 0)
	}

	// Initialize firmware manager
	brick.firmwareManager = NewFirmwareManager(brick)

	// Initialize connections
	for i := range NumPorts {
		brick.connections[i] = &Connection{
			TypeID:    -1,
			Connected: false,
		}
	}

	// Create scanner from input
	brick.scanner = bufio.NewScanner(reader)

	// Start the reader thread
	brick.wg.Add(1)
	go brick.reader()

	return brick
}

// Initialize initializes the BuildHat after creation
func (b *Brick) Initialize() error {
	b.logger.Info("Initializing BuildHat...")

	// Wait a moment for the reader thread to start
	time.Sleep(500 * time.Millisecond)

	// Check and update firmware if needed
	if err := b.firmwareManager.CheckAndUpdateFirmware(); err != nil {
		b.logger.Error("Firmware update failed", "error", err)
		return fmt.Errorf("firmware update failed: %w", err)
	}

	// Send version command to check state
	if err := b.writeCommand(Version()); err != nil {
		return err
	}

	// Wait a bit for initialization
	time.Sleep(2 * time.Second)

	// Send list command to scan devices
	if err := b.writeCommand(List()); err != nil {
		return err
	}

	// Wait for device scanning
	time.Sleep(3 * time.Second)

	b.logger.Info("BuildHat initialized")
	return nil
}

// Close closes the BuildHat connection
func (b *Brick) Close() error {
	// Close the serial port first to unblock the scanner
	if closer, ok := b.writer.(io.Closer); ok {
		closer.Close()
	}

	// Then cancel the context and wait for goroutines to finish
	b.cancel()
	b.wg.Wait()

	return nil
}

// reader is the main serial data reader thread
func (b *Brick) reader() {
	defer b.wg.Done()

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			if b.scanner.Scan() {
				line := strings.TrimSpace(b.scanner.Text())
				b.parseLine(line)
			} else {
				// Check for scanner error
				if err := b.scanner.Err(); err != nil {
					b.logger.Error("Scanner error", "error", err)
				}
				// Small delay to prevent busy waiting
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

// parseLine parses incoming serial data
func (b *Brick) parseLine(line string) {
	// Log all received lines for debugging
	if line != "" {
		b.logger.Debug("RX", "line", line)
	}

	if b.tryParsePortMessage(line) {
		return
	}

	if b.tryParseVoltageReading(line) {
		return
	}

	if b.tryParseVersionResponse(line) {
		return
	}

	if b.tryParseSensorData(line) {
		return
	}

	b.logger.Debug("Unhandled line", "line", line)
}

// tryParsePortMessage attempts to parse port connection messages
func (b *Brick) tryParsePortMessage(line string) bool {
	if len(line) < 3 || line[0] != 'P' || line[2] != ':' {
		return false
	}

	portID := int(line[1] - '0')
	if portID < 0 || portID > 3 {
		return false
	}

	msg := line[2:]
	b.handlePortMessage(portID, msg)
	return true
}

// tryParseVoltageReading attempts to parse voltage readings
func (b *Brick) tryParseVoltageReading(line string) bool {
	if len(line) < 3 || !strings.HasSuffix(line, " V") {
		return false
	}

	parts := strings.Split(line, " ")
	if len(parts) < 1 {
		return false
	}

	voltage, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return false
	}

	b.handleVoltageReading(voltage)
	return true
}

// tryParseVersionResponse attempts to parse version responses
func (b *Brick) tryParseVersionResponse(line string) bool {
	// Handle firmware version responses
	if version, ok := strings.CutPrefix(line, "Firmware version: "); ok {
		b.handleVersionResponse(version)
		return true
	}

	// Handle bootloader version responses
	if strings.HasPrefix(line, "BuildHAT bootloader version") {
		b.handleVersionResponse(line)
		return true
	}

	return false
}

// tryParseSensorData attempts to parse sensor data
func (b *Brick) tryParseSensorData(line string) bool {
	if len(line) < 4 || line[0] != 'P' {
		return false
	}

	if line[2] != 'M' && line[2] != 'C' {
		return false
	}

	portID := int(line[1] - '0')
	if portID < 0 || portID > 3 {
		return false
	}

	b.handleSensorData(portID, line)
	return true
}

// handlePortMessage handles port connection/disconnection messages
func (b *Brick) handlePortMessage(portID int, msg string) {
	switch {
	case strings.Contains(msg, "ramp done"):
		// Handle ramp completion
		b.mu.Lock()
		if len(b.rampFutures[portID]) > 0 {
			// Pop the first future and signal completion
			future := b.rampFutures[portID][0]
			b.rampFutures[portID] = b.rampFutures[portID][1:]
			b.mu.Unlock()
			future <- true
			close(future)
		} else {
			b.mu.Unlock()
			b.logger.Debug("Received ramp done with no pending future", "port", portID)
		}
		return

	case strings.Contains(msg, "pulse done"):
		// Handle pulse completion
		b.mu.Lock()
		if len(b.pulseFutures[portID]) > 0 {
			// Pop the first future and signal completion
			future := b.pulseFutures[portID][0]
			b.pulseFutures[portID] = b.pulseFutures[portID][1:]
			b.mu.Unlock()
			future <- true
			close(future)
		} else {
			b.mu.Unlock()
			b.logger.Debug("Received pulse done with no pending future", "port", portID)
		}
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	switch {
	case strings.Contains(msg, "connected to active ID"):
		// Extract type ID from message (e.g., "P3: connected to active ID 3D")
		parts := strings.Split(msg, " ")
		if len(parts) >= 6 {
			hexStr := parts[5] // The type ID is the 6th part (index 5)
			if typeID, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
				b.connections[portID].TypeID = int(typeID)
				b.connections[portID].Connected = true
			} else {
				b.logger.Error("Failed to parse type ID", "port", portID, "hex", hexStr, "error", err)
			}
		}
	case strings.Contains(msg, "connected to passive ID"):
		// Extract type ID from message
		parts := strings.Split(msg, " ")
		if len(parts) >= 6 {
			hexStr := parts[5] // The type ID is the 6th part (index 5)
			if typeID, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
				b.connections[portID].TypeID = int(typeID)
				b.connections[portID].Connected = true
			} else {
				b.logger.Error("Failed to parse passive type ID", "port", portID, "hex", hexStr, "error", err)
			}
		}
	case strings.Contains(msg, "disconnected"):
		b.connections[portID].TypeID = -1
		b.connections[portID].Connected = false
	case strings.Contains(msg, "no device detected"):
		b.connections[portID].TypeID = -1
		b.connections[portID].Connected = false
	}
}

// handleVoltageReading handles voltage readings
func (b *Brick) handleVoltageReading(voltage float64) {
	if len(b.vinFutures) > 0 {
		future := b.vinFutures[0]
		b.vinFutures = b.vinFutures[1:]
		select {
		case future <- voltage:
		default:
		}
	}
}

// handleVersionResponse handles version responses
func (b *Brick) handleVersionResponse(version string) {
	if len(b.versionFutures) > 0 {
		future := b.versionFutures[0]
		b.versionFutures = b.versionFutures[1:]
		select {
		case future <- version:
		default:
		}
	}
}

// handleSensorData handles sensor data
func (b *Brick) handleSensorData(portID int, line string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Ensure line is long enough to slice
	if len(line) < 5 {
		b.logger.Debug("Sensor data too short", "port", portID, "line", line)
		return
	}

	// Parse sensor data (simplified)
	parts := strings.Split(line[5:], " ")
	data := make([]any, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}
		if strings.Contains(part, ".") {
			if val, err := strconv.ParseFloat(part, 64); err == nil {
				data = append(data, val)
			}
		} else {
			if val, err := strconv.ParseInt(part, 10, 32); err == nil {
				data = append(data, int(val))
			}
		}
	}

	b.connections[portID].Data = data
	b.logger.Debug("Sensor data", "port", portID, "data", data)

	// Notify any waiting sensor futures
	if len(b.sensorFutures[portID]) > 0 {
		future := b.sensorFutures[portID][0]
		b.sensorFutures[portID] = b.sensorFutures[portID][1:]
		select {
		case future <- data:
		default:
		}
	}
}

// writeCommand sends a command to the BuildHat
func (b *Brick) writeCommand(command Command) error {
	cmd := command.CommandString()
	if !strings.HasSuffix(cmd, "\r") {
		cmd += "\r"
	}

	b.logger.Debug("TX", "cmd", strings.TrimSuffix(cmd, "\r"))
	_, err := b.writer.Write([]byte(cmd))
	return err
}

// GetHardwareVersion gets the hardware version
func (b *Brick) GetHardwareVersion() (string, error) {
	future := make(chan string, 1)
	b.versionFutures = append(b.versionFutures, future)

	if err := b.writeCommand(Version()); err != nil {
		return "", err
	}

	select {
	case version := <-future:
		return version, nil
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("timeout waiting for version response")
	}
}

// GetVoltage gets the input voltage
func (b *Brick) GetVoltage() (float64, error) {
	future := make(chan float64, 1)
	b.vinFutures = append(b.vinFutures, future)

	if err := b.writeCommand(Vin()); err != nil {
		return 0, err
	}

	select {
	case voltage := <-future:
		return voltage, nil
	case <-time.After(5 * time.Second):
		return 0, fmt.Errorf("timeout waiting for voltage response")
	}
}

// ScanDevices scans for connected devices
func (b *Brick) ScanDevices() error {
	return b.writeCommand(List())
}

// getSensorData waits for sensor data from a specific port
func (b *Brick) getSensorData(port int) ([]any, error) {
	b.mu.Lock()

	// Check if we already have cached data
	if len(b.connections[port].Data) > 0 {
		data := b.connections[port].Data
		// Clear cached data so next call gets fresh data
		b.connections[port].Data = nil
		b.mu.Unlock()
		return data, nil
	}

	// No cached data, create a future and wait for new data
	future := make(chan []any, 1)
	b.sensorFutures[port] = append(b.sensorFutures[port], future)
	b.mu.Unlock()

	select {
	case data := <-future:
		return data, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout waiting for sensor data on port %d", port)
	}
}

// GetDeviceInfo returns information about devices on all ports
func (b *Brick) GetDeviceInfo() map[Port]DeviceInfo {
	b.mu.RLock()
	defer b.mu.RUnlock()

	devices := make(map[Port]DeviceInfo)
	for i := range NumPorts {
		port := Port(i)
		conn := b.connections[i]

		devices[port] = DeviceInfo{
			Port:      port,
			TypeID:    conn.TypeID,
			Connected: conn.Connected,
			Name:      getDeviceName(conn.TypeID),
			Category:  getDeviceCategory(conn.TypeID),
		}
	}

	return devices
}

// DeviceInfo represents information about a device
type DeviceInfo struct {
	Port      Port
	TypeID    int
	Connected bool
	Name      string
	Category  DeviceCategory
}

// GetEmbeddedFirmwareVersion returns the version of the embedded firmware
func (b *Brick) GetEmbeddedFirmwareVersion() (string, error) {
	return b.firmwareManager.GetEmbeddedFirmwareVersion()
}

// CheckFirmwareVersion compares current firmware version with embedded version
func (b *Brick) CheckFirmwareVersion() (bool, error) {
	return b.firmwareManager.CheckFirmwareVersion()
}

// UpdateFirmware manually triggers a firmware update
func (b *Brick) UpdateFirmware() error {
	return b.firmwareManager.CheckAndUpdateFirmware()
}
