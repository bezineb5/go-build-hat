package sensors

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// UltrasonicDistanceSensor represents a Spike distance sensor
type UltrasonicDistanceSensor struct {
	*ActiveSensor

	distance           int
	hasDistanceUpdated bool
	continuous         bool

	mu sync.RWMutex
}

// NewUltrasonicDistanceSensor creates a new ultrasonic distance sensor
func NewUltrasonicDistanceSensor(brick BrickInterface, port models.SensorPort) *UltrasonicDistanceSensor {
	sensor := &UltrasonicDistanceSensor{
		ActiveSensor: NewActiveSensor(brick, port, models.SpikePrimeUltrasonicDistanceSensor),
		distance:     0,
		continuous:   false,
	}

	// Initialize sensor - set to distance mode (-1)
	command := fmt.Sprintf("port %d ; plimit 1 ; set -1\r", port.Byte())
	sensor.GetBrick().SendRawCommand(command)

	return sensor
}

// GetSensorName gets the name of the sensor
func (s *UltrasonicDistanceSensor) GetSensorName() string {
	return "SPIKE distance sensor"
}

// Distance gets the distance in millimeters
func (s *UltrasonicDistanceSensor) Distance() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.distance
}

// SetDistance sets the distance
func (s *UltrasonicDistanceSensor) SetDistance(distance int) {
	s.mu.Lock()
	oldValue := s.distance
	if s.distance != distance {
		s.distance = distance
		s.mu.Unlock()

		// Trigger property change events
		s.OnPropertyChanged("Distance", oldValue, distance)
		s.OnPropertyUpdated("Distance")
	} else {
		s.mu.Unlock()
		// Still trigger property updated event even if value didn't change
		s.OnPropertyUpdated("Distance")
	}
}

// ContinuousMeasurement gets or sets the continuous measurement mode
func (s *UltrasonicDistanceSensor) ContinuousMeasurement() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.continuous
}

// SetContinuousMeasurement sets the continuous measurement mode
func (s *UltrasonicDistanceSensor) SetContinuousMeasurement(continuous bool) error {
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

// GetDistance gets the distance by reading from the sensor
func (s *UltrasonicDistanceSensor) GetDistance() (int, error) {
	trigger := false
	if s.SetupModeAndRead(0, &trigger, !s.ContinuousMeasurement()) {
		return s.Distance(), nil
	}

	return 0, errors.New("can't measure the distance")
}

// UpdateFromSensorData updates sensor values from raw sensor data
func (s *UltrasonicDistanceSensor) UpdateFromSensorData(data []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update raw values in base sensor
	s.ActiveSensor.UpdateFromSensorData(data)

	// Parse distance sensor specific data
	if len(data) >= 1 {
		if distance, err := strconv.Atoi(data[0]); err == nil {
			s.distance = distance
			s.hasDistanceUpdated = true
		}
	}

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this sensor
func (s *UltrasonicDistanceSensor) GetTriggerFlag() *bool {
	return s.ActiveSensor.GetTriggerFlag()
}
