package sensors

import (
	"context"
	"testing"
	"time"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// MockBrickInterface is a mock implementation of BrickInterface for testing
type MockBrickInterface struct {
	// Motor control method calls
	SetMotorPowerCalls               []SetMotorPowerCall
	SetMotorLimitsCalls              []SetMotorLimitsCall
	SetMotorBiasCalls                []SetMotorBiasCall
	MoveMotorForSecondsCalls         []MoveMotorForSecondsCall
	MoveMotorToPositionCalls         []MoveMotorToPositionCall
	MoveMotorToAbsolutePositionCalls []MoveMotorToAbsolutePositionCall
	MoveMotorForDegreesCalls         []MoveMotorForDegreesCall
	FloatMotorCalls                  []FloatMotorCall

	// Sensor control method calls
	SelectModeAndReadCalls           []SelectModeAndReadCall
	SelectCombiModesAndReadCalls     []SelectCombiModesAndReadCall
	StopContinuousReadingSensorCalls []StopContinuousReadingSensorCall
	SwitchSensorOnCalls              []SwitchSensorOnCall
	SwitchSensorOffCalls             []SwitchSensorOffCall
	WriteBytesToSensorCalls          []WriteBytesToSensorCall
	SendRawCommandCalls              []SendRawCommandCall
}

type SetMotorPowerCall struct {
	Port         models.SensorPort
	PowerPercent int
}

type SetMotorLimitsCall struct {
	Port       models.SensorPort
	PowerLimit float64
}

type SetMotorBiasCall struct {
	Port models.SensorPort
	Bias float64
}

type MoveMotorForSecondsCall struct {
	Port     models.SensorPort
	Seconds  float64
	Speed    int
	Blocking bool
	Ctx      context.Context
}

type MoveMotorToPositionCall struct {
	Port           models.SensorPort
	TargetPosition int
	Speed          int
	Blocking       bool
	Ctx            context.Context
}

type MoveMotorToAbsolutePositionCall struct {
	Port           models.SensorPort
	TargetPosition int
	Way            models.PositionWay
	Speed          int
	Blocking       bool
	Ctx            context.Context
}

type MoveMotorForDegreesCall struct {
	Port           models.SensorPort
	TargetPosition int
	Speed          int
	Blocking       bool
	Ctx            context.Context
}

type FloatMotorCall struct {
	Port models.SensorPort
}

type SelectModeAndReadCall struct {
	Port     models.SensorPort
	Mode     int
	ReadOnce bool
}

type SelectCombiModesAndReadCall struct {
	Port     models.SensorPort
	Modes    []int
	ReadOnce bool
}

type StopContinuousReadingSensorCall struct {
	Port models.SensorPort
}

type SwitchSensorOnCall struct {
	Port models.SensorPort
}

type SwitchSensorOffCall struct {
	Port models.SensorPort
}

type WriteBytesToSensorCall struct {
	Port         models.SensorPort
	Data         []byte
	SingleHeader bool
}

type SendRawCommandCall struct {
	Command string
}

// Motor control methods
func (m *MockBrickInterface) SetPowerLevel(port models.SensorPort, powerPercent int) error {
	m.SetMotorPowerCalls = append(m.SetMotorPowerCalls, SetMotorPowerCall{port, powerPercent})
	return nil
}

func (m *MockBrickInterface) SetMotorPower(port models.SensorPort, powerPercent int) error {
	// Backward compatibility - delegate to SetPowerLevel
	return m.SetPowerLevel(port, powerPercent)
}

func (m *MockBrickInterface) SetMotorLimits(port models.SensorPort, powerLimit float64) error {
	m.SetMotorLimitsCalls = append(m.SetMotorLimitsCalls, SetMotorLimitsCall{port, powerLimit})
	return nil
}

func (m *MockBrickInterface) SetMotorBias(port models.SensorPort, bias float64) error {
	m.SetMotorBiasCalls = append(m.SetMotorBiasCalls, SetMotorBiasCall{port, bias})
	return nil
}

func (m *MockBrickInterface) MoveMotorForSeconds(port models.SensorPort, seconds float64, speed int, blocking bool, ctx context.Context) error {
	m.MoveMotorForSecondsCalls = append(m.MoveMotorForSecondsCalls, MoveMotorForSecondsCall{port, seconds, speed, blocking, ctx})
	return nil
}

func (m *MockBrickInterface) MoveMotorToPosition(port models.SensorPort, targetPosition int, speed int, blocking bool, ctx context.Context) error {
	m.MoveMotorToPositionCalls = append(m.MoveMotorToPositionCalls, MoveMotorToPositionCall{port, targetPosition, speed, blocking, ctx})
	return nil
}

func (m *MockBrickInterface) MoveMotorToAbsolutePosition(port models.SensorPort, targetPosition int, way models.PositionWay, speed int, blocking bool, ctx context.Context) error {
	m.MoveMotorToAbsolutePositionCalls = append(m.MoveMotorToAbsolutePositionCalls, MoveMotorToAbsolutePositionCall{port, targetPosition, way, speed, blocking, ctx})
	return nil
}

func (m *MockBrickInterface) MoveMotorForDegrees(port models.SensorPort, targetPosition int, speed int, blocking bool, ctx context.Context) error {
	m.MoveMotorForDegreesCalls = append(m.MoveMotorForDegreesCalls, MoveMotorForDegreesCall{port, targetPosition, speed, blocking, ctx})
	return nil
}

func (m *MockBrickInterface) FloatMotor(port models.SensorPort) error {
	m.FloatMotorCalls = append(m.FloatMotorCalls, FloatMotorCall{port})
	return nil
}

// Sensor control methods
func (m *MockBrickInterface) SelectModeAndRead(port models.SensorPort, mode int, readOnce bool) error {
	m.SelectModeAndReadCalls = append(m.SelectModeAndReadCalls, SelectModeAndReadCall{port, mode, readOnce})
	return nil
}

func (m *MockBrickInterface) SelectCombiModesAndRead(port models.SensorPort, modes []int, readOnce bool) error {
	m.SelectCombiModesAndReadCalls = append(m.SelectCombiModesAndReadCalls, SelectCombiModesAndReadCall{port, modes, readOnce})
	return nil
}

func (m *MockBrickInterface) StopContinuousReadingSensor(port models.SensorPort) error {
	m.StopContinuousReadingSensorCalls = append(m.StopContinuousReadingSensorCalls, StopContinuousReadingSensorCall{port})
	return nil
}

func (m *MockBrickInterface) SwitchSensorOn(port models.SensorPort) error {
	m.SwitchSensorOnCalls = append(m.SwitchSensorOnCalls, SwitchSensorOnCall{port})
	return nil
}

func (m *MockBrickInterface) SwitchSensorOff(port models.SensorPort) error {
	m.SwitchSensorOffCalls = append(m.SwitchSensorOffCalls, SwitchSensorOffCall{port})
	return nil
}

func (m *MockBrickInterface) WriteBytesToSensor(port models.SensorPort, data []byte, singleHeader bool) error {
	m.WriteBytesToSensorCalls = append(m.WriteBytesToSensorCalls, WriteBytesToSensorCall{port, data, singleHeader})
	return nil
}

func (m *MockBrickInterface) SendRawCommand(command string) error {
	m.SendRawCommandCalls = append(m.SendRawCommandCalls, SendRawCommandCall{command})
	return nil
}

func TestBaseSensor(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewBaseSensor(mockBrick, models.PortA, models.SpikePrimeLargeMotor)

	// Test basic properties
	if sensor.GetSensorName() != "Generic sensor" {
		t.Errorf("Expected sensor name 'Generic sensor', got '%s'", sensor.GetSensorName())
	}

	if sensor.GetPort() != models.PortA {
		t.Errorf("Expected port PortA, got %v", sensor.GetPort())
	}

	if sensor.GetSensorType() != models.SpikePrimeLargeMotor {
		t.Errorf("Expected sensor type SpikePrimeLargeMotor, got %v", sensor.GetSensorType())
	}

	if !sensor.IsConnected() {
		t.Error("Expected sensor to be connected")
	}

	// Test connection status change
	sensor.SetConnected(false)
	if sensor.IsConnected() {
		t.Error("Expected sensor to be disconnected")
	}

	if sensor.GetBrick() != mockBrick {
		t.Error("Expected brick to be the mock brick")
	}
}

func TestActiveSensor(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewActiveSensor(mockBrick, models.PortA, models.SpikePrimeColorSensor)

	// Test basic properties
	if sensor.GetSensorName() != "Active sensor" {
		t.Errorf("Expected sensor name 'Active sensor', got '%s'", sensor.GetSensorName())
	}

	// Test values as string
	values := []string{"P0C0", "123", "456"}
	sensor.SetValuesAsString(values)

	retrievedValues := sensor.ValuesAsString()
	if len(retrievedValues) != len(values) {
		t.Errorf("Expected %d values, got %d", len(values), len(retrievedValues))
	}

	for i, v := range values {
		if retrievedValues[i] != v {
			t.Errorf("Expected value '%s' at index %d, got '%s'", v, i, retrievedValues[i])
		}
	}

	// Test baud rate
	sensor.SetBaudRate(9600)
	if sensor.BaudRate() != 9600 {
		t.Errorf("Expected baud rate 9600, got %d", sensor.BaudRate())
	}

	// Test hardware version
	sensor.SetHardwareVersion(0x12345678)
	if sensor.HardwareVersion() != 0x12345678 {
		t.Errorf("Expected hardware version 0x12345678, got 0x%x", sensor.HardwareVersion())
	}

	// Test software version
	sensor.SetSoftwareVersion(0x87654321)
	if sensor.SoftwareVersion() != 0x87654321 {
		t.Errorf("Expected software version 0x87654321, got 0x%x", sensor.SoftwareVersion())
	}

	// Test mode details
	modeDetails := []models.ModeDetail{
		{Number: 0, Name: "Mode 0"},
		{Number: 1, Name: "Mode 1"},
	}
	sensor.SetModeDetails(modeDetails)

	retrievedDetails := sensor.ModeDetails()
	if len(retrievedDetails) != len(modeDetails) {
		t.Errorf("Expected %d mode details, got %d", len(modeDetails), len(retrievedDetails))
	}

	if sensor.NumberOfModes() != len(modeDetails) {
		t.Errorf("Expected %d modes, got %d", len(modeDetails), sensor.NumberOfModes())
	}

	// Test combi modes
	combiModes := []models.CombiModes{
		{Number: 0, Modes: []int{0, 1}},
		{Number: 1, Modes: []int{2, 3}},
	}
	sensor.SetCombiModes(combiModes)

	retrievedCombiModes := sensor.CombiModes()
	if len(retrievedCombiModes) != len(combiModes) {
		t.Errorf("Expected %d combi modes, got %d", len(combiModes), len(retrievedCombiModes))
	}

	// Test PID settings
	speedPid := models.RecommendedPid{Pid1: 1, Pid2: 2, Pid3: 3, Pid4: 4}
	sensor.SetSpeedPid(speedPid)
	if sensor.SpeedPid() != speedPid {
		t.Errorf("Expected speed PID %v, got %v", speedPid, sensor.SpeedPid())
	}

	positionPid := models.RecommendedPid{Pid1: 5, Pid2: 6, Pid3: 7, Pid4: 8}
	sensor.SetPositionPid(positionPid)
	if sensor.PositionPid() != positionPid {
		t.Errorf("Expected position PID %v, got %v", positionPid, sensor.PositionPid())
	}

	// Test sensor control methods
	err := sensor.SelectModeAndRead(0, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockBrick.SelectModeAndReadCalls) != 1 {
		t.Errorf("Expected 1 SelectModeAndRead call, got %d", len(mockBrick.SelectModeAndReadCalls))
	}

	call := mockBrick.SelectModeAndReadCalls[0]
	if call.Port != models.PortA || call.Mode != 0 || !call.ReadOnce {
		t.Errorf("Expected call with PortA, mode 0, readOnce true, got %v", call)
	}

	err = sensor.SelectCombiModesAndRead([]int{0, 1}, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockBrick.SelectCombiModesAndReadCalls) != 1 {
		t.Errorf("Expected 1 SelectCombiModesAndRead call, got %d", len(mockBrick.SelectCombiModesAndReadCalls))
	}

	err = sensor.StopReading()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockBrick.StopContinuousReadingSensorCalls) != 1 {
		t.Errorf("Expected 1 StopContinuousReadingSensor call, got %d", len(mockBrick.StopContinuousReadingSensorCalls))
	}

	err = sensor.SwitchOn()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockBrick.SwitchSensorOnCalls) != 1 {
		t.Errorf("Expected 1 SwitchSensorOn call, got %d", len(mockBrick.SwitchSensorOnCalls))
	}

	err = sensor.SwitchOff()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockBrick.SwitchSensorOffCalls) != 1 {
		t.Errorf("Expected 1 SwitchSensorOff call, got %d", len(mockBrick.SwitchSensorOffCalls))
	}

	data := []byte{0x01, 0x02, 0x03}
	err = sensor.WriteBytes(data, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockBrick.WriteBytesToSensorCalls) != 1 {
		t.Errorf("Expected 1 WriteBytesToSensor call, got %d", len(mockBrick.WriteBytesToSensorCalls))
	}
}

