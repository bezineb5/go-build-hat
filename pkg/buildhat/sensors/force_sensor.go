package sensors

import (
	"errors"
	"strconv"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// ForceSensor represents a Spike force sensor
type ForceSensor struct {
	*ActiveSensor

	force               int
	hasForceUpdated     bool
	continuous          bool
	isPressed           bool
	hasIsPressedUpdated bool

	mu sync.RWMutex
}

// NewForceSensor creates a new force sensor
func NewForceSensor(brick BrickInterface, port models.SensorPort) *ForceSensor {
	return &ForceSensor{
		ActiveSensor: NewActiveSensor(brick, port, models.SpikePrimeForceSensor),
		force:        0,
		continuous:   false,
		isPressed:    false,
	}
}

// GetSensorName gets the name of the sensor
func (s *ForceSensor) GetSensorName() string {
	return "SPIKE force sensor"
}

// Force gets the force in Newtons
func (s *ForceSensor) Force() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.force
}

// SetForce sets the force
func (s *ForceSensor) SetForce(force int) {
	s.mu.Lock()
	oldValue := s.force
	if s.force != force {
		s.force = force
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("Force", oldValue, force)
		s.OnPropertyUpdated("Force")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("Force")
	}
}

// IsPressed gets true if the sensor is pressed
func (s *ForceSensor) IsPressed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isPressed
}

// SetIsPressed sets the pressed state
func (s *ForceSensor) SetIsPressed(pressed bool) {
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

// ContinuousMeasurement gets or sets the continuous measurement mode
func (s *ForceSensor) ContinuousMeasurement() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.continuous
}

// SetContinuousMeasurement sets the continuous measurement mode
func (s *ForceSensor) SetContinuousMeasurement(continuous bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.continuous != continuous {
		s.continuous = continuous

		if s.continuous {
			// Start continuous reading
			return s.GetBrick().SelectModeAndRead(s.GetPort(), 0, s.continuous)
		} else {
			// Stop continuous reading
			return s.GetBrick().StopContinuousReadingSensor(s.GetPort())
		}
	}

	return nil
}

// GetForce gets the force by reading from the sensor
func (s *ForceSensor) GetForce() (int, error) {
	trigger := false
	if s.SetupModeAndRead(0, &trigger, !s.ContinuousMeasurement()) {
		return s.Force(), nil
	}

	return 0, errors.New("can't measure the force")
}

// GetPressed gets if the sensor is pressed by reading from the sensor
func (s *ForceSensor) GetPressed() (bool, error) {
	trigger := false
	if s.SetupModeAndRead(0, &trigger, !s.ContinuousMeasurement()) {
		return s.IsPressed(), nil
	}

	return false, errors.New("can't measure if the sensor is pressed")
}

// UpdateFromSensorData updates sensor values from raw sensor data
func (s *ForceSensor) UpdateFromSensorData(data []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update raw values in base sensor
	s.ActiveSensor.UpdateFromSensorData(data)

	// Parse force sensor specific data
	if len(data) >= 1 {
		if force, err := strconv.Atoi(data[0]); err == nil {
			s.force = force
			s.hasForceUpdated = true
			s.isPressed = force > 0
			s.hasIsPressedUpdated = true
		}
	}

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this sensor
func (s *ForceSensor) GetTriggerFlag() *bool {
	return s.ActiveSensor.GetTriggerFlag()
}
