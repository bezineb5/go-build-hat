package motors

import (
	"context"
	"testing"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// MockBrick implements sensors.BrickInterface for testing
type MockBrick struct {
	motorPowerCalls  []MotorPowerCall
	motorLimitsCalls []MotorLimitsCall
	motorBiasCalls   []MotorBiasCall
	floatMotorCalls  []models.SensorPort
	moveMotorCalls   []MoveMotorCall
}

type MotorPowerCall struct {
	Port         models.SensorPort
	PowerPercent int
}

type MotorLimitsCall struct {
	Port       models.SensorPort
	PowerLimit float64
}

type MotorBiasCall struct {
	Port models.SensorPort
	Bias float64
}

type MoveMotorCall struct {
	Port           models.SensorPort
	Seconds        float64
	Speed          int
	Blocking       bool
	TargetPosition int
	Way            models.PositionWay
	Degrees        int
	Type           string // "seconds", "position", "absolute", "degrees"
	Ctx            context.Context
}

func (m *MockBrick) SetPowerLevel(port models.SensorPort, powerPercent int) error {
	m.motorPowerCalls = append(m.motorPowerCalls, MotorPowerCall{port, powerPercent})
	return nil
}

func (m *MockBrick) SetMotorPower(port models.SensorPort, powerPercent int) error {
	// Backward compatibility - delegate to SetPowerLevel
	return m.SetPowerLevel(port, powerPercent)
}

func (m *MockBrick) SetMotorLimits(port models.SensorPort, powerLimit float64) error {
	m.motorLimitsCalls = append(m.motorLimitsCalls, MotorLimitsCall{port, powerLimit})
	return nil
}

func (m *MockBrick) SetMotorBias(port models.SensorPort, bias float64) error {
	m.motorBiasCalls = append(m.motorBiasCalls, MotorBiasCall{port, bias})
	return nil
}

func (m *MockBrick) MoveMotorForSeconds(port models.SensorPort, seconds float64, speed int, blocking bool, ctx context.Context) error {
	m.moveMotorCalls = append(m.moveMotorCalls, MoveMotorCall{
		Port:     port,
		Seconds:  seconds,
		Speed:    speed,
		Blocking: blocking,
		Type:     "seconds",
		Ctx:      ctx,
	})
	return nil
}

func (m *MockBrick) MoveMotorToPosition(port models.SensorPort, targetPosition int, speed int, blocking bool, ctx context.Context) error {
	m.moveMotorCalls = append(m.moveMotorCalls, MoveMotorCall{
		Port:           port,
		TargetPosition: targetPosition,
		Speed:          speed,
		Blocking:       blocking,
		Type:           "position",
		Ctx:            ctx,
	})
	return nil
}

func (m *MockBrick) MoveMotorToAbsolutePosition(port models.SensorPort, targetPosition int, way models.PositionWay, speed int, blocking bool, ctx context.Context) error {
	m.moveMotorCalls = append(m.moveMotorCalls, MoveMotorCall{
		Port:           port,
		TargetPosition: targetPosition,
		Way:            way,
		Speed:          speed,
		Blocking:       blocking,
		Type:           "absolute",
		Ctx:            ctx,
	})
	return nil
}

func (m *MockBrick) MoveMotorForDegrees(port models.SensorPort, targetPosition int, speed int, blocking bool, ctx context.Context) error {
	m.moveMotorCalls = append(m.moveMotorCalls, MoveMotorCall{
		Port:           port,
		TargetPosition: targetPosition,
		Speed:          speed,
		Blocking:       blocking,
		Type:           "degrees",
		Ctx:            ctx,
	})
	return nil
}

func (m *MockBrick) FloatMotor(port models.SensorPort) error {
	m.floatMotorCalls = append(m.floatMotorCalls, port)
	return nil
}

// Sensor control methods
func (m *MockBrick) SelectModeAndRead(port models.SensorPort, mode int, readOnce bool) error {
	return nil
}

func (m *MockBrick) SelectCombiModesAndRead(port models.SensorPort, modes []int, readOnce bool) error {
	return nil
}

func (m *MockBrick) StopContinuousReadingSensor(port models.SensorPort) error {
	return nil
}

func (m *MockBrick) SwitchSensorOn(port models.SensorPort) error {
	return nil
}

func (m *MockBrick) SwitchSensorOff(port models.SensorPort) error {
	return nil
}

func (m *MockBrick) WriteBytesToSensor(port models.SensorPort, data []byte, singleHeader bool) error {
	return nil
}

func (m *MockBrick) SendRawCommand(command string) error {
	return nil
}

func TestActiveMotor(t *testing.T) {
	mockBrick := &MockBrick{}
	motor := NewActiveMotor(mockBrick, models.PortA, models.SpikePrimeLargeMotor)

	// Test basic properties
	if motor.GetPort() != models.PortA {
		t.Error("Expected port to be PortA")
	}

	if motor.GetSensorType() != models.SpikePrimeLargeMotor {
		t.Error("Expected sensor type to be SpikePrimeLargeMotor")
	}

	if motor.GetMotorName() != "SPIKE Prime large motor" {
		t.Error("Expected motor name to be 'SPIKE Prime large motor'")
	}

	// Test speed setting
	if err := motor.SetSpeed(50); err != nil {
		t.Errorf("Failed to set speed: %v", err)
	}

	if motor.GetSpeed() != 0 { // Speed should still be 0 until Start() is called
		t.Error("Expected speed to be 0 before starting")
	}

	// Test starting motor
	if err := motor.Start(); err != nil {
		t.Errorf("Failed to start motor: %v", err)
	}

	if len(mockBrick.motorPowerCalls) != 1 {
		t.Error("Expected one motor power call")
	}

	if mockBrick.motorPowerCalls[0].PowerPercent != 50 {
		t.Error("Expected power to be 50")
	}

	// Test stopping motor
	if err := motor.Stop(); err != nil {
		t.Errorf("Failed to stop motor: %v", err)
	}

	if len(mockBrick.motorPowerCalls) != 2 {
		t.Error("Expected two motor power calls")
	}

	if mockBrick.motorPowerCalls[1].PowerPercent != 0 {
		t.Error("Expected power to be 0 when stopping")
	}

	// Test bias setting (initial call + our call = 2 total)
	if err := motor.SetBias(0.5); err != nil {
		t.Errorf("Failed to set bias: %v", err)
	}

	if len(mockBrick.motorBiasCalls) != 2 {
		t.Error("Expected two bias calls (initial + our call)")
	}

	// Test power limit setting (initial call + our call = 2 total)
	if err := motor.SetPowerLimit(0.8); err != nil {
		t.Errorf("Failed to set power limit: %v", err)
	}

	if len(mockBrick.motorLimitsCalls) != 2 {
		t.Error("Expected two power limit calls (initial + our call)")
	}

	// Test floating motor
	if err := motor.Float(); err != nil {
		t.Errorf("Failed to float motor: %v", err)
	}

	if len(mockBrick.floatMotorCalls) != 1 {
		t.Error("Expected one float motor call")
	}

	// Test position methods
	motor.SetPosition(100)
	if motor.GetPosition() != 100 {
		t.Error("Expected position to be 100")
	}

	motor.SetAbsolutePosition(200)
	if motor.GetAbsolutePosition() != 200 {
		t.Error("Expected absolute position to be 200")
	}

	// Test movement methods
	ctx := context.Background()
	if err := motor.MoveForSeconds(2.0, true, ctx); err != nil {
		t.Errorf("Failed to move motor for seconds: %v", err)
	}

	if len(mockBrick.moveMotorCalls) != 1 {
		t.Error("Expected one move motor call")
	}

	if mockBrick.moveMotorCalls[0].Type != "seconds" {
		t.Error("Expected move type to be 'seconds'")
	}
}

func TestPassiveMotor(t *testing.T) {
	mockBrick := &MockBrick{}
	motor := NewPassiveMotor(mockBrick, models.PortB, models.SystemMediumMotor)

	// Test basic properties
	if motor.GetPort() != models.PortB {
		t.Error("Expected port to be PortB")
	}

	if motor.GetSensorType() != models.SystemMediumMotor {
		t.Error("Expected sensor type to be SystemMediumMotor")
	}

	if motor.GetMotorName() != "System medium Motor" {
		t.Error("Expected motor name to be 'System medium Motor'")
	}

	// Test speed setting
	if err := motor.SetSpeed(75); err != nil {
		t.Errorf("Failed to set speed: %v", err)
	}

	if motor.GetSpeed() != 75 {
		t.Error("Expected speed to be 75")
	}

	// Test starting motor
	if err := motor.Start(); err != nil {
		t.Errorf("Failed to start motor: %v", err)
	}

	if !motor.IsRunning() {
		t.Error("Expected motor to be running")
	}

	if len(mockBrick.motorPowerCalls) != 1 {
		t.Error("Expected one motor power call")
	}

	if mockBrick.motorPowerCalls[0].PowerPercent != 75 {
		t.Error("Expected power to be 75")
	}

	// Test stopping motor
	if err := motor.Stop(); err != nil {
		t.Errorf("Failed to stop motor: %v", err)
	}

	if motor.IsRunning() {
		t.Error("Expected motor to not be running")
	}

	if motor.GetSpeed() != 0 {
		t.Error("Expected speed to be 0 when stopped")
	}

	// Test bias setting (initial call + our call = 2 total)
	if err := motor.SetBias(0.4); err != nil {
		t.Errorf("Failed to set bias: %v", err)
	}

	if len(mockBrick.motorBiasCalls) != 2 {
		t.Error("Expected two bias calls (initial + our call)")
	}

	// Test power limit setting (initial call + our call = 2 total)
	if err := motor.SetPowerLimit(0.9); err != nil {
		t.Errorf("Failed to set power limit: %v", err)
	}

	if len(mockBrick.motorLimitsCalls) != 2 {
		t.Error("Expected two power limit calls (initial + our call)")
	}

	// Test floating motor
	if err := motor.Float(); err != nil {
		t.Errorf("Failed to float motor: %v", err)
	}

	if len(mockBrick.floatMotorCalls) != 1 {
		t.Error("Expected one float motor call")
	}
}

func TestMotorValidation(t *testing.T) {
	mockBrick := &MockBrick{}
	motor := NewActiveMotor(mockBrick, models.PortA, models.SpikePrimeLargeMotor)

	// Test invalid speed
	if err := motor.SetSpeed(150); err == nil {
		t.Error("Expected error for speed > 100")
	}

	if err := motor.SetSpeed(-150); err == nil {
		t.Error("Expected error for speed < -100")
	}

	// Test invalid bias
	if err := motor.SetBias(1.5); err == nil {
		t.Error("Expected error for bias > 1")
	}

	if err := motor.SetBias(-0.5); err == nil {
		t.Error("Expected error for bias < 0")
	}

	// Test invalid power limit
	if err := motor.SetPowerLimit(1.5); err == nil {
		t.Error("Expected error for power limit > 1")
	}

	if err := motor.SetPowerLimit(-0.5); err == nil {
		t.Error("Expected error for power limit < 0")
	}
}
