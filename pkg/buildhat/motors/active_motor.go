package motors

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
	"github.com/bezineb5/go-build-hat/pkg/buildhat/sensors"
)

// ActiveMotor represents an active motor with position feedback
type ActiveMotor struct {
	*sensors.BaseSensor

	// Internal state
	tacho         int
	absoluteTacho int
	speed         int
	targetSpeed   int
	powerLimit    float64

	// Thread safety
	mu sync.RWMutex
}

// NewActiveMotor creates a new active motor
func NewActiveMotor(brick sensors.BrickInterface, port models.SensorPort, motorType models.SensorType) (*ActiveMotor, error) {
	motor := &ActiveMotor{
		BaseSensor: sensors.NewBaseSensor(brick, port, motorType),
		powerLimit: 0.7, // Default power limit
	}

	// Initialize motor with default bias and power limit
	if err := brick.SetMotorBias(port, 0.3); err != nil {
		return nil, err
	}
	if err := brick.SetMotorLimits(port, 0.7); err != nil {
		return nil, err
	}

	return motor, nil
}

const (
	technicXLMotorName       = "Technic XL motor"
	technicLargeMotorName    = "Technic large motor"
	spikePrimeLargeMotorName = "SPIKE Prime large motor"
)

// GetMotorName gets the name of the motor based on its type
func (m *ActiveMotor) GetMotorName() string {
	switch m.GetSensorType() {
	case models.TechnicXLMotorID:
		return technicXLMotorName
	case models.TechnicLargeMotorID:
		return technicLargeMotorName
	case models.MediumLinearMotor:
		return "Medium linear motor"
	case models.SpikeEssentialSmallAngularMotor:
		return "SPIKE Essential small angular motor"
	case models.SpikePrimeLargeMotor:
		return spikePrimeLargeMotorName
	case models.SpikePrimeMediumMotor:
		return "SPIKE Prime medium motor"
	case models.TechnicMediumAngularMotor:
		return "Technical medium angular motor"
	case models.TechnicMotor:
		return "Technical motor"
	default:
		return "Unknown active motor"
	}
}

// GetSensorName gets the sensor name (implements sensors.Sensor interface)
func (m *ActiveMotor) GetSensorName() string {
	return m.GetMotorName()
}

// GetSpeed gets the current speed
func (m *ActiveMotor) GetSpeed() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.speed
}

// Speed gets the current speed (alias for GetSpeed)
func (m *ActiveMotor) Speed() int {
	return m.GetSpeed()
}

// SetSpeed sets the target speed
func (m *ActiveMotor) SetSpeed(speed int) error {
	if speed < -100 || speed > 100 {
		return fmt.Errorf("speed must be between -100 and 100, got %d", speed)
	}

	m.mu.Lock()
	m.targetSpeed = speed
	m.mu.Unlock()

	return nil
}

// GetPosition gets the current tachometer count
func (m *ActiveMotor) GetPosition() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tacho
}

// SetPosition sets the current tachometer count (internal use)
func (m *ActiveMotor) SetPosition(position int) {
	m.mu.Lock()
	m.tacho = position
	m.mu.Unlock()
}

// GetAbsolutePosition gets the current absolute tachometer count
func (m *ActiveMotor) GetAbsolutePosition() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.absoluteTacho
}

// SetAbsolutePosition sets the current absolute tachometer count (internal use)
func (m *ActiveMotor) SetAbsolutePosition(position int) {
	m.mu.Lock()
	m.absoluteTacho = position
	m.mu.Unlock()
}

// SetSpeed sets the current speed (internal use)
func (m *ActiveMotor) SetSpeedInternal(speed int) {
	m.mu.Lock()
	m.speed = speed
	m.mu.Unlock()
}

// Start starts the motor with the current target speed
func (m *ActiveMotor) Start() error {
	m.mu.RLock()
	speed := m.targetSpeed
	m.mu.RUnlock()

	return m.GetBrick().SetPowerLevel(m.GetPort(), speed)
}

