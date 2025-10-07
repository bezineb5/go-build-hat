package sensors

import (
	"context"
	"reflect"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// PropertyChangedEvent represents a property change event
type PropertyChangedEvent struct {
	PropertyName string
	OldValue     interface{}
	NewValue     interface{}
}

// PropertyChangedHandler is a function type for handling property change events
type PropertyChangedHandler func(event PropertyChangedEvent)

// PropertyUpdatedHandler is a function type for handling property update events
type PropertyUpdatedHandler func(propertyName string)

// Sensor represents the base interface for all sensors
type Sensor interface {
	// GetSensorName gets the name of the sensor
	GetSensorName() string

	// GetPort gets the sensor port
	GetPort() models.SensorPort

	// GetSensorType gets the sensor type
	GetSensorType() models.SensorType

	// IsConnected gets true if the sensor is connected
	IsConnected() bool
}

// BaseSensor provides common functionality for all sensors
type BaseSensor struct {
	brick      BrickInterface
	port       models.SensorPort
	sensorType models.SensorType
	connected  bool
	mu         sync.RWMutex

	// Event handlers
	propertyChangedHandlers []PropertyChangedHandler
	propertyUpdatedHandlers []PropertyUpdatedHandler
	eventMu                 sync.RWMutex
}

// NewBaseSensor creates a new base sensor
func NewBaseSensor(brick BrickInterface, port models.SensorPort, sensorType models.SensorType) *BaseSensor {
	return &BaseSensor{
		brick:      brick,
		port:       port,
		sensorType: sensorType,
		connected:  true,
	}
}

const (
	genericSensorName = "Generic sensor"
)

// GetSensorName gets the name of the sensor
func (s *BaseSensor) GetSensorName() string {
	return genericSensorName
}

// GetPort gets the sensor port
func (s *BaseSensor) GetPort() models.SensorPort {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.port
}

// GetSensorType gets the sensor type
func (s *BaseSensor) GetSensorType() models.SensorType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sensorType
}

// IsConnected gets true if the sensor is connected
func (s *BaseSensor) IsConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connected
}

// SetConnected sets the connection status
func (s *BaseSensor) SetConnected(connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connected = connected
}

// GetBrick gets the brick interface
func (s *BaseSensor) GetBrick() BrickInterface {
	return s.brick
}

// AddPropertyChangedHandler adds a handler for property change events
func (s *BaseSensor) AddPropertyChangedHandler(handler PropertyChangedHandler) {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.propertyChangedHandlers = append(s.propertyChangedHandlers, handler)
}

// RemovePropertyChangedHandler removes a handler for property change events
func (s *BaseSensor) RemovePropertyChangedHandler(handler PropertyChangedHandler) {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	for i, h := range s.propertyChangedHandlers {
		// Use reflect to compare function pointers
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			s.propertyChangedHandlers = append(s.propertyChangedHandlers[:i], s.propertyChangedHandlers[i+1:]...)
			break
		}
	}
}

// AddPropertyUpdatedHandler adds a handler for property update events
func (s *BaseSensor) AddPropertyUpdatedHandler(handler PropertyUpdatedHandler) {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.propertyUpdatedHandlers = append(s.propertyUpdatedHandlers, handler)
}

// RemovePropertyUpdatedHandler removes a handler for property update events
func (s *BaseSensor) RemovePropertyUpdatedHandler(handler PropertyUpdatedHandler) {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	for i, h := range s.propertyUpdatedHandlers {
		// Use reflect to compare function pointers
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			s.propertyUpdatedHandlers = append(s.propertyUpdatedHandlers[:i], s.propertyUpdatedHandlers[i+1:]...)
			break
		}
	}
}

// OnPropertyChanged raises the property changed event
func (s *BaseSensor) OnPropertyChanged(propertyName string, oldValue, newValue interface{}) {
	s.eventMu.RLock()
	handlers := make([]PropertyChangedHandler, len(s.propertyChangedHandlers))
	copy(handlers, s.propertyChangedHandlers)
	s.eventMu.RUnlock()

	event := PropertyChangedEvent{
		PropertyName: propertyName,
		OldValue:     oldValue,
		NewValue:     newValue,
	}

	for _, handler := range handlers {
		go handler(event) // Run handlers in goroutines to avoid blocking
	}
}

// OnPropertyUpdated raises the property updated event
func (s *BaseSensor) OnPropertyUpdated(propertyName string) {
	s.eventMu.RLock()
	handlers := make([]PropertyUpdatedHandler, len(s.propertyUpdatedHandlers))
	copy(handlers, s.propertyUpdatedHandlers)
	s.eventMu.RUnlock()

	for _, handler := range handlers {
		go handler(propertyName) // Run handlers in goroutines to avoid blocking
	}
}

// BrickInterface represents the interface that sensors need from the Brick
type BrickInterface interface {
	// Motor control methods
	SetPowerLevel(port models.SensorPort, powerPercent int) error
	SetMotorLimits(port models.SensorPort, powerLimit float64) error
	SetMotorBias(port models.SensorPort, bias float64) error
	MoveMotorForSeconds(ctx context.Context, port models.SensorPort, seconds float64, speed int, blocking bool) error
	MoveMotorToPosition(ctx context.Context, port models.SensorPort, targetPosition, speed int, blocking bool) error
	MoveMotorToAbsolutePosition(ctx context.Context, port models.SensorPort, targetPosition int, way models.PositionWay, speed int, blocking bool) error
	MoveMotorForDegrees(ctx context.Context, port models.SensorPort, targetPosition, speed int, blocking bool) error
	FloatMotor(port models.SensorPort) error

	// Sensor control methods
	SelectModeAndRead(port models.SensorPort, mode int, readOnce bool) error
	SelectCombiModesAndRead(port models.SensorPort, modes []int, readOnce bool) error
	StopContinuousReadingSensor(port models.SensorPort) error
	SwitchSensorOn(port models.SensorPort) error
	SwitchSensorOff(port models.SensorPort) error
	WriteBytesToSensor(port models.SensorPort, data []byte, singleHeader bool) error
	SendRawCommand(command string) error
}
