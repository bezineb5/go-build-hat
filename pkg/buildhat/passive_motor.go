package buildhat

import (
	"fmt"
)

// PassiveMotor creates a passive motor interface for the specified port
func (b *Brick) PassiveMotor(port Port) *PassiveMotor {
	return &PassiveMotor{
		brick: b,
		port:  port,
	}
}

// PassiveMotor provides a Python-like passive motor interface (WeDo motors)
type PassiveMotor struct {
	brick *Brick
	port  Port
}

// Start starts the passive motor at the specified speed (-100 to 100)
func (m *PassiveMotor) Start(speed int) error {
	if speed == 0 {
		speed = 50 // Default speed
	}
	if speed < -100 || speed > 100 {
		return fmt.Errorf("speed must be between -100 and 100, got %d", speed)
	}

	// Passive motors use simple speed control
	return m.brick.writeCommand(Compound(SelectPort(m.port), SetConstant(float64(speed))))
}

// Stop stops the passive motor
func (m *PassiveMotor) Stop() error {
	return m.brick.writeCommand(Compound(SelectPort(m.port), SetConstant(0)))
}

// SetSpeed sets the speed of the passive motor (-100 to 100)
func (m *PassiveMotor) SetSpeed(speed int) error {
	if speed < -100 || speed > 100 {
		return fmt.Errorf("speed must be between -100 and 100, got %d", speed)
	}
	return m.brick.writeCommand(Compound(SelectPort(m.port), SetConstant(float64(speed))))
}
