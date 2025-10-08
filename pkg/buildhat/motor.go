package buildhat

import (
	"fmt"
	"time"
)

// MotorRunMode represents the current mode the motor is in
type MotorRunMode int

const (
	MotorRunModeNone MotorRunMode = iota
	MotorRunModeFree
	MotorRunModeDegrees
	MotorRunModeSeconds
)

// MotorDirection represents the direction for position-based movements
type MotorDirection int

const (
	DirectionShortest MotorDirection = iota
	DirectionClockwise
	DirectionAnticlockwise
)

// Motor creates a motor interface for the specified port
func (b *Brick) Motor(port BuildHatPort) *Motor {
	motor := &Motor{
		brick:        b,
		port:         port.Int(),
		defaultSpeed: 20,
		currentSpeed: 0,
		runMode:      MotorRunModeNone,
		release:      true,
		rpm:          false,
	}

	// Initialize motor with default settings
	_ = motor.SetPowerLimit(0.7)
	_ = motor.SetPWMParams(0.65, 0.01)

	return motor
}

// Motor provides a Python-like motor interface
type Motor struct {
	brick        *Brick
	port         int
	defaultSpeed int
	currentSpeed int
	runMode      MotorRunMode
	release      bool
	rpm          bool
}

// SetDefaultSpeed sets the default speed of the motor (-100 to 100)
func (m *Motor) SetDefaultSpeed(speed int) error {
	if speed < -100 || speed > 100 {
		return fmt.Errorf("invalid speed: must be between -100 and 100")
	}
	m.defaultSpeed = speed
	return nil
}

// SetSpeedUnitRPM sets whether to use RPM for speed units or not
func (m *Motor) SetSpeedUnitRPM(rpm bool) {
	m.rpm = rpm
}

// RunForRotations runs the motor for N rotations
func (m *Motor) RunForRotations(rotations float64, speed int) error {
	return m.RunForDegrees(int(rotations*360), speed)
}

// RunForDegrees runs the motor for the specified number of degrees
func (m *Motor) RunForDegrees(degrees, speed int) error {
	if speed == 0 {
		speed = m.defaultSpeed
	}
	if speed < -100 || speed > 100 {
		return fmt.Errorf("invalid speed: must be between -100 and 100")
	}

	m.runMode = MotorRunModeDegrees

	// Get current position
	position, err := m.GetPosition()
	if err != nil {
		// If we can't get position, use a simple approach
		position = 0
	}

	// Calculate target position
	mul := 1
	actualSpeed := speed
	if speed < 0 {
		actualSpeed = -speed
		mul = -1
	}

	newPos := float64(position+degrees*mul) / 360.0
	currentPos := float64(position) / 360.0

	// Process speed
	processedSpeed := float64(actualSpeed) * 0.05 // Collapse speed range to 0-5

	// Calculate duration
	var duration time.Duration
	if processedSpeed != 0 {
		durationSecs := (newPos - currentPos) / processedSpeed
		if durationSecs < 0 {
			durationSecs = -durationSecs
		}
		duration = time.Duration(durationSecs * float64(time.Second))
	}

	// Send ramp command
	durationSecs := duration.Seconds()
	cmd := fmt.Sprintf("port %d ; select 0 ; pid %d 0 1 s4 0.0027777778 0 5 0 .1 3 0.01 ; set ramp %.6f %.6f %.6f 0",
		m.port, m.port, currentPos, newPos, durationSecs)
	if err := m.brick.writeCommand(cmd); err != nil {
		return err
	}

	// Wait for completion (simplified - in real implementation would use futures)
	if duration > 0 {
		time.Sleep(duration)
	}

	// Coast to stop if release is enabled
	if m.release {
		time.Sleep(200 * time.Millisecond)
		_ = m.Coast()
	}

	m.runMode = MotorRunModeNone
	return nil
}

// RunForDuration runs the motor for the specified duration
func (m *Motor) RunForDuration(duration time.Duration, speed int) error {
	if speed == 0 {
		speed = m.defaultSpeed
	}
	if speed < -100 || speed > 100 {
		return fmt.Errorf("invalid speed: must be between -100 and 100")
	}

	m.runMode = MotorRunModeSeconds

	// Process speed
	processedSpeed := float64(speed) * 0.05 // Collapse to -5 to 5 range

	// Set up PID for speed control
	pidCmd := "pid %d 0 0 s1 1 0 0.003 0.01 0 100 0.01"
	if m.rpm {
		pidCmd = "pid_diff %d 0 5 s2 0.0027777778 1 0 2.5 0 .4 0.01"
	}

	seconds := duration.Seconds()
	cmd := fmt.Sprintf("port %d ; select 0 ; %s ; set pulse %.6f 0.0 %.6f 0",
		m.port, fmt.Sprintf(pidCmd, m.port), processedSpeed, seconds)
	if err := m.brick.writeCommand(cmd); err != nil {
		return err
	}

	// Wait for the specified time
	time.Sleep(duration)

	// Coast to stop if release is enabled
	if m.release {
		_ = m.Coast()
	}

	m.runMode = MotorRunModeNone
	return nil
}

