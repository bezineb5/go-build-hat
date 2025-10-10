package buildhat

import (
	"fmt"
)

// ForceSensor creates a force sensor interface for the specified port
func (b *Brick) ForceSensor(port Port) *ForceSensor {
	return &ForceSensor{
		brick: b,
		port:  port.Int(),
	}
}

// ForceSensor provides a Python-like force sensor interface
type ForceSensor struct {
	brick *Brick
	port  int
}

// GetForce gets the current force reading in Newtons
func (s *ForceSensor) GetForce() (int, error) {
	// Python uses combi mode: [(0, 0), (1, 0), (3, 0)]
	// For simplicity, we'll just use mode 0
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(0))); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no force data received")
	}

	// Force is the first value
	if force, ok := data[0].(int); ok {
		return force, nil
	}

	return 0, fmt.Errorf("invalid force data type")
}

// IsPressed checks if the force sensor is pressed
func (s *ForceSensor) IsPressed() (bool, error) {
	force, err := s.GetForce()
	if err != nil {
		return false, err
	}

	// Consider pressed if force > 0
	return force > 0, nil
}
