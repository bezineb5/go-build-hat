package buildhat

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
	"github.com/bezineb5/go-build-hat/pkg/buildhat/motors"
	"github.com/bezineb5/go-build-hat/pkg/buildhat/sensors"
)

//go:embed data/*
var embeddedData embed.FS

// SensorDataHandler interface for sensors that can receive data updates
type SensorDataHandler interface {
	// UpdateFromSensorData updates sensor values from raw sensor data
	UpdateFromSensorData(data []string) error
	// GetTriggerFlag returns a pointer to the trigger flag for this sensor
	GetTriggerFlag() *bool
}

// Brick represents the main Brick class allowing low level access to motors and sensors
type Brick struct {
	reader io.Reader
	writer io.Writer
	logger *slog.Logger

	// Internal state
	sensorTypes  [4]models.SensorType // 4 ports, can be any of motors, sensors, active elements
	sensors      [4]SensorDataHandler // Sensor data handlers for each port
	ledMode      models.LedMode
	inputVoltage float64

	// Communication
	scanner *bufio.Scanner
	mu      sync.RWMutex

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Colors mapping (from C# implementation)
	colors map[int]string
}

// NewBrick creates a new Brick instance with io.Reader and io.Writer interfaces
func NewBrick(reader io.Reader, writer io.Writer, logger *slog.Logger) *Brick {
	if logger == nil {
		logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	brick := &Brick{
		reader: reader,
		writer: writer,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		colors: map[int]string{
			-1: "black",
			0:  "black",
			1:  "brown",
			2:  "magenta",
			3:  "blue",
			4:  "cyan",
			5:  "palegreen",
			6:  "green",
			7:  "yellow",
			8:  "yellow",
			9:  "red",
			10: "white",
		},
	}

	// Initialize scanner for reading
	brick.scanner = bufio.NewScanner(reader)

	// Initialize the brick
	if err := brick.initialize(); err != nil {
		logger.Error("Failed to initialize brick", "error", err)
	}

	// Start the background processing
	brick.startRunning()

	return brick
}

// Close closes the brick and cleans up resources
func (b *Brick) Close() error {
	b.cancel()
	b.wg.Wait()
	return nil
}

// GetSensorType gets the sensor type connected at a specific port
func (b *Brick) GetSensorType(port models.SensorPort) models.SensorType {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sensorTypes[port.Byte()]
}

// SetSensorType sets the sensor type for a port (for testing purposes)
func (b *Brick) SetSensorType(port models.SensorPort, sensorType models.SensorType) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sensorTypes[port.Byte()] = sensorType
}

// SetLedMode sets the LED mode
func (b *Brick) SetLedMode(mode models.LedMode) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.ledMode = mode
	return b.setLedMode(mode)
}

// GetLedMode gets the current LED mode
func (b *Brick) GetLedMode() models.LedMode {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ledMode
}

// GetInputVoltage gets the input voltage
func (b *Brick) GetInputVoltage() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.inputVoltage
}

// initialize initializes the brick
func (b *Brick) initialize() error {
	// Clear the port first
	b.readExisting()

	// Check if there is firmware, if not, upload one
	if err := b.checkForFirmwareAndUpload(); err != nil {
		return fmt.Errorf("failed to check/upload firmware: %w", err)
	}

	b.readExisting()

	// Clear the output and find the version
	var line string

	for i := range 10 {
		if err := b.writeCommand("version\r"); err != nil {
			return fmt.Errorf("failed to write version command: %w", err)
		}

		// Wait longer for response after firmware upload
		time.Sleep(100 * time.Millisecond)

		line = b.readExisting()
		b.logger.Debug("Version response", "attempt", i+1, "response", line)

		if strings.Contains(line, "Firmware version: ") || strings.Contains(line, "version") {
			b.logger.Info("Version detected", "response", line)
			break
		}

		if i == 9 {
			b.logger.Warn("Could not read version after 10 retries, continuing anyway", "last_response", line)
			// Don't fail initialization if we can't read version
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	b.readExisting()

	// No echo and read the voltage
	if err := b.setEcho(false); err != nil {
		return fmt.Errorf("failed to set echo: %w", err)
	}

	// Try to read voltage quickly, but don't block initialization
	b.inputVoltage = 12.0 // Default voltage

	// Quick voltage reading attempt (non-blocking)
	go func() {
		time.Sleep(100 * time.Millisecond) // Brief delay
		if err := b.setLedMode(models.VoltageDependant); err == nil {
			time.Sleep(50 * time.Millisecond)
			b.readLine()
			rawVoltage, err := b.getRawVoltage()
			if err != nil {
				b.logger.Error("Failed to get raw voltage", "error", err)
				return
			}
			rawV := strings.Fields(rawVoltage)
			if len(rawV) > 0 {
				if voltage, err := strconv.ParseFloat(rawV[0], 64); err == nil {
					b.mu.Lock()
					b.inputVoltage = voltage
					b.mu.Unlock()
					b.logger.Info("Voltage read successfully", "voltage", voltage)
				}
			}
		}
	}()

	return nil
}

// startRunning starts the background processing thread
func (b *Brick) startRunning() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.running()
	}()
}