// RunToPosition runs motor to a specific position (in degrees, -180 to 180)
func (m *Motor) RunToPosition(degrees, speed int, direction MotorDirection) error {
	if err := m.validatePositionParams(degrees, speed, direction); err != nil {
		return err
	}

	m.runMode = MotorRunModeDegrees

	pos, apos, err := m.getCurrentAndAbsolutePosition()
	if err != nil {
		return err
	}

	newPos := m.calculateTargetPosition(pos, apos, degrees, direction)
	duration := m.calculateMovementDuration(float64(pos), newPos, speed)

	if err := m.executeRampMovement(float64(pos), newPos, duration); err != nil {
		return err
	}

	m.waitForMovementCompletion(duration)
	m.runMode = MotorRunModeNone
	return nil
}

// validatePositionParams validates parameters for RunToPosition
func (m *Motor) validatePositionParams(degrees, speed int, direction MotorDirection) error {
	if speed < 0 || speed > 100 {
		return fmt.Errorf("invalid speed: must be between 0 and 100")
	}
	if degrees < -180 || degrees > 180 {
		return fmt.Errorf("invalid angle: must be between -180 and 180")
	}
	if direction != DirectionShortest && direction != DirectionClockwise && direction != DirectionAnticlockwise {
		return fmt.Errorf("invalid direction: must be DirectionShortest, DirectionClockwise, or DirectionAnticlockwise")
	}
	return nil
}

// getCurrentAndAbsolutePosition retrieves current motor position data
func (m *Motor) getCurrentAndAbsolutePosition() (pos, apos int, err error) {
	data, err := m.getData()
	if err != nil {
		return 0, 0, err
	}
	if len(data) < 3 {
		return 0, 0, fmt.Errorf("insufficient motor data")
	}
	return data[1].(int), data[2].(int), nil
}

// calculateTargetPosition calculates the target position based on direction
func (m *Motor) calculateTargetPosition(pos, apos, degrees int, direction MotorDirection) float64 {
	diff := (degrees-apos+180)%360 - 180

	if direction == DirectionShortest {
		return float64(pos+diff) / 360.0
	}

	path1, path2 := m.calculatePaths(degrees, apos, diff)
	return m.selectPathByDirection(pos, path1, path2, direction)
}

// calculatePaths calculates alternate paths for position movement
func (m *Motor) calculatePaths(degrees, apos, diff int) (path1, path2 int) {
	v1 := (degrees - apos) % 360
	v2 := (apos - degrees) % 360

	mul := 1
	if diff > 0 {
		mul = -1
	}

	path1 = diff
	path2 = mul * v2
	if diff == v1 {
		path2 = mul * v1
	}
	return path1, path2
}

// selectPathByDirection selects the appropriate path based on direction
func (m *Motor) selectPathByDirection(pos, path1, path2 int, direction MotorDirection) float64 {
	switch direction {
	case DirectionClockwise:
		if path2 > path1 {
			return float64(pos+path2) / 360.0
		}
		return float64(pos+path1) / 360.0
	case DirectionAnticlockwise:
		if path1 < path2 {
			return float64(pos+path1) / 360.0
		}
		return float64(pos+path2) / 360.0
	default:
		return float64(pos+path1) / 360.0
	}
}

// calculateMovementDuration calculates how long the movement will take
func (m *Motor) calculateMovementDuration(currentPos, newPos float64, speed int) time.Duration {
	processedSpeed := float64(speed) * 0.05
	if processedSpeed == 0 {
		return 0
	}

	durationSecs := (newPos - currentPos) / processedSpeed
	if durationSecs < 0 {
		durationSecs = -durationSecs
	}
	return time.Duration(durationSecs * float64(time.Second))
}

// executeRampMovement sends the ramp command to the motor
func (m *Motor) executeRampMovement(currentPos, newPos float64, duration time.Duration) error {
	durationSecs := duration.Seconds()
	cmd := fmt.Sprintf("port %d ; select 0 ; pid %d 0 1 s4 0.0027777778 0 5 0 .1 3 0.01 ; set ramp %.6f %.6f %.6f 0",
		m.port, m.port, currentPos/360.0, newPos, durationSecs)
	return m.brick.writeCommand(cmd)
}

// waitForMovementCompletion waits for movement to complete and optionally coasts
func (m *Motor) waitForMovementCompletion(duration time.Duration) {
	if duration > 0 {
		time.Sleep(duration)
	}

	if m.release {
		time.Sleep(200 * time.Millisecond)
		_ = m.Coast()
	}
}

