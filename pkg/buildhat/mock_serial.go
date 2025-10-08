package buildhat

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"go.bug.st/serial"
)

// MockSerialPort implements a mock serial port for testing
type MockSerialPort struct {
	mu           sync.RWMutex
	readBuffer   []byte
	writeBuffer  []byte
	writeHistory []string
	readHistory  []string
	closed       bool
	logger       *slog.Logger
}

// NewMockSerialPort creates a new mock serial port
func NewMockSerialPort(logger *slog.Logger) *MockSerialPort {
	if logger == nil {
		logger = slog.Default()
	}
	return &MockSerialPort{
		readBuffer:   make([]byte, 0),
		writeBuffer:  make([]byte, 0),
		writeHistory: make([]string, 0),
		readHistory:  make([]string, 0),
		logger:       logger,
	}
}

// Write implements io.Writer
func (m *MockSerialPort) Write(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return 0, fmt.Errorf("port is closed")
	}

	m.writeBuffer = append(m.writeBuffer, data...)
	m.writeHistory = append(m.writeHistory, string(data))

	m.logger.Debug("MockSerialPort.Write", "data", string(data), "hex", fmt.Sprintf("%x", data))
	return len(data), nil
}

// Read implements io.Reader
func (m *MockSerialPort) Read(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return 0, fmt.Errorf("port is closed")
	}

	if len(m.readBuffer) == 0 {
		return 0, io.EOF
	}

	n := copy(data, m.readBuffer)
	m.readBuffer = m.readBuffer[n:]

	m.logger.Debug("MockSerialPort.Read", "data", string(data[:n]), "hex", fmt.Sprintf("%x", data[:n]))
	return n, nil
}

// Close implements io.Closer
func (m *MockSerialPort) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	m.logger.Debug("MockSerialPort.Close")
	return nil
}

// SetReadTimeout sets a mock read timeout (not implemented in real serial)
func (m *MockSerialPort) SetReadTimeout(_ time.Duration) error {
	return nil
}

// SetWriteTimeout sets a mock write timeout (not implemented in real serial)
func (m *MockSerialPort) SetWriteTimeout(_ time.Duration) error {
	return nil
}

// Break sends a break signal (not implemented in mock)
func (m *MockSerialPort) Break(_ time.Duration) error {
	return nil
}

// Drain waits until all data in the buffer are sent (not implemented in mock)
func (m *MockSerialPort) Drain() error {
	return nil
}

// ResetInputBuffer purges port read buffer (not implemented in mock)
func (m *MockSerialPort) ResetInputBuffer() error {
	return nil
}

// ResetOutputBuffer purges port write buffer (not implemented in mock)
func (m *MockSerialPort) ResetOutputBuffer() error {
	return nil
}

// SetDTR sets the modem status bit DataTerminalReady (not implemented in mock)
func (m *MockSerialPort) SetDTR(_ bool) error {
	return nil
}

// SetRTS sets the modem status bit RequestToSend (not implemented in mock)
func (m *MockSerialPort) SetRTS(_ bool) error {
	return nil
}

// GetModemStatusBits returns modem status bits (not implemented in mock)
func (m *MockSerialPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	return nil, nil
}

// SetMode sets all parameters of the serial port (not implemented in mock)
func (m *MockSerialPort) SetMode(_ *serial.Mode) error {
	return nil
}

// Reset resets the mock port state
func (m *MockSerialPort) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.readBuffer = make([]byte, 0)
	m.writeBuffer = make([]byte, 0)
	m.writeHistory = make([]string, 0)
	m.readHistory = make([]string, 0)
	m.closed = false
}

// QueueReadData queues data to be read by the mock port
func (m *MockSerialPort) QueueReadData(data string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.readBuffer = append(m.readBuffer, []byte(data)...)
	m.readHistory = append(m.readHistory, data)
	m.logger.Debug("MockSerialPort.QueueReadData", "data", data)
}

// GetWriteHistory returns all data written to the port
func (m *MockSerialPort) GetWriteHistory() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]string{}, m.writeHistory...)
}

// GetReadHistory returns all data read from the port
func (m *MockSerialPort) GetReadHistory() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]string{}, m.readHistory...)
}

// GetLastWrite returns the last data written to the port
func (m *MockSerialPort) GetLastWrite() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.writeHistory) == 0 {
		return ""
	}
	return m.writeHistory[len(m.writeHistory)-1]
}

// GetWriteCount returns the number of write operations
func (m *MockSerialPort) GetWriteCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.writeHistory)
}

// SimulateDeviceResponse simulates a device response to a command
func (m *MockSerialPort) SimulateDeviceResponse(command, response string) {
	// Queue the echo of the command
	m.QueueReadData(command + "\r\n")
	// Queue the response
	m.QueueReadData(response + "\r\n")
}

// SimulateFirmwareBootloader simulates bootloader mode responses
func (m *MockSerialPort) SimulateFirmwareBootloader() {
	m.QueueReadData("BuildHAT bootloader version 1.0\r\n")
}

// SimulateFirmwareNormal simulates normal firmware mode responses
func (m *MockSerialPort) SimulateFirmwareNormal() {
	m.QueueReadData("version\r\n")
	m.QueueReadData("1737564117 2025-01-22T16:41:57+00:00\r\n")
}

// SimulateDeviceList simulates device list responses
func (m *MockSerialPort) SimulateDeviceList() {
	m.QueueReadData("list\r\n")
	m.QueueReadData("P0: connected to active ID 4B\r\n")
	m.QueueReadData("type 4B\r\n")
	m.QueueReadData("nmodes =5\r\n")
	m.QueueReadData("P1: connected to passive ID 1A\r\n")
	m.QueueReadData("type 1A\r\n")
	m.QueueReadData("P2: no device detected\r\n")
	m.QueueReadData("P3: no device detected\r\n")
}

// SimulateMotorResponse simulates motor control responses
func (m *MockSerialPort) SimulateMotorResponse(port string, speed int) {
	cmd := fmt.Sprintf("port %s ; set %d", port, speed)
	m.QueueReadData(cmd + "\r\n")
	m.QueueReadData("OK\r\n")
}

// SimulateSensorResponse simulates sensor reading responses
func (m *MockSerialPort) SimulateSensorResponse(port string, mode int, value string) {
	// Convert port letter to number (A=0, B=1, C=2, D=3)
	portNum := 0
	if port != "" && port[0] >= 'A' && port[0] <= 'D' {
		portNum = int(port[0] - 'A')
	} else if port != "" && port[0] >= '0' && port[0] <= '3' {
		portNum = int(port[0] - '0')
	}

	// Queue the sensor data in the format expected by parseLine
	// Format: P<port>M<mode> <data>
	m.QueueReadData(fmt.Sprintf("P%dM%d %s\r\n", portNum, mode, value))
}
