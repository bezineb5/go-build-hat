package buildhat

import (
	"fmt"
)

// Light creates a light interface for the specified port
func (b *Brick) Light(port string) *Light {
	portNum := int(port[0] - 'A')
	return &Light{
		brick: b,
		port:  portNum,
	}
}

// Light provides a Python-like light interface
type Light struct {
	brick *Brick
	port  int
}

// SetBrightness sets the brightness of the light (0-100)
func (l *Light) SetBrightness(brightness int) error {
	if brightness < 0 || brightness > 100 {
		return fmt.Errorf("brightness must be between 0 and 100")
	}

	if brightness == 0 {
		return l.Off()
	}

	// Convert brightness to 0.0-1.0 range
	value := float64(brightness) / 100.0
	cmd := fmt.Sprintf("port %d ; on ; set %.2f", l.port, value)
	return l.brick.writeCommand(cmd)
}

// On turns the light on at full brightness
func (l *Light) On() error {
	return l.SetBrightness(100)
}

// Off turns the light off
func (l *Light) Off() error {
	// Using coast to turn off lights completely
	cmd := fmt.Sprintf("port %d ; coast", l.port)
	return l.brick.writeCommand(cmd)
}

// GetBrightness gets the current brightness reading (not supported on all lights)
// This returns the stored data value, which may not be available for all light types
func (l *Light) GetBrightness() (int, error) {
	l.brick.mu.RLock()
	defer l.brick.mu.RUnlock()

	if len(l.brick.connections[l.port].Data) > 0 {
		if brightness, ok := l.brick.connections[l.port].Data[0].(int); ok {
			return brightness, nil
		}
	}

	// Default to a middle value if no data available
	return 50, nil
}
