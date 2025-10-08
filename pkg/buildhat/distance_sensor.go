package buildhat

import (
	"fmt"
)

// DistanceSensor creates a distance sensor interface for the specified port
func (b *Brick) DistanceSensor(port BuildHatPort) *DistanceSensor {
	return &DistanceSensor{
		brick: b,
		port:  port.Int(),
	}
}

// DistanceSensor provides a Python-like distance sensor interface
type DistanceSensor struct {
	brick *Brick
	port  int
}

// GetDistance gets the current distance reading in millimeters
func (s *DistanceSensor) GetDistance() (int, error) {
	// Set to distance mode (mode 0)
	cmd := fmt.Sprintf("port %d ; select 0", s.port)
	if err := s.brick.writeCommand(cmd); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no distance data received")
	}

	// Distance is the first value
	if distance, ok := data[0].(int); ok {
		return distance, nil
	}

	return 0, fmt.Errorf("invalid distance data type")
}
