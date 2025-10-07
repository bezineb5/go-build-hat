package motors

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
	"github.com/bezineb5/go-build-hat/pkg/buildhat/sensors"
)

// PassiveMotor represents a passive motor without position feedback
type PassiveMotor struct {
	*sensors.BaseSensor

	// Internal state
	isRunning bool
	speed     int

	// Thread safety
	mu sync.RWMutex
}

// NewPassiveMotor creates a new passive motor
func NewPassiveMotor(brick sensors.BrickInterface, port models.SensorPort, motorType models.SensorType) *PassiveMotor {
	motor := &PassiveMotor{
		BaseSensor: sensors.NewBaseSensor(brick, port, motorType),
	}

	// Initialize motor with default bias and power limit
	brick.SetMotorBias(port, 0.0)   // Default bias
	brick.SetMotorLimits(port, 1.0) // Default power limit (full power)

	return motor
}

// GetMotorName gets the name of the motor based on its type
func (m *PassiveMotor) GetMotorName() string {
	switch m.GetSensorType() {
	case models.SystemTrainMotor:
		return "System train motor"
	case models.SystemTurntableMotor:
		return "System turntable motor"
	case models.SystemMediumMotor:
		return "System medium Motor"
	case models.TechnicLargeMotor:
		return "Technic large motor"
	case models.TechnicXLMotor:
		return "Technic XL motor"
	default:
		return "Unknown passive motor"
	}
}

// GetSensorName gets the sensor name (implements sensors.Sensor interface)
func (m *PassiveMotor) GetSensorName() string {
	return m.GetMotorName()
}

// GetSpeed gets the current speed
func (m *PassiveMotor) GetSpeed() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.speed
}

// Speed gets the current speed (alias for GetSpeed)
func (m *PassiveMotor) Speed() int {
	return m.GetSpeed()
}

// SetSpeed sets the speed of the motor
func (m *PassiveMotor) SetSpeed(speed int) error {
	if speed < -100 || speed > 100 {
		return fmt.Errorf("speed must be between -100 and 100, got %d", speed)
	}

	m.mu.Lock()
	m.speed = speed
	isRunning := m.isRunning
	m.mu.Unlock()

	// If motor is running, apply the new speed immediately
	if isRunning {
		return m.Start()
	}

	return nil
}

// Start starts the motor with the current speed
func (m *PassiveMotor) Start() error {
	m.mu.RLock()
	speed := m.speed
	m.mu.RUnlock()

	if err := m.GetBrick().SetPowerLevel(m.GetPort(), speed); err != nil {
		return err
	}

	m.mu.Lock()
	m.isRunning = true
	m.mu.Unlock()

	return nil
}

// StartWithSpeed starts the motor with the specified speed
func (m *PassiveMotor) StartWithSpeed(speed int) error {
	if err := m.SetSpeed(speed); err != nil {
		return err
	}
	return m.Start()
}

// Stop stops the motor
func (m *PassiveMotor) Stop() error {
	if err := m.SetSpeed(0); err != nil {
		return err
	}

	m.mu.Lock()
	m.isRunning = false
	m.mu.Unlock()

	return nil
}

// SetBias sets the bias of the motor
func (m *PassiveMotor) SetBias(bias float64) error {
	if bias < 0 || bias > 1 {
		return fmt.Errorf("bias must be between 0 and 1, got %f", bias)
	}

	return m.GetBrick().SetMotorBias(m.GetPort(), bias)
}

// SetPowerLimit sets the power consumption limit
func (m *PassiveMotor) SetPowerLimit(plimit float64) error {
	if plimit < 0 || plimit > 1 {
		return fmt.Errorf("power limit must be between 0 and 1, got %f", plimit)
	}

	return m.GetBrick().SetMotorLimits(m.GetPort(), plimit)
}

// Float floats the motor and stops all constraints on it
func (m *PassiveMotor) Float() error {
	return m.GetBrick().FloatMotor(m.GetPort())
}

// IsRunning gets whether the motor is currently running
func (m *PassiveMotor) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isRunning
}

// GetBrick gets the brick interface (helper method)
func (m *PassiveMotor) GetBrick() sensors.BrickInterface {
	return m.BaseSensor.GetBrick()
}

// UpdateFromSensorData updates motor values from raw sensor data
func (m *PassiveMotor) UpdateFromSensorData(data []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(data) >= 1 {
		// Parse speed from data
		if speed, err := strconv.Atoi(data[0]); err == nil {
			m.speed = speed
		}
	}

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this motor
func (m *PassiveMotor) GetTriggerFlag() *bool {
	// PassiveMotor doesn't have a trigger flag, return nil
	return nil
}