// Start starts the motor at the specified speed
func (m *Motor) Start(speed int) error {
	if speed == 0 {
		speed = m.defaultSpeed
	}
	if speed < -100 || speed > 100 {
		return fmt.Errorf("invalid speed: must be between -100 and 100")
	}

	// If already running at this speed, do nothing
	if m.runMode == MotorRunModeFree && m.currentSpeed == speed {
		return nil
	}

	// If motor is running in another mode, don't interrupt
	if m.runMode != MotorRunModeNone && m.runMode != MotorRunModeFree {
		return fmt.Errorf("motor is busy in another mode")
	}

	// Process speed
	processedSpeed := float64(speed) * 0.05

	// Set up PID
	pidCmd := "pid %d 0 0 s1 1 0 0.003 0.01 0 100 0.01"
	if m.rpm {
		pidCmd = "pid_diff %d 0 5 s2 0.0027777778 1 0 2.5 0 .4 0.01"
	}

	cmd := fmt.Sprintf("port %d ; select 0 ; %s ; set %.6f",
		m.port, fmt.Sprintf(pidCmd, m.port), processedSpeed)

	if err := m.brick.writeCommand(cmd); err != nil {
		return err
	}

	m.runMode = MotorRunModeFree
	m.currentSpeed = speed
	return nil
}

// Stop stops the motor
func (m *Motor) Stop() error {
	m.runMode = MotorRunModeNone
	m.currentSpeed = 0
	return m.Coast()
}

// Coast puts the motor into coast mode (freely spinning)
func (m *Motor) Coast() error {
	cmd := fmt.Sprintf("port %d ; coast", m.port)
	return m.brick.writeCommand(cmd)
}

// Float puts the motor into float mode (same as coast)
func (m *Motor) Float() error {
	return m.Coast()
}

// GetPosition gets the position of motor relative to preset position
func (m *Motor) GetPosition() (int, error) {
	data, err := m.getData()
	if err != nil {
		return 0, err
	}
	if len(data) < 2 {
		return 0, fmt.Errorf("insufficient motor data")
	}
	if pos, ok := data[1].(int); ok {
		return pos, nil
	}
	return 0, fmt.Errorf("invalid position data type")
}

// GetAbsolutePosition gets the absolute position of motor (-180 to 180)
func (m *Motor) GetAbsolutePosition() (int, error) {
	data, err := m.getData()
	if err != nil {
		return 0, err
	}
	if len(data) < 3 {
		return 0, fmt.Errorf("no absolute position available for this motor")
	}
	if apos, ok := data[2].(int); ok {
		return apos, nil
	}
	return 0, fmt.Errorf("invalid absolute position data type")
}

// GetSpeed gets the current speed of the motor
func (m *Motor) GetSpeed() (int, error) {
	data, err := m.getData()
	if err != nil {
		return 0, err
	}
	if len(data) < 1 {
		return 0, fmt.Errorf("insufficient motor data")
	}
	if speed, ok := data[0].(int); ok {
		return speed, nil
	}
	return 0, fmt.Errorf("invalid speed data type")
}

// SetPowerLimit limits the power to the motor (0.0 to 1.0)
func (m *Motor) SetPowerLimit(limit float64) error {
	if limit < 0 || limit > 1 {
		return fmt.Errorf("power limit must be between 0 and 1")
	}
	cmd := fmt.Sprintf("port %d ; port_plimit %.2f", m.port, limit)
	return m.brick.writeCommand(cmd)
}

// SetPWMParams sets PWM thresholds
func (m *Motor) SetPWMParams(pwmThresh, minPWM float64) error {
	if pwmThresh < 0 || pwmThresh > 1 {
		return fmt.Errorf("pwmThresh must be between 0 and 1")
	}
	if minPWM < 0 || minPWM > 1 {
		return fmt.Errorf("minPWM must be between 0 and 1")
	}
	cmd := fmt.Sprintf("port %d ; pwmparams %.2f %.2f", m.port, pwmThresh, minPWM)
	return m.brick.writeCommand(cmd)
}

// PWM sets the motor to PWM mode with the specified value (-1.0 to 1.0)
func (m *Motor) PWM(value float64) error {
	if value < -1 || value > 1 {
		return fmt.Errorf("PWM value must be between -1 and 1")
	}
	cmd := fmt.Sprintf("port %d ; pwm ; set %.2f", m.port, value)
	return m.brick.writeCommand(cmd)
}

// PresetPosition presets the motor position to 0
func (m *Motor) PresetPosition() error {
	cmd := fmt.Sprintf("port %d ; preset", m.port)
	return m.brick.writeCommand(cmd)
}

// SetRelease sets whether the motor should coast after completing a movement
func (m *Motor) SetRelease(release bool) {
	m.release = release
}

// getData gets the current motor data (speed, position, absolute position)
func (m *Motor) getData() ([]interface{}, error) {
	// Set motor to combi mode to read data
	cmd := fmt.Sprintf("port %d ; select 0", m.port)
	if err := m.brick.writeCommand(cmd); err != nil {
		return nil, err
	}

	// Wait for sensor data
	data, err := m.brick.getSensorData(m.port)
	if err != nil {
		return nil, err
	}

	return data, nil
}
