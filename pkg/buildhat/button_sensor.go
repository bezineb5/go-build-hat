package buildhat

import (
	"fmt"
)

// ButtonSensor creates a button sensor interface for the specified port
func (b *Brick) ButtonSensor(port Port) *ButtonSensor {
	return &ButtonSensor{
		brick: b,
		port:  port.Int(),
	}
}

// ButtonSensor provides a Python-like button sensor interface
type ButtonSensor struct {
	brick *Brick
	port  int
}

// IsPressed checks if the button is pressed
func (s *ButtonSensor) IsPressed() (bool, error) {
	// Set to button mode (mode 0)
	if err := s.brick.writeCommand(Compound(SelectPort(s.port), Select(0))); err != nil {
		return false, err
	}

	// Wait for sensor data
	data, err := s.brick.getSensorData(s.port)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return false, fmt.Errorf("no button data received")
	}

	// Button state is the first value (1 = pressed, 0 = not pressed)
	if state, ok := data[0].(int); ok {
		return state == 1, nil
	}

	return false, fmt.Errorf("invalid button data type")
}
