package buildhat

import (
	"fmt"
)

// PassiveMotor creates a passive motor interface for the specified port
func (b *Brick) PassiveMotor(port string) *PassiveMotor {
	portNum := int(port[0] - 'A')
	return &PassiveMotor{
		brick: b,
		port:  portNum,
	}
}

// PassiveMotor provides a Python-like passive motor interface (WeDo motors)
type PassiveMotor struct {
	brick *Brick
	port  int
}

// Start starts the passive motor at the specified speed
func (m *PassiveMotor) Start(speed int) error {
	if speed == 0 {
		speed = 50 // Default speed
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

// SetSpeed sets the speed of the passive motor
func (m *PassiveMotor) SetSpeed(speed int) error {
	cmd := fmt.Sprintf("port %d ; set %d", m.port, speed)
	return m.brick.writeCommand(cmd)
}