// running is the main background processing loop
func (b *Brick) running() {
	// Force an update
	if err := b.writeCommand("list\r"); err != nil {
		b.logger.Error("Failed to send list command", "error", err)
	}

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			if b.scanner.Scan() {
				line := b.scanner.Text()
				b.processOutput(line)
			} else {
				// Scanner error or EOF
				if err := b.scanner.Err(); err != nil {
					// Only log non-timeout errors
					if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "no data") {
						b.logger.Error("Scanner error", "error", err)
					}
				}
				// Wait longer to avoid excessive polling
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// processOutput processes output from the BuildHat
func (b *Brick) processOutput(line string) {
	if line == "" {
		return
	}

	b.logger.Debug("Processing output", "line", line)

	// Handle different types of output
	switch {
	case strings.Contains(line, ": connected to active ID ") || strings.Contains(line, ": connected to passive ID "):
		b.handleDeviceConnected(line)
	case strings.Contains(line, ": disconnected") || strings.Contains(line, ": timeout during data phase: disconnecting") || strings.Contains(line, ": no device detected"):
		b.handleDeviceDisconnected(line)
	case len(line) > 3 && line[0] == 'P' && (line[2] == 'C' || line[2] == 'M'):
		b.handleSensorData(line)
	case strings.Contains(line, ": button released"):
		b.handleButtonReleased(line)
	case strings.Contains(line, ": button pressed"):
		b.handleButtonPressed(line)
	case strings.Contains(line, "power fault"):
		b.handlePowerFault()
	}
}

// writeCommand writes a command to the BuildHat
func (b *Brick) writeCommand(command string) error {
	b.logger.Debug("Sending command", "command", strings.TrimSpace(command))
	_, err := b.writer.Write([]byte(command))
	return err
}

// readLine reads a line from the BuildHat
func (b *Brick) readLine() string {
	if b.scanner.Scan() {
		line := b.scanner.Text()
		b.logger.Debug("Received line", "line", line)
		return line
	}
	return ""
}

// readExisting reads all available data from the BuildHat
func (b *Brick) readExisting() string {
	var result strings.Builder

	// Read all available data from the serial port
	for {
		// Try to read a chunk of data
		buffer := make([]byte, 1024)
		n, err := b.reader.Read(buffer)

		if n > 0 {
			result.Write(buffer[:n])
		}

		if err != nil {
			// If it's EOF or no more data, break
			if err == io.EOF {
				break
			}
			// For other errors, log and break
			b.logger.Debug("Error reading from serial port", "error", err)
			break
		}

		// If we read less than the buffer size, we've likely read all available data
		if n < len(buffer) {
			break
		}
	}

	return result.String()
}

// Helper methods for initialization
func (b *Brick) checkForFirmwareAndUpload() error {
	const bootloaderSignature = "BuildHAT bootloader version"
	const getFirmwareVersion = "version\r"

	// Clear the port first
	b.readExisting()

	// We have to do it 2 times to make sure, the first time, it may not provide the proper elements
	if err := b.writeCommand(getFirmwareVersion); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	b.readExisting() // Clear first response

	if err := b.writeCommand(getFirmwareVersion); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	prompt := b.readExisting()

	// Check if we need to upload firmware
	if strings.Contains(prompt, bootloaderSignature) {
		b.logger.Info("Bootloader detected, uploading firmware")

		// Send carriage return to get prompt
		if err := b.writeCommand("\r"); err != nil {
			return fmt.Errorf("failed to send carriage return: %w", err)
		}
		b.readLine()
		b.readExisting()

		// Load firmware and signature files
		firmware, err := b.LoadFirmwareFile()
		if err != nil {
			return fmt.Errorf("failed to load firmware file: %w", err)
		}

		signature, err := b.LoadSignatureFile()
		if err != nil {
			return fmt.Errorf("failed to load signature file: %w", err)
		}

		// Step 1: clear and get the prompt
		if err := b.writeCommand("clear\r"); err != nil {
			return fmt.Errorf("failed to send clear command: %w", err)
		}
		b.readLine()
		b.readExisting()

		// Step 2: load the firmware
		checksum := b.GetFirmwareChecksum(firmware)
		if err := b.writeCommand(fmt.Sprintf("load %d %d\r", len(firmware), checksum)); err != nil {
			return fmt.Errorf("failed to send load command: %w", err)
		}
		b.readExisting()

		// Write firmware data with STX/ETX markers
		if err := b.writeBinaryData(firmware); err != nil {
			return fmt.Errorf("failed to write firmware data: %w", err)
		}
		time.Sleep(10 * time.Millisecond)

		// Step 3: load the signature
		b.readExisting()
		if err := b.writeCommand(fmt.Sprintf("signature %d\r", len(signature))); err != nil {
			return fmt.Errorf("failed to send signature command: %w", err)
		}

		// Write signature data with STX/ETX markers
		if err := b.writeBinaryData(signature); err != nil {
			return fmt.Errorf("failed to write signature data: %w", err)
		}
		time.Sleep(10 * time.Millisecond)
		b.readExisting()

		// Step 4: reboot
		if err := b.writeCommand("reboot\r"); err != nil {
			return fmt.Errorf("failed to send reboot command: %w", err)
		}
		b.readLine()
		b.readExisting()
		time.Sleep(1500 * time.Millisecond)

		b.logger.Info("Firmware upload completed")
	}

	return nil
}

func (b *Brick) setEcho(on bool) error {
	echoValue := "0"
	if on {
		echoValue = "1"
	}
	return b.writeCommand(fmt.Sprintf("echo %s\r", echoValue))
}

func (b *Brick) setLedMode(mode models.LedMode) error {
	return b.writeCommand(fmt.Sprintf("ledmode %d\r", int(mode)))
}

func (b *Brick) getRawVoltage() (string, error) {
	if err := b.writeCommand("vin\r"); err != nil {
		return "", fmt.Errorf("failed to send vin command: %w", err)
	}
	b.readLine()
	response := b.readLine()
	// Remove "V" suffix if present for consistent parsing
	response = strings.TrimSuffix(response, "V")
	return strings.TrimSpace(response), nil
}

// createSensorForType creates the appropriate sensor/motor object based on sensor type
func (b *Brick) createSensorForType(sensorPort models.SensorPort, sensorType models.SensorType) (SensorDataHandler, error) {
	var sensor SensorDataHandler
	var err error

	switch sensorType {
	case models.ButtonOrTouchSensor:
		sensor = sensors.NewButtonSensor(b, sensorPort)
	case models.SpikePrimeColorSensor, models.ColourAndDistanceSensor:
		sensor, err = sensors.NewColorSensor(b, sensorPort, sensorType)
		if err != nil {
			return nil, fmt.Errorf("failed to create color sensor: %w", err)
		}
	case models.SpikePrimeForceSensor:
		sensor = sensors.NewForceSensor(b, sensorPort)
	case models.SpikePrimeUltrasonicDistanceSensor:
		sensor, err = sensors.NewUltrasonicDistanceSensor(b, sensorPort)
		if err != nil {
			return nil, fmt.Errorf("failed to create ultrasonic distance sensor: %w", err)
		}
	case models.SimpleLights:
		sensor = sensors.NewPassiveLight(b, sensorPort)
	case models.SpikePrimeLargeMotor, models.SpikePrimeMediumMotor, models.TechnicMediumAngularMotor, models.TechnicLargeMotorID, models.TechnicXLMotorID, models.SpikeEssentialSmallAngularMotor, models.MediumLinearMotor:
		// Create active motor for motor types
		sensor, err = motors.NewActiveMotor(b, sensorPort, sensorType)
		if err != nil {
			return nil, fmt.Errorf("failed to create active motor: %w", err)
		}
	case models.SystemMediumMotor, models.SystemTrainMotor, models.SystemTurntableMotor, models.TechnicLargeMotor, models.TechnicXLMotor:
		// Create passive motor for passive motor types
		sensor, err = motors.NewPassiveMotor(b, sensorPort, sensorType)
		if err != nil {
			return nil, fmt.Errorf("failed to create passive motor: %w", err)
		}
	default:
		// For other sensor types, create a generic active sensor
		if sensorType.IsActiveSensor() {
			sensor = sensors.NewActiveSensor(b, sensorPort, sensorType)
		}
	}

	return sensor, nil
}

// Device connection/disconnection handlers
func (b *Brick) handleDeviceConnected(line string) {
	// Parse port and sensor type from line like "P0: connected to active ID 4B"
	portStr := string(line[1])
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 3 {
		return
	}

	// Extract sensor type (hex value at the end)
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	lastPart := parts[len(parts)-1]
	sensorTypeHex, err := strconv.ParseInt(lastPart, 16, 32)
	if err != nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.sensorTypes[port] = models.SensorType(sensorTypeHex)

	// Create appropriate sensor/motor object based on type
	sensorPort := models.SensorPort(port)
	sensorType := models.SensorType(sensorTypeHex)

	sensor, err := b.createSensorForType(sensorPort, sensorType)
	if err != nil {
		b.logger.Error("Failed to create sensor", "port", sensorPort, "error", err)
		return
	}

	if sensor != nil {
		b.sensors[port] = sensor
	}

	b.logger.Info("Device connected", "port", port, "sensor_type", sensorTypeHex)
}

func (b *Brick) handleDeviceDisconnected(line string) {
	portStr := string(line[1])
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 3 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.sensorTypes[port] = models.None
	b.sensors[port] = nil

	b.logger.Info("Device disconnected", "port", port)
}

func (b *Brick) handleSensorData(line string) {
	// Handle sensor data like "P0C0: +18 +5489 +12"
	b.logger.Debug("Sensor data received", "line", line)

	// Parse port from line (P0, P1, P2, P3)
	if len(line) < 3 || line[0] != 'P' {
		return
	}

	portStr := string(line[1])
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 3 {
		return
	}

	b.mu.RLock()
	sensor := b.sensors[port]
	b.mu.RUnlock()

	if sensor == nil {
		return
	}

	// Parse the data part (everything after ": ")
	dataStart := strings.Index(line, ": ")
	if dataStart == -1 {
		return
	}

	dataPart := line[dataStart+2:]
	dataValues := strings.Fields(dataPart)

	// Update sensor with parsed data
	if err := sensor.UpdateFromSensorData(dataValues); err != nil {
		b.logger.Error("Failed to update sensor data", "port", port, "error", err)
		return
	}

	// Set trigger flag to indicate data was received
	if triggerFlag := sensor.GetTriggerFlag(); triggerFlag != nil {
		*triggerFlag = true
	}
}

func (b *Brick) handleButtonReleased(line string) {
	portStr := string(line[1])
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 3 {
		return
	}

	b.mu.RLock()
	sensor := b.sensors[port]
	b.mu.RUnlock()

	if sensor != nil {
		// Update button sensor with released state
		if err := sensor.UpdateFromSensorData([]string{"0"}); err != nil {
			b.logger.Error("Failed to update sensor data", "error", err)
		}
		if triggerFlag := sensor.GetTriggerFlag(); triggerFlag != nil {
			*triggerFlag = true
		}
	}

	b.logger.Debug("Button released", "port", port)
}

func (b *Brick) handleButtonPressed(line string) {
	portStr := string(line[1])
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 3 {
		return
	}

	b.mu.RLock()
	sensor := b.sensors[port]
	b.mu.RUnlock()

	if sensor != nil {
		// Update button sensor with pressed state
		if err := sensor.UpdateFromSensorData([]string{"1"}); err != nil {
			b.logger.Error("Failed to update sensor data", "error", err)
		}
		if triggerFlag := sensor.GetTriggerFlag(); triggerFlag != nil {
			*triggerFlag = true
		}
	}

	b.logger.Debug("Button pressed", "port", port)
}

func (b *Brick) handlePowerFault() {
	b.logger.Error("Power fault detected")
	// Could emit an event here if needed
}

// Motor control methods

// SetPowerLevel sets the power level in percent for motors and lights
func (b *Brick) SetPowerLevel(port models.SensorPort, powerPercent int) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	// Check if the sensor can have its power level controlled
	if !sensorType.CanSetPowerLevel() {
		return fmt.Errorf("sensor type %s on port %s cannot have power level controlled", sensorType.String(), port.String())
	}

	// Clamp power between -100 and 100
	if powerPercent < -100 {
		powerPercent = -100
	} else if powerPercent > 100 {
		powerPercent = 100
	}

	if sensorType.IsActiveSensor() {
		// Set continuous reading for active motors
		if err := b.SelectCombiModesAndRead(port, []int{1, 2, 3}, false); err != nil {
			return fmt.Errorf("failed to select combi modes: %w", err)
		}
		return b.writeCommand(fmt.Sprintf("port %d ; pid %d 0 0 s1 1 0 0.003 0.01 0 100; set %d\r", port.Byte(), port.Byte(), powerPercent))
	}
	// Passive motor or light
	powerFloat := float64(powerPercent) / 100.0
	return b.writeCommand(fmt.Sprintf("port %d ; pwm ; set %.3f\r", port.Byte(), powerFloat))
}