func TestButtonSensor(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewButtonSensor(mockBrick, models.PortA)

	// Test basic properties
	if sensor.GetSensorName() != "Button sensor" {
		t.Errorf("Expected sensor name 'Button sensor', got '%s'", sensor.GetSensorName())
	}

	if sensor.GetSensorType() != models.ButtonOrTouchSensor {
		t.Errorf("Expected sensor type ButtonOrTouchSensor, got %v", sensor.GetSensorType())
	}

	// Test pressed state
	if sensor.IsPressed() {
		t.Error("Expected button to not be pressed initially")
	}

	sensor.SetIsPressed(true)
	if !sensor.IsPressed() {
		t.Error("Expected button to be pressed")
	}

	sensor.SetIsPressed(false)
	if sensor.IsPressed() {
		t.Error("Expected button to not be pressed")
	}
}

func TestColorSensor(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewColorSensor(mockBrick, models.PortA, models.SpikePrimeColorSensor)

	// Test that initialization sent the proper command for SpikePrimeColorSensor
	if len(mockBrick.SendRawCommandCalls) != 1 {
		t.Errorf("Expected 1 SendRawCommand call during initialization, got %d", len(mockBrick.SendRawCommandCalls))
	}

	expectedCommand := "port 0 ; plimit 1 ; set -1\r"
	if mockBrick.SendRawCommandCalls[0].Command != expectedCommand {
		t.Errorf("Expected command '%s', got '%s'", expectedCommand, mockBrick.SendRawCommandCalls[0].Command)
	}

	// Test basic properties
	if sensor.GetSensorName() != "Color sensor" {
		t.Errorf("Expected sensor name 'Color sensor', got '%s'", sensor.GetSensorName())
	}

	// Test color
	color := sensor.Color()
	if color.R != 0 || color.G != 0 || color.B != 0 || color.A != 255 {
		t.Errorf("Expected initial color to be black with alpha 255, got %v", color)
	}

	newColor := struct{ R, G, B, A uint8 }{R: 255, G: 128, B: 64, A: 255}
	sensor.SetColor(newColor)

	retrievedColor := sensor.Color()
	if retrievedColor.R != newColor.R || retrievedColor.G != newColor.G || retrievedColor.B != newColor.B || retrievedColor.A != newColor.A {
		t.Errorf("Expected color %v, got %v", newColor, retrievedColor)
	}

	// Test color detection
	if sensor.IsColorDetected() {
		t.Error("Expected color to not be detected initially")
	}

	sensor.SetIsColorDetected(true)
	if !sensor.IsColorDetected() {
		t.Error("Expected color to be detected")
	}

	// Test reflected light
	sensor.SetReflectedLight(500)
	if sensor.ReflectedLight() != 500 {
		t.Errorf("Expected reflected light 500, got %d", sensor.ReflectedLight())
	}

	// Test ambient light
	sensor.SetAmbientLight(300)
	if sensor.AmbientLight() != 300 {
		t.Errorf("Expected ambient light 300, got %d", sensor.AmbientLight())
	}
}

