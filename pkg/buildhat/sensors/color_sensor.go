package sensors

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// ColorSensor represents a color sensor
type ColorSensor struct {
	*ActiveSensor

	color               color.RGBA
	isColorDetected     bool
	hasColorUpdated     bool
	reflectedLight      int
	hasReflectedUpdated bool
	ambientLight        int
	hasAmbientUpdated   bool

	mu sync.RWMutex
}

// NewColorSensor creates a new color sensor
func NewColorSensor(brick BrickInterface, port models.SensorPort, sensorType models.SensorType) (*ColorSensor, error) {
	cs := &ColorSensor{
		ActiveSensor:    NewActiveSensor(brick, port, sensorType),
		color:           color.RGBA{R: 0, G: 0, B: 0, A: 255},
		isColorDetected: false,
		reflectedLight:  0,
		ambientLight:    0,
	}

	// Initialize sensor based on type
	switch sensorType {
	case models.SpikePrimeColorSensor:
		// Send initialization command for Spike Prime color sensor
		// This sets the sensor to color mode (-1)
		command := fmt.Sprintf("port %d ; plimit 1 ; set -1\r", port.Byte())
		if err := cs.GetBrick().SendRawCommand(command); err != nil {
			return nil, err
		}
	case models.ColourAndDistanceSensor:
		// Switch sensor on for color and distance sensor
		if err := cs.SwitchOn(); err != nil {
			return nil, err
		}
	}

	return cs, nil
}

const colorSensorName = "Color sensor"

// GetSensorName gets the name of the sensor
func (s *ColorSensor) GetSensorName() string {
	return colorSensorName
}

// Color gets the last measured color
func (s *ColorSensor) Color() color.RGBA {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.color
}

// SetColor sets the color
func (s *ColorSensor) SetColor(c color.RGBA) {
	s.mu.Lock()
	oldValue := s.color
	if s.color != c {
		s.color = c
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("Color", oldValue, c)
		s.OnPropertyUpdated("Color")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("Color")
	}
}

// IsColorDetected gets true if a color is detected
func (s *ColorSensor) IsColorDetected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isColorDetected
}

// SetIsColorDetected sets the color detected state
func (s *ColorSensor) SetIsColorDetected(detected bool) {
	s.mu.Lock()
	oldValue := s.isColorDetected
	if s.isColorDetected != detected {
		s.isColorDetected = detected
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("IsColorDetected", oldValue, detected)
		s.OnPropertyUpdated("IsColorDetected")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("IsColorDetected")
	}
}

// ReflectedLight gets the reflected light
func (s *ColorSensor) ReflectedLight() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.reflectedLight
}

// SetReflectedLight sets the reflected light
func (s *ColorSensor) SetReflectedLight(light int) {
	s.mu.Lock()
	oldValue := s.reflectedLight
	if s.reflectedLight != light {
		s.reflectedLight = light
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("ReflectedLight", oldValue, light)
		s.OnPropertyUpdated("ReflectedLight")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("ReflectedLight")
	}
}

// AmbientLight gets the ambient light
func (s *ColorSensor) AmbientLight() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ambientLight
}

// SetAmbientLight sets the ambient light
func (s *ColorSensor) SetAmbientLight(light int) {
	s.mu.Lock()
	oldValue := s.ambientLight
	if s.ambientLight != light {
		s.ambientLight = light
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("AmbientLight", oldValue, light)
		s.OnPropertyUpdated("AmbientLight")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("AmbientLight")
	}
}

// GetColor gets the color by reading from the sensor
func (s *ColorSensor) GetColor() (color.RGBA, error) {
	trigger := false
	if s.SetupModeAndRead(6, &trigger, true) {
		return s.Color(), nil
	}

	return color.RGBA{}, errors.New("can't measure the color")
}

// GetReflectedLight gets the reflected light by reading from the sensor
func (s *ColorSensor) GetReflectedLight() (int, error) {
	trigger := false
	if s.SetupModeAndRead(3, &trigger, true) {
		return s.ReflectedLight(), nil
	}

	return 0, errors.New("can't measure the reflected light")
}

// GetAmbientLight gets the ambient light by reading from the sensor
func (s *ColorSensor) GetAmbientLight() (int, error) {
	trigger := false
	if s.SetupModeAndRead(4, &trigger, true) {
		return s.AmbientLight(), nil
	}

	return 0, errors.New("can't measure the ambient light")
}

// UpdateFromSensorData updates sensor values from raw sensor data
func (s *ColorSensor) UpdateFromSensorData(data []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update raw values in base sensor
	if err := s.ActiveSensor.UpdateFromSensorData(data); err != nil {
		return err
	}

	// Parse color sensor specific data
	if len(data) >= 3 {
		// Parse color values (R, G, B)
		if r, err := strconv.Atoi(data[0]); err == nil {
			if g, err := strconv.Atoi(data[1]); err == nil {
				if b, err := strconv.Atoi(data[2]); err == nil {
					s.color = color.RGBA{
						R: clamp8(r),
						G: clamp8(g),
						B: clamp8(b),
						A: 255,
					}
					s.isColorDetected = true
					s.hasColorUpdated = true
				}
			}
		}
	}

	if len(data) >= 4 {
		// Parse reflected light
		if reflected, err := strconv.Atoi(data[3]); err == nil {
			s.reflectedLight = reflected
			s.hasReflectedUpdated = true
		}
	}

	if len(data) >= 5 {
		// Parse ambient light
		if ambient, err := strconv.Atoi(data[4]); err == nil {
			s.ambientLight = ambient
			s.hasAmbientUpdated = true
		}
	}

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this sensor
func (s *ColorSensor) GetTriggerFlag() *bool {
	return s.ActiveSensor.GetTriggerFlag()
}

// Helper function to clamp values
func clamp8(value int) uint8 {
	const minVal = 0
	const maxVal = 255

	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return uint8(value)
}