// SetMotorLimits sets the motor speed limit
func (b *Brick) SetMotorLimits(port models.SensorPort, powerLimit float64) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected to port %s", port.String())
	}

	// Clamp power limit between 0 and 1
	if powerLimit < 0 {
		powerLimit = 0
	} else if powerLimit > 1 {
		powerLimit = 1
	}

	return b.writeCommand(fmt.Sprintf("port %d ; plimit %.3f\r", port.Byte(), powerLimit))
}

// SetMotorBias sets the motor bias
func (b *Brick) SetMotorBias(port models.SensorPort, bias float64) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected to port %s", port.String())
	}

	// Clamp bias between 0 and 1
	if bias < 0 {
		bias = 0
	} else if bias > 1 {
		bias = 1
	}

	return b.writeCommand(fmt.Sprintf("port %d ; bias %.3f\r", port.Byte(), bias))
}

// MoveMotorForSeconds runs the specified motor for an amount of seconds
func (b *Brick) MoveMotorForSeconds(ctx context.Context, port models.SensorPort, seconds float64, speed int, blocking bool) error {
	if seconds <= 0 {
		return nil // No need to move
	}

	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected to port %s", port.String())
	}

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("not an active motor connected to port %s", port.String())
	}

	if speed == 0 {
		return fmt.Errorf("speed can't be 0")
	}

	// Clamp speed between -100 and 100
	if speed < -100 {
		speed = -100
	} else if speed > 100 {
		speed = 100
	}

	// Set continuous reading
	if err := b.SelectCombiModesAndRead(port, []int{1, 2, 3}, false); err != nil {
		return fmt.Errorf("failed to select combi modes: %w", err)
	}

	command := fmt.Sprintf("port %d ; pid %d 0 0 s1 1 0 0.003 0.01 0 100; set pulse %d 0.0 %.3f 0\r", port.Byte(), port.Byte(), speed, seconds)
	if err := b.writeCommand(command); err != nil {
		return fmt.Errorf("failed to write move command: %w", err)
	}

	if blocking {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(seconds * float64(time.Second))):
			return nil
		}
	}

	return nil
}

