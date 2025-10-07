package sensors

import (
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// ButtonSensor represents a simple passive button sensor
type ButtonSensor struct {
	*BaseSensor
	isPressed   bool
	triggerFlag bool
	mu          sync.RWMutex
}

// NewButtonSensor creates a new button sensor
func NewButtonSensor(brick BrickInterface, port models.SensorPort) *ButtonSensor {
	return &ButtonSensor{
		BaseSensor: NewBaseSensor(brick, port, models.ButtonOrTouchSensor),
		isPressed:  false,
	}
}

const (
	buttonSensorName = "Button sensor"
)

// GetSensorName gets the name of the sensor
func (s *ButtonSensor) GetSensorName() string {
	return buttonSensorName
}

// IsPressed gets true when the button is pressed
func (s *ButtonSensor) IsPressed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isPressed
}

// SetIsPressed sets the pressed state
func (s *ButtonSensor) SetIsPressed(pressed bool) {
	s.mu.Lock()
	oldValue := s.isPressed
	if s.isPressed != pressed {
		s.isPressed = pressed
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("IsPressed", oldValue, pressed)
		s.OnPropertyUpdated("IsPressed")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("IsPressed")
	}
}

// UpdateFromSensorData updates sensor values from raw sensor data
func (s *ButtonSensor) UpdateFromSensorData(data []string) error {
	s.mu.Lock()
	oldValue := s.isPressed

	if len(data) >= 1 {
		newPressed := data[0] == "1"
		if s.isPressed != newPressed {
			s.isPressed = newPressed
			s.mu.Unlock()

			// Trigger property change events
			s.OnPropertyChanged("IsPressed", oldValue, newPressed)
			s.OnPropertyUpdated("IsPressed")
		} else {
			s.mu.Unlock()
			// Still trigger property updated event even if value didn't change
			s.OnPropertyUpdated("IsPressed")
		}
	} else {
		s.mu.Unlock()
	}

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this sensor
func (s *ButtonSensor) GetTriggerFlag() *bool {
	return &s.triggerFlag
}
