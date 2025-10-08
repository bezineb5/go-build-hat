package buildhat

import (
	"io"
	"log/slog"
	"testing"
)

// TestBrick creates a Brick instance with a mock serial port for testing
func TestBrick(_ *testing.T) *Brick {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors in tests
	}))

	// Create mock serial port
	mockPort := NewMockSerialPort(logger)

	// Use the normal NewBrick constructor which starts the reader goroutine
	brick := NewBrick(mockPort, mockPort, logger)

	return brick
}

// GetMockPort returns the mock serial port from a test brick
func (b *Brick) GetMockPort() *MockSerialPort {
	return b.input.(*MockSerialPort)
}

// CleanupTestBrick cleans up a test brick
func CleanupTestBrick(brick *Brick) {
	if brick != nil {
		brick.Close()
	}
}
