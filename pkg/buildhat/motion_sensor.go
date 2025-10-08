package buildhat

import (
	"fmt"
)

// MotionSensor creates a motion sensor interface for the specified port
func (b *Brick) MotionSensor(port string) *MotionSensor {
	portNum := int(port[0] - 'A')
	return &MotionSensor{
		brick: b,
		port:  portNum,
	}
}

// MotionSensor provides a Python-like motion sensor interface (WeDo sensor)
type MotionSensor struct {
	brick *Brick
	port  int
}

// GetDistance gets the distance reading from the motion sensor
func (s *MotionSensor) GetDistance() (int, error) {
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

// GetMovementCount gets the movement count (number of detected motions)
func (s *MotionSensor) GetMovementCount() (int, error) {
	// Set to movement count mode (mode 1)
	cmd := fmt.Sprintf("port %d ; select 1", s.port)
	if err := s.brick.writeCommand(cmd); err != nil {
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
