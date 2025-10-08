package buildhat

import (
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestMockSerialPort_Read(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Test reading when no data is queued
	buffer := make([]byte, 10)
	n, err := mockPort.Read(buffer)
	if err != nil {
		t.Errorf("Expected no error when no data is available, got %v", err)
	}
	if n != 0 {
		t.Errorf("Expected 0 bytes read, got %d", n)
	}

	// Queue some data
	mockPort.QueueReadData("test data")

	// Test reading queued data
	n, err = mockPort.Read(buffer)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}
	if n == 0 {
		t.Error("Expected to read some data")
	}
}

func TestMockSerialPort_Read_WhenClosed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	mockPort.Close()

	buffer := make([]byte, 10)
	_, err := mockPort.Read(buffer)
	if err == nil {
		t.Error("Expected error when reading from closed port")
	}
}

func TestMockSerialPort_Write_WhenClosed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	mockPort.Close()

	_, err := mockPort.Write([]byte("test"))
	if err == nil {
		t.Error("Expected error when writing to closed port")
	}
}

func TestMockSerialPort_SetReadTimeout(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.SetReadTimeout(100 * time.Millisecond)
	if err != nil {
		t.Errorf("SetReadTimeout failed: %v", err)
	}
}

func TestMockSerialPort_SetWriteTimeout(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.SetWriteTimeout(100 * time.Millisecond)
	if err != nil {
		t.Errorf("SetWriteTimeout failed: %v", err)
	}
}

func TestMockSerialPort_Break(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.Break(100 * time.Millisecond)
	if err != nil {
		t.Errorf("Break failed: %v", err)
	}
}

func TestMockSerialPort_Drain(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.Drain()
	if err != nil {
		t.Errorf("Drain failed: %v", err)
	}
}

func TestMockSerialPort_ResetInputBuffer(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.ResetInputBuffer()
	if err != nil {
		t.Errorf("ResetInputBuffer failed: %v", err)
	}
}

func TestMockSerialPort_ResetOutputBuffer(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.ResetOutputBuffer()
	if err != nil {
		t.Errorf("ResetOutputBuffer failed: %v", err)
	}
}

func TestMockSerialPort_SetDTR(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.SetDTR(true)
	if err != nil {
		t.Errorf("SetDTR failed: %v", err)
	}

	err = mockPort.SetDTR(false)
	if err != nil {
		t.Errorf("SetDTR(false) failed: %v", err)
	}
}

func TestMockSerialPort_SetRTS(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.SetRTS(true)
	if err != nil {
		t.Errorf("SetRTS failed: %v", err)
	}

	err = mockPort.SetRTS(false)
	if err != nil {
		t.Errorf("SetRTS(false) failed: %v", err)
	}
}

func TestMockSerialPort_GetModemStatusBits(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	bits, err := mockPort.GetModemStatusBits()
	if err != nil {
		t.Errorf("GetModemStatusBits failed: %v", err)
	}
	if bits != nil {
		t.Error("Expected nil modem status bits")
	}
}

func TestMockSerialPort_SetMode(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	err := mockPort.SetMode(nil)
	if err != nil {
		t.Errorf("SetMode failed: %v", err)
	}
}

func TestMockSerialPort_Reset(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Write some data
	if _, err := mockPort.Write([]byte("test")); err != nil {
		t.Errorf("Write failed: %v", err)
	}

	// Queue some read data
	mockPort.QueueReadData("test read")

	// Reset the port
	mockPort.Reset()

	// Verify write history is cleared
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) != 0 {
		t.Errorf("Expected write history to be cleared, got %d entries", len(writeHistory))
	}

	// Verify read buffer is cleared
	buffer := make([]byte, 10)
	n, err := mockPort.Read(buffer)
	if err != nil {
		t.Errorf("Expected no error when reading from reset port, got %v", err)
	}
	if n != 0 {
		t.Error("Expected 0 bytes after reset")
	}

	// Verify port is not closed after reset
	_, err = mockPort.Write([]byte("test"))
	if err != nil {
		t.Errorf("Write after reset failed: %v", err)
	}
}