// StartWithSpeed starts the motor with the specified speed
func (m *ActiveMotor) StartWithSpeed(speed int) error {
	if err := m.SetSpeed(speed); err != nil {
		return err
	}
	return m.Start()
}

// Stop stops the motor
func (m *ActiveMotor) Stop() error {
	if err := m.SetSpeed(0); err != nil {
		return err
	}
	return m.Start() // Apply the zero speed
}

// SetBias sets the bias of the motor
func (m *ActiveMotor) SetBias(bias float64) error {
	if bias < 0 || bias > 1 {
		return fmt.Errorf("bias must be between 0 and 1, got %f", bias)
	}

	return m.GetBrick().SetMotorBias(m.GetPort(), bias)
}

// SetPowerLimit sets the power consumption limit
func (m *ActiveMotor) SetPowerLimit(plimit float64) error {
	if plimit < 0 || plimit > 1 {
		return fmt.Errorf("power limit must be between 0 and 1, got %f", plimit)
	}

	if err := m.GetBrick().SetMotorLimits(m.GetPort(), plimit); err != nil {
		return err
	}

	m.mu.Lock()
	m.powerLimit = plimit
	m.mu.Unlock()

	return nil
}

// GetPowerLimit gets the current power limit
func (m *ActiveMotor) GetPowerLimit() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.powerLimit
}

// Float floats the motor and stops all constraints on it
func (m *ActiveMotor) Float() error {
	return m.GetBrick().FloatMotor(m.GetPort())
}

// MoveToAbsolutePosition runs the motor to an absolute position
func (m *ActiveMotor) MoveToAbsolutePosition(ctx context.Context, targetPosition int, way models.PositionWay, blocking bool) error {
	m.mu.RLock()
	speed := m.targetSpeed
	m.mu.RUnlock()

	return m.GetBrick().MoveMotorToAbsolutePosition(ctx, m.GetPort(), targetPosition, way, speed, blocking)
}

// MoveForSeconds runs the motor for an amount of seconds
func (m *ActiveMotor) MoveForSeconds(ctx context.Context, seconds float64, blocking bool) error {
	m.mu.RLock()
	speed := m.targetSpeed
	m.mu.RUnlock()

	return m.GetBrick().MoveMotorForSeconds(ctx, m.GetPort(), seconds, speed, blocking)
}

// MoveToPosition runs the motor to an absolute position
func (m *ActiveMotor) MoveToPosition(ctx context.Context, targetPosition int, blocking bool) error {
	m.mu.RLock()
	speed := m.targetSpeed
	m.mu.RUnlock()

	return m.GetBrick().MoveMotorToPosition(ctx, m.GetPort(), targetPosition, speed, blocking)
}

// MoveForDegrees runs the motor for a specific number of degrees
func (m *ActiveMotor) MoveForDegrees(ctx context.Context, targetPosition int, blocking bool) error {
	m.mu.RLock()
	speed := m.targetSpeed
	m.mu.RUnlock()

	return m.GetBrick().MoveMotorForDegrees(ctx, m.GetPort(), targetPosition, speed, blocking)
}

// GetBrick gets the brick interface (helper method)
func (m *ActiveMotor) GetBrick() sensors.BrickInterface {
	return m.BaseSensor.GetBrick()
}

// UpdateFromSensorData updates motor values from raw sensor data
func (m *ActiveMotor) UpdateFromSensorData(data []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(data) >= 3 {
		// Parse tacho, absolute tacho, and speed from data
		if tacho, err := strconv.Atoi(data[0]); err == nil {
			m.tacho = tacho
		}
		if absoluteTacho, err := strconv.Atoi(data[1]); err == nil {
			m.absoluteTacho = absoluteTacho
		}
		if speed, err := strconv.Atoi(data[2]); err == nil {
			m.speed = speed
		}
	}

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this motor
func (m *ActiveMotor) GetTriggerFlag() *bool {
	// ActiveMotor doesn't have a trigger flag, return nil
	return nil
}
