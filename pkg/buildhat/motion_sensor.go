package buildhat

import (
	"fmt"
)

// MotionSensor creates a motion sensor interface for the specified port
func (b *Brick) MotionSensor(port Port) *MotionSensor {
	return &MotionSensor{
		brick: b,
		port:  port,
	}
}

// MotionSensor provides a Python-like motion sensor interface (WeDo sensor)
type MotionSensor struct {
	brick *Brick
	port  Port
}

// GetDistance gets the distance reading from the motion sensor
func (s *MotionSensor) GetDistance() (int, error) {
	// Set to distance mode (mode 0)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(0))); err != nil {
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

// GetMovementCount gets the movement count (number of detected motions)
func (s *MotionSensor) GetMovementCount() (int, error) {
	// Set to movement count mode (mode 1)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(1))); err != nil {
		return 0, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("no movement count data received")
	}

	// Movement count is the first value
	if count, ok := data[0].(int); ok {
		return count, nil
	}

	return 0, fmt.Errorf("invalid movement count data type")
}