// FloatMotor floats the motor and stops all constraints on it
func (b *Brick) FloatMotor(port models.SensorPort) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected to port %s", port.String())
	}

	return b.writeCommand("coast\r")
}

// Sensor control methods

// SelectModeAndRead selects modes on a specific port (only possible on active sensors and motors)
func (b *Brick) SelectModeAndRead(port models.SensorPort, mode int, readOnce bool) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("mode can be changed only on active sensors")
	}

	// TODO: Validate mode exists for the sensor
	// This would require implementing sensor-specific mode validation

	var command string
	if readOnce {
		command = fmt.Sprintf("port %d ; combi %d ; selonce %d\r", port.Byte(), mode, mode)
	} else {
		command = fmt.Sprintf("port %d ; combi %d ; select %d\r", port.Byte(), mode, mode)
	}

	return b.writeCommand(command)
}

// SelectCombiModesAndRead selects multiple modes on a specific port
func (b *Brick) SelectCombiModesAndRead(port models.SensorPort, modes []int, readOnce bool) error {
	if len(modes) == 0 {
		return fmt.Errorf("modes can't be empty")
	}

	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("mode can be changed only on active sensors")
	}

	// TODO: Validate all modes exist for the sensor

	var command strings.Builder
	command.WriteString(fmt.Sprintf("port %d ; combi 0 ", port.Byte()))

	for _, mode := range modes {
		command.WriteString(fmt.Sprintf("%d 0 ", mode))
	}

	if readOnce {
		command.WriteString("; selonce 0\r")
	} else {
		command.WriteString("; select 0\r")
	}

	return b.writeCommand(command.String())
}

