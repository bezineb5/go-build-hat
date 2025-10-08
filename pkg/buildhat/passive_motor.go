package buildhat

import (
	"fmt"
)

// PassiveMotor creates a passive motor interface for the specified port
func (b *Brick) PassiveMotor(port BuildHatPort) *PassiveMotor {
	return &PassiveMotor{
		brick: b,
		port:  port.Int(),
	}
}

// PassiveMotor provides a Python-like passive motor interface (WeDo motors)
type PassiveMotor struct {
	brick *Brick
	port  int
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
	cmd := fmt.Sprintf("port %d ; set %d", m.port, speed)
	return m.brick.writeCommand(cmd)
}

// Stop stops the passive motor
func (m *PassiveMotor) Stop() error {
	stopCmd := fmt.Sprintf("port %d ; set 0", m.port)
	return m.brick.writeCommand(stopCmd)
}

// SetSpeed sets the speed of the passive motor (-100 to 100)
func (m *PassiveMotor) SetSpeed(speed int) error {
	if speed < -100 || speed > 100 {
		return fmt.Errorf("speed must be between -100 and 100, got %d", speed)
	}
	cmd := fmt.Sprintf("port %d ; set %d", m.port, speed)
	return m.brick.writeCommand(cmd)
}