func TestMockSerialPort_GetReadHistory(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Queue some read data
	mockPort.QueueReadData("first")
	mockPort.QueueReadData("second")

	// Get read history
	readHistory := mockPort.GetReadHistory()
	if len(readHistory) != 2 {
		t.Errorf("Expected 2 read history entries, got %d", len(readHistory))
	}

	if readHistory[0] != "first" {
		t.Errorf("Expected first read to be 'first', got '%s'", readHistory[0])
	}

	if readHistory[1] != "second" {
		t.Errorf("Expected second read to be 'second', got '%s'", readHistory[1])
	}
}

func TestMockSerialPort_GetLastWrite(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Test when no writes have occurred
	lastWrite := mockPort.GetLastWrite()
	if lastWrite != "" {
		t.Errorf("Expected empty last write, got '%s'", lastWrite)
	}

	// Write some data
	if _, err := mockPort.Write([]byte("first")); err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if _, err := mockPort.Write([]byte("second")); err != nil {
		t.Errorf("Write failed: %v", err)
	}

	// Get last write
	lastWrite = mockPort.GetLastWrite()
	if lastWrite != "second" {
		t.Errorf("Expected last write to be 'second', got '%s'", lastWrite)
	}
}

func TestMockSerialPort_GetWriteCount(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Test initial count
	count := mockPort.GetWriteCount()
	if count != 0 {
		t.Errorf("Expected initial write count to be 0, got %d", count)
	}

	// Write some data
	if _, err := mockPort.Write([]byte("test1")); err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if _, err := mockPort.Write([]byte("test2")); err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if _, err := mockPort.Write([]byte("test3")); err != nil {
		t.Errorf("Write failed: %v", err)
	}

	// Check count
	count = mockPort.GetWriteCount()
	if count != 3 {
		t.Errorf("Expected write count to be 3, got %d", count)
	}
}

func TestMockSerialPort_SimulateDeviceResponse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Simulate device response
	mockPort.SimulateDeviceResponse("test command", "test response")

	// Check read history
	readHistory := mockPort.GetReadHistory()
	if len(readHistory) != 2 {
		t.Errorf("Expected 2 read history entries, got %d", len(readHistory))
	}

	if readHistory[0] != "test command\r\n" {
		t.Errorf("Expected first read to be 'test command\\r\\n', got '%s'", readHistory[0])
	}

	if readHistory[1] != "test response\r\n" {
		t.Errorf("Expected second read to be 'test response\\r\\n', got '%s'", readHistory[1])
	}
}

func TestMockSerialPort_SimulateFirmwareBootloader(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Simulate bootloader
	mockPort.SimulateFirmwareBootloader()

	// Check read history
	readHistory := mockPort.GetReadHistory()
	if len(readHistory) != 1 {
		t.Errorf("Expected 1 read history entry, got %d", len(readHistory))
	}

	expected := "BuildHAT bootloader version 1.0\r\n"
	if readHistory[0] != expected {
		t.Errorf("Expected bootloader message '%s', got '%s'", expected, readHistory[0])
	}
}

func TestMockSerialPort_SimulateFirmwareNormal(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Simulate normal firmware
	mockPort.SimulateFirmwareNormal()

	// Check read history
	readHistory := mockPort.GetReadHistory()
	if len(readHistory) != 2 {
		t.Errorf("Expected 2 read history entries, got %d", len(readHistory))
	}

	if readHistory[0] != "version\r\n" {
		t.Errorf("Expected first read to be 'version\\r\\n', got '%s'", readHistory[0])
	}

	if readHistory[1] != "1737564117 2025-01-22T16:41:57+00:00\r\n" {
		t.Errorf("Expected second read to be version string, got '%s'", readHistory[1])
	}
}

func TestMockSerialPort_SimulateDeviceList(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	mockPort := NewMockSerialPort(logger)
	defer mockPort.Close()

	// Simulate device list
	mockPort.SimulateDeviceList()

	// Check read history
	readHistory := mockPort.GetReadHistory()
	if len(readHistory) != 8 {
		t.Errorf("Expected 8 read history entries, got %d", len(readHistory))
	}

	// Check first few entries
	expected := []string{
		"list\r\n",
		"P0: connected to active ID 4B\r\n",
		"type 4B\r\n",
		"nmodes =5\r\n",
	}

	for i, exp := range expected {
		if readHistory[i] != exp {
			t.Errorf("Expected read[%d] to be '%s', got '%s'", i, exp, readHistory[i])
		}
	}
}