// StopContinuousReadingSensor stops reading continuous data from a specific sensor
func (b *Brick) StopContinuousReadingSensor(port models.SensorPort) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("mode can be changed only on active sensors")
	}

	return b.writeCommand(fmt.Sprintf("port %d ; select\r", port.Byte()))
}

// SwitchSensorOn switches a sensor on
func (b *Brick) SwitchSensorOn(port models.SensorPort) error {
	return b.writeCommand(fmt.Sprintf("port %d ; plimit 1 ; on\r", port.Byte()))
}

// SwitchSensorOff switches a sensor off
func (b *Brick) SwitchSensorOff(port models.SensorPort) error {
	return b.writeCommand(fmt.Sprintf("port %d ; plimit 1 ; off\r", port.Byte()))
}

// WriteBytesToSensor writes bytes directly to a sensor
func (b *Brick) WriteBytesToSensor(port models.SensorPort, data []byte, singleHeader bool) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to write")
	}

	// Format the command according to C# implementation
	var command strings.Builder
	command.WriteString(fmt.Sprintf("port %d ; ", port.Byte()))

	if singleHeader {
		command.WriteString("write1 ")
	} else {
		command.WriteString("write2 ")
	}

	// Convert bytes to hex format with spaces
	for i, byteVal := range data {
		if i > 0 {
			command.WriteString(" ")
		}
		command.WriteString(fmt.Sprintf("%02X", byteVal))
	}

	command.WriteString("\r")

	return b.writeCommand(command.String())
}

