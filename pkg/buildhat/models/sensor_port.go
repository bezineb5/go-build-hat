package models

// SensorPort represents the sensor ports 1, 2, 3 and 4
type SensorPort byte

const (
	// PortA represents port A (0)
	PortA SensorPort = iota
	// PortB represents port B (1)
	PortB
	// PortC represents port C (2)
	PortC
	// PortD represents port D (3)
	PortD
)

// String returns the string representation of the sensor port
func (sp SensorPort) String() string {
	switch sp {
	case PortA:
		return "Port A"
	case PortB:
		return "Port B"
	case PortC:
		return "Port C"
	case PortD:
		return "Port D"
	default:
		return "Unknown Port"
	}
}

// Byte returns the byte value of the sensor port
func (sp SensorPort) Byte() byte {
	return byte(sp)
}