func TestForceSensor(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewForceSensor(mockBrick, models.PortA)

	// Test basic properties
	if sensor.GetSensorName() != "SPIKE force sensor" {
		t.Errorf("Expected sensor name 'SPIKE force sensor', got '%s'", sensor.GetSensorName())
	}

	if sensor.GetSensorType() != models.SpikePrimeForceSensor {
		t.Errorf("Expected sensor type SpikePrimeForceSensor, got %v", sensor.GetSensorType())
	}

	// Test force
	sensor.SetForce(100)
	if sensor.Force() != 100 {
		t.Errorf("Expected force 100, got %d", sensor.Force())
	}

	// Test pressed state
	if sensor.IsPressed() {
		t.Error("Expected sensor to not be pressed initially")
	}

	sensor.SetIsPressed(true)
	if !sensor.IsPressed() {
		t.Error("Expected sensor to be pressed")
	}

	// Test continuous measurement
	if sensor.ContinuousMeasurement() {
		t.Error("Expected continuous measurement to be false initially")
	}

	err := sensor.SetContinuousMeasurement(true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !sensor.ContinuousMeasurement() {
		t.Error("Expected continuous measurement to be true")
	}

	if len(mockBrick.SelectModeAndReadCalls) != 1 {
		t.Errorf("Expected 1 SelectModeAndRead call, got %d", len(mockBrick.SelectModeAndReadCalls))
	}

	err = sensor.SetContinuousMeasurement(false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sensor.ContinuousMeasurement() {
		t.Error("Expected continuous measurement to be false")
	}

	if len(mockBrick.StopContinuousReadingSensorCalls) != 1 {
		t.Errorf("Expected 1 StopContinuousReadingSensor call, got %d", len(mockBrick.StopContinuousReadingSensorCalls))
	}
}

func TestUltrasonicDistanceSensor(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewUltrasonicDistanceSensor(mockBrick, models.PortA)

	// Test that initialization sent the proper command for UltrasonicDistanceSensor
	if len(mockBrick.SendRawCommandCalls) != 1 {
		t.Errorf("Expected 1 SendRawCommand call during initialization, got %d", len(mockBrick.SendRawCommandCalls))
	}

	expectedCommand := "port 0 ; plimit 1 ; set -1\r"
	if mockBrick.SendRawCommandCalls[0].Command != expectedCommand {
		t.Errorf("Expected command '%s', got '%s'", expectedCommand, mockBrick.SendRawCommandCalls[0].Command)
	}

	// Test basic properties
	if sensor.GetSensorName() != "SPIKE distance sensor" {
		t.Errorf("Expected sensor name 'SPIKE distance sensor', got '%s'", sensor.GetSensorName())
	}

	if sensor.GetSensorType() != models.SpikePrimeUltrasonicDistanceSensor {
		t.Errorf("Expected sensor type SpikePrimeUltrasonicDistanceSensor, got %v", sensor.GetSensorType())
	}

	// Test distance
	sensor.SetDistance(150)
	if sensor.Distance() != 150 {
		t.Errorf("Expected distance 150, got %d", sensor.Distance())
	}

	// Test continuous measurement
	if sensor.ContinuousMeasurement() {
		t.Error("Expected continuous measurement to be false initially")
	}

	err := sensor.SetContinuousMeasurement(true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !sensor.ContinuousMeasurement() {
		t.Error("Expected continuous measurement to be true")
	}

	if len(mockBrick.SelectModeAndReadCalls) != 1 {
		t.Errorf("Expected 1 SelectModeAndRead call, got %d", len(mockBrick.SelectModeAndReadCalls))
	}

	err = sensor.SetContinuousMeasurement(false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sensor.ContinuousMeasurement() {
		t.Error("Expected continuous measurement to be false")
	}

	if len(mockBrick.StopContinuousReadingSensorCalls) != 1 {
		t.Errorf("Expected 1 StopContinuousReadingSensor call, got %d", len(mockBrick.StopContinuousReadingSensorCalls))
	}
}

func TestPassiveLight(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewPassiveLight(mockBrick, models.PortA)

	// Test basic properties
	if sensor.GetSensorName() != "Passive light" {
		t.Errorf("Expected sensor name 'Passive light', got '%s'", sensor.GetSensorName())
	}

	if sensor.GetSensorType() != models.SimpleLights {
		t.Errorf("Expected sensor type SimpleLights, got %v", sensor.GetSensorType())
	}

	// Test brightness
	if sensor.Brightness() != 0 {
		t.Errorf("Expected initial brightness 0, got %d", sensor.Brightness())
	}

	err := sensor.SetBrightness(50)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sensor.Brightness() != 50 {
		t.Errorf("Expected brightness 50, got %d", sensor.Brightness())
	}

	// Test brightness validation
	err = sensor.SetBrightness(-1)
	if err == nil {
		t.Error("Expected error for negative brightness")
	}

	err = sensor.SetBrightness(101)
	if err == nil {
		t.Error("Expected error for brightness > 100")
	}

	// Test on/off state
	if sensor.IsOn() {
		t.Error("Expected light to be off initially")
	}

	err = sensor.TurnOn()
	if err != nil {
		t.Errorf("Expected no error turning on, got %v", err)
	}
	if !sensor.IsOn() {
		t.Error("Expected light to be on")
	}

	err = sensor.TurnOff()
	if err != nil {
		t.Errorf("Expected no error turning off, got %v", err)
	}
	if sensor.IsOn() {
		t.Error("Expected light to be off")
	}

	err = sensor.SetIsOn(true)
	if err != nil {
		t.Errorf("Expected no error setting on, got %v", err)
	}
	if !sensor.IsOn() {
		t.Error("Expected light to be on")
	}
}

func TestButtonSensorEvents(t *testing.T) {
	mockBrick := &MockBrickInterface{}
	sensor := NewButtonSensor(mockBrick, models.PortA)

	// Track events with channels for synchronization
	propertyChangedEvents := make(chan PropertyChangedEvent, 10)
	propertyUpdatedEvents := make(chan string, 10)

	// Add event handlers
	propertyChangedHandler := func(event PropertyChangedEvent) {
		propertyChangedEvents <- event
	}
	propertyUpdatedHandler := func(propertyName string) {
		propertyUpdatedEvents <- propertyName
	}

	sensor.AddPropertyChangedHandler(propertyChangedHandler)
	sensor.AddPropertyUpdatedHandler(propertyUpdatedHandler)

	// Test initial state
	if sensor.IsPressed() {
		t.Error("Expected button to not be pressed initially")
	}

	// Test setting pressed state
	sensor.SetIsPressed(true)
	if !sensor.IsPressed() {
		t.Error("Expected button to be pressed")
	}

	// Wait for events with timeout
	select {
	case event := <-propertyChangedEvents:
		if event.PropertyName != "IsPressed" {
			t.Errorf("Expected property name 'IsPressed', got '%s'", event.PropertyName)
		}
		if event.OldValue != false {
			t.Errorf("Expected old value false, got %v", event.OldValue)
		}
		if event.NewValue != true {
			t.Errorf("Expected new value true, got %v", event.NewValue)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected property changed event, but got timeout")
	}

	select {
	case propertyName := <-propertyUpdatedEvents:
		if propertyName != "IsPressed" {
			t.Errorf("Expected property name 'IsPressed', got '%s'", propertyName)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected property updated event, but got timeout")
	}

	// Test setting same value (should trigger updated but not changed)
	sensor.SetIsPressed(true)

	// Should not get a property changed event
	select {
	case <-propertyChangedEvents:
		t.Error("Expected no property changed event for same value")
	case <-time.After(50 * time.Millisecond):
		// This is expected - no property changed event
	}

	// Should get a property updated event
	select {
	case propertyName := <-propertyUpdatedEvents:
		if propertyName != "IsPressed" {
			t.Errorf("Expected property name 'IsPressed', got '%s'", propertyName)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected property updated event, but got timeout")
	}

	// Test UpdateFromSensorData
	err := sensor.UpdateFromSensorData([]string{"0"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sensor.IsPressed() {
		t.Error("Expected button to not be pressed after UpdateFromSensorData")
	}

	// Wait for events from UpdateFromSensorData
	select {
	case event := <-propertyChangedEvents:
		if event.PropertyName != "IsPressed" {
			t.Errorf("Expected property name 'IsPressed', got '%s'", event.PropertyName)
		}
		if event.OldValue != true {
			t.Errorf("Expected old value true, got %v", event.OldValue)
		}
		if event.NewValue != false {
			t.Errorf("Expected new value false, got %v", event.NewValue)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected property changed event from UpdateFromSensorData, but got timeout")
	}

	select {
	case propertyName := <-propertyUpdatedEvents:
		if propertyName != "IsPressed" {
			t.Errorf("Expected property name 'IsPressed', got '%s'", propertyName)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected property updated event from UpdateFromSensorData, but got timeout")
	}

	// Test removing handlers
	sensor.RemovePropertyChangedHandler(propertyChangedHandler)
	sensor.RemovePropertyUpdatedHandler(propertyUpdatedHandler)

	sensor.SetIsPressed(true)

	// Should not get any events after removing handlers
	select {
	case <-propertyChangedEvents:
		t.Error("Expected no property changed event after removing handler")
	case <-time.After(50 * time.Millisecond):
		// This is expected - no events after removing handlers
	}

	select {
	case <-propertyUpdatedEvents:
		t.Error("Expected no property updated event after removing handler")
	case <-time.After(50 * time.Millisecond):
		// This is expected - no events after removing handlers
	}
}