// SendRawCommand sends a raw command to the BuildHat
func (b *Brick) SendRawCommand(command string) error {
	if command == "" {
		return nil
	}

	if !strings.HasSuffix(command, "\r") {
		command += "\r"
	}

	return b.writeCommand(command)
}

// ClearFaults clears any fault
func (b *Brick) ClearFaults() error {
	return b.writeCommand("clear_faults\r")
}

// Additional methods needed for BrickInterface

// MoveMotorToPosition runs the motor to a relative position
func (b *Brick) MoveMotorToPosition(ctx context.Context, port models.SensorPort, targetPosition, speed int, blocking bool) error {
	const angleTolerance = 2.028

	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected")
	}

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("not an active motor connected")
	}

	if speed == 0 {
		return fmt.Errorf("speed can't be 0")
	}

	// Clamp speed to valid range
	if speed < -100 {
		speed = -100
	} else if speed > 100 {
		speed = 100
	}

	b.mu.RLock()
	motor, ok := b.sensors[port.Byte()].(*motors.ActiveMotor)
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("motor not found")
	}

	// Get current relative position
	actualPosition := motor.GetPosition()
	actualPositionDouble := float64(actualPosition) / 360.0
	newPosition := (float64(targetPosition) - actualPositionDouble) / 360.0

	// Calculate duration based on speed and power limit
	duration := math.Abs(newPosition-actualPositionDouble) / (float64(speed) * 0.05 * motor.GetPowerLimit())

	// Set continuous reading
	if err := b.SelectCombiModesAndRead(port, []int{1, 2, 3}, false); err != nil {
		return fmt.Errorf("failed to set continuous reading: %w", err)
	}

	// Send ramp command
	command := fmt.Sprintf("port %d ; pid %d 0 1 s4 0.0027777778 0 5 0 .1 3 ; set ramp %.10f %.10f %.10f 0\r",
		port.Byte(), port.Byte(), actualPositionDouble, newPosition, duration)

	if err := b.writeCommand(command); err != nil {
		return fmt.Errorf("failed to send ramp command: %w", err)
	}

	if blocking {
		pos := motor.GetPosition()
		for !((float64(pos)/360.0 < newPosition+angleTolerance) && (float64(pos)/360.0 > newPosition-angleTolerance)) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				time.Sleep(5 * time.Millisecond)
				pos = motor.GetPosition()
			}
		}
	}

	return nil
}

