package sensors

import (
	"errors"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// PassiveLight represents a simple passive light sensor
type PassiveLight struct {
	*BaseSensor

	isOn        bool
	brightness  int
	triggerFlag bool

	mu sync.RWMutex
}

// NewPassiveLight creates a new passive light sensor
func NewPassiveLight(brick BrickInterface, port models.SensorPort) *PassiveLight {
	return &PassiveLight{
		BaseSensor: NewBaseSensor(brick, port, models.SimpleLights),
		isOn:       false,
		brightness: 0,
	}
}

// GetSensorName gets the name of the sensor
func (s *PassiveLight) GetSensorName() string {
	return "Passive light"
}

// Brightness gets the brightness from 0 to 100
func (s *PassiveLight) Brightness() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.brightness
}

// SetBrightness sets the brightness from 0 to 100
func (s *PassiveLight) SetBrightness(brightness int) error {
	if brightness < 0 || brightness > 100 {
		return errors.New("brightness can only be between 0 (off) and 100 (full bright)")
	}

	s.mu.Lock()
	s.brightness = brightness
	s.mu.Unlock()

	if s.isOn {
		// Update the light with new brightness using power level control
		return s.GetBrick().SetPowerLevel(s.GetPort(), brightness)
	}

	return nil
}

// IsOn gets true if the light is on
func (s *PassiveLight) IsOn() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isOn
}

// SetIsOn sets the light on/off state
func (s *PassiveLight) SetIsOn(on bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isOn != on {
		s.isOn = on

		if s.isOn {
			// Turn on the light with current brightness
			return s.GetBrick().SetPowerLevel(s.GetPort(), s.brightness)
		} else {
			// Turn off the light by setting brightness to 0
			return s.GetBrick().SetPowerLevel(s.GetPort(), 0)
		}
	}

	return nil
}

// TurnOn turns the light on
func (s *PassiveLight) TurnOn() error {
	return s.SetIsOn(true)
}

// TurnOff turns the light off
func (s *PassiveLight) TurnOff() error {
	return s.SetIsOn(false)
}

// On turns the light on with the current brightness
func (s *PassiveLight) On() error {
	return s.SetIsOn(true)
}

// OnWithBrightness turns the light on with a specific brightness
func (s *PassiveLight) OnWithBrightness(brightness int) error {
	if err := s.SetBrightness(brightness); err != nil {
		return err
	}
	return s.SetIsOn(true)
}

// Off turns the light off (sets brightness to 0)
func (s *PassiveLight) Off() error {
	s.mu.Lock()
	s.brightness = 0
	s.isOn = false
	s.mu.Unlock()

	return s.GetBrick().SetPowerLevel(s.GetPort(), 0)
}

// UpdateFromSensorData updates sensor values from raw sensor data
func (s *PassiveLight) UpdateFromSensorData(data []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Passive lights don't typically receive data updates
	// This is mainly for interface compliance
	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this sensor
func (s *PassiveLight) GetTriggerFlag() *bool {
	return &s.triggerFlag
}