// MoveMotorToAbsolutePosition runs the motor to an absolute position
func (b *Brick) MoveMotorToAbsolutePosition(ctx context.Context, port models.SensorPort, targetPosition int, way models.PositionWay, speed int, blocking bool) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected")
	}

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("not an active motor connected")
	}

	if targetPosition < -180 || targetPosition > 179 {
		return fmt.Errorf("target position can only be between -180 and +179")
	}

	if speed == 0 {
		return fmt.Errorf("speed can't be 0")
	}

	// Clamp speed between -100 and 100
	if speed < -100 {
		speed = -100
	} else if speed > 100 {
		speed = 100
	}

	// Get motor from sensors
	b.mu.RLock()
	sensorHandler := b.sensors[port.Byte()]
	b.mu.RUnlock()

	if sensorHandler == nil {
		return fmt.Errorf("no sensor found on port %s", port.String())
	}

	// Cast to ActiveMotor to get position data
	activeMotor, ok := sensorHandler.(*motors.ActiveMotor)
	if !ok {
		return fmt.Errorf("sensor on port %s is not an active motor", port.String())
	}

	// Get current positions
	actualPosition := activeMotor.GetPosition()
	actualAbsolutePosition := activeMotor.GetAbsolutePosition()

	// Calculate new position using ToAbsolutePosition logic
	newPosition := b.toAbsolutePosition(targetPosition, way, actualPosition, actualAbsolutePosition)

	// Calculate duration based on speed and power limit
	powerLimit := activeMotor.GetPowerLimit()
	actualPositionDouble := float64(actualPosition) / 360.0
	duration := math.Abs(newPosition-actualPositionDouble) / (float64(speed) * 0.05 * powerLimit)

	// Set continuous reading
	if err := b.SelectCombiModesAndRead(port, []int{1, 2, 3}, false); err != nil {
		return fmt.Errorf("failed to select combi modes: %w", err)
	}

	// Send ramp command
	command := fmt.Sprintf("port %d ; pid %d 0 1 s4 0.0027777778 0 5 0 .1 3 ; set ramp %.10f %.10f %.10f 0\r",
		port.Byte(), port.Byte(), actualPositionDouble, newPosition, duration)

	if err := b.writeCommand(command); err != nil {
		return fmt.Errorf("failed to send ramp command: %w", err)
	}

	if blocking {
		// Wait for the duration
		select {
		case <-time.After(time.Duration(duration * float64(time.Second))):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// validateMotorForDegrees validates the motor and speed for degree-based movement
func (b *Brick) validateMotorForDegrees(port models.SensorPort, speed int) error {
	b.mu.RLock()
	sensorType := b.sensorTypes[port.Byte()]
	b.mu.RUnlock()

	if !sensorType.IsMotor() {
		return fmt.Errorf("not a motor connected")
	}

	if !sensorType.IsActiveSensor() {
		return fmt.Errorf("not an active motor connected")
	}

	if speed == 0 {
		return fmt.Errorf("speed can't be 0")
	}

	return nil
}

// getActiveMotorFromPort retrieves and validates an active motor from the specified port
func (b *Brick) getActiveMotorFromPort(port models.SensorPort) (*motors.ActiveMotor, error) {
	b.mu.RLock()
	sensorHandler := b.sensors[port.Byte()]
	b.mu.RUnlock()

	if sensorHandler == nil {
		return nil, fmt.Errorf("no sensor found on port %s", port.String())
	}

	activeMotor, ok := sensorHandler.(*motors.ActiveMotor)
	if !ok {
		return nil, fmt.Errorf("sensor on port %s is not an active motor", port.String())
	}

	return activeMotor, nil
}

// clampSpeed clamps speed to valid range [-100, 100]
func clampSpeed(speed int) int {
	if speed < -100 {
		return -100
	}
	if speed > 100 {
		return 100
	}
	return speed
}

// MoveMotorForDegrees runs the motor for a specific number of degrees
func (b *Brick) MoveMotorForDegrees(ctx context.Context, port models.SensorPort, targetPosition, speed int, blocking bool) error {
	const angleTolerance = 2.028

	// Validate motor and speed
	if err := b.validateMotorForDegrees(port, speed); err != nil {
		return err
	}

	// Clamp speed to valid range
	speed = clampSpeed(speed)

	// Get active motor
	activeMotor, err := b.getActiveMotorFromPort(port)
	if err != nil {
		return err
	}

	// Calculate movement parameters
	actualPosition := activeMotor.GetPosition()
	actualPositionDouble := float64(actualPosition) / 360.0

	// Adjust target position based on speed direction
	if speed < 0 {
		targetPosition = -targetPosition
	}

	newPosition := float64(targetPosition) - actualPositionDouble/360.0

	// Calculate duration based on speed and power limit
	powerLimit := activeMotor.GetPowerLimit()
	duration := math.Abs(newPosition-actualPositionDouble) / (float64(speed) * 0.05 / powerLimit)

	// Set continuous reading
	if err := b.SelectCombiModesAndRead(port, []int{1, 2, 3}, false); err != nil {
		return fmt.Errorf("failed to select combi modes: %w", err)
	}

	// Send ramp command
	command := fmt.Sprintf("port %d ; pid %d 0 1 s4 0.0027777778 0 5 0 .1 3 ; set ramp %.10f %.10f %.10f 0\r",
		port.Byte(), port.Byte(), actualPositionDouble, newPosition, duration)

	if err := b.writeCommand(command); err != nil {
		return fmt.Errorf("failed to send ramp command: %w", err)
	}

	if blocking {
		// Wait for the motor to reach the target position with tolerance
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentPos := activeMotor.GetPosition()
				currentPosDouble := float64(currentPos) / 360.0

				if currentPosDouble >= newPosition-angleTolerance && currentPosDouble <= newPosition+angleTolerance {
					return nil
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

// toAbsolutePosition calculates the new position based on target position, way, and current positions
func (b *Brick) toAbsolutePosition(targetPosition int, way models.PositionWay, actualPosition, actualAbsolutePosition int) float64 {
	difference := (targetPosition-actualAbsolutePosition+180)%360 - 180
	var newPosition float64

	switch way {
	case models.Clockwise:
		if difference > 0 {
			newPosition = (float64(actualPosition) + float64(difference)) / 360.0
		} else {
			newPosition = (float64(actualPosition) + 360 + float64(difference)) / 360.0
		}
	case models.AntiClockwise:
		if difference < 0 {
			newPosition = -(float64(actualPosition) - float64(difference)) / 360.0
		} else {
			newPosition = -(float64(actualPosition) + 360 - float64(difference)) / 360.0
		}
	case models.Shortest:
	default:
		if math.Abs(float64(difference)) > 180 {
			newPosition = (float64(actualPosition) + 360 + float64(difference)) / 360.0
		} else {
			newPosition = (float64(actualPosition) + float64(difference)) / 360.0
		}
	}

	return newPosition
}

// LoadFirmwareFile loads the firmware.bin file from embedded resources
func (b *Brick) LoadFirmwareFile() ([]byte, error) {
	data, err := embeddedData.ReadFile("data/firmware.bin")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded firmware.bin: %w", err)
	}
	b.logger.Debug("Loaded firmware", "size", len(data))
	return data, nil
}

// LoadSignatureFile loads the signature.bin file from embedded resources
func (b *Brick) LoadSignatureFile() ([]byte, error) {
	data, err := embeddedData.ReadFile("data/signature.bin")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded signature.bin: %w", err)
	}
	b.logger.Debug("Loaded signature", "size", len(data))
	return data, nil
}

// GetEmbeddedVersion reads the version file from embedded resources
func (b *Brick) GetEmbeddedVersion() (string, error) {
	data, err := embeddedData.ReadFile("data/version")
	if err != nil {
		return "", fmt.Errorf("failed to read embedded version file: %w", err)
	}
	version := strings.TrimSpace(string(data))
	b.logger.Debug("Read embedded version", "version", version)
	return version, nil
}

// GetHardwareVersion reads the firmware version from the BuildHat hardware
func (b *Brick) GetHardwareVersion() (string, error) {
	// Send version command
	if err := b.writeCommand("version\r"); err != nil {
		return "", fmt.Errorf("failed to send version command: %w", err)
	}

	// Wait for response
	time.Sleep(100 * time.Millisecond)

	// Read response
	response := b.readExisting()
	b.logger.Debug("Hardware version response", "response", response)

	// Parse version from response like "Firmware version: 1636109636 2021-11-05T10:53:56+00:00"
	if strings.Contains(response, "Firmware version: ") {
		parts := strings.Fields(response)
		for i, part := range parts {
			if part == "version:" && i+1 < len(parts) {
				return parts[i+1], nil
			}
		}
	}

	return "", fmt.Errorf("could not parse version from response: %s", response)
}

// ReadInputVoltage reads the input voltage from the BuildHat hardware
func (b *Brick) ReadInputVoltage() (float64, error) {
	rawVoltage, err := b.getRawVoltage()
	if err != nil {
		return 0, fmt.Errorf("failed to get raw voltage: %w", err)
	}
	if rawVoltage == "" {
		return 0, fmt.Errorf("no voltage data received")
	}

	voltage, err := strconv.ParseFloat(rawVoltage, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse voltage '%s': %w", rawVoltage, err)
	}

	return voltage, nil
}

// GetFirmwareChecksum calculates the checksum for the firmware
func (b *Brick) GetFirmwareChecksum(firmware []byte) uint32 {
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

// writeBinaryData writes binary data with STX/ETX markers
func (b *Brick) writeBinaryData(data []byte) error {
	// Write STX (0x02)
	if _, err := b.writer.Write([]byte{0x02}); err != nil {
		return err
	}

	// Write the data
	if _, err := b.writer.Write(data); err != nil {
		return err
	}

	// Write ETX (0x03) and carriage return
	if _, err := b.writer.Write([]byte{0x03, '\r'}); err != nil {
		return err
	}

	return nil
}

func GetDevice[T any](b *Brick, port models.SensorPort) (T, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handler := b.sensors[port.Byte()]
	if handler == nil {
		var zero T
		return zero, fmt.Errorf("no device connected to port %s", port.String())
	}

	if device, ok := handler.(T); ok {
		return device, nil
	}

	var zero T
	return zero, fmt.Errorf("device on port %s is not of requested type", port.String())
}
