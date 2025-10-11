package buildhat

import "fmt"

// Port represents a physical port on the BuildHat (A, B, C, or D)
type Port int

const (
	// NumPorts is the number of ports on the BuildHat
	NumPorts = 4

	PortA Port = 0
	PortB Port = 1
	PortC Port = 2
	PortD Port = 3
)

// String returns the string representation of the port ("A", "B", "C", or "D")
func (p Port) String() string {
	switch p {
	case PortA:
		return "A"
	case PortB:
		return "B"
	case PortC:
		return "C"
	case PortD:
		return "D"
	default:
		return fmt.Sprintf("Invalid(%d)", p)
	}
}

// Int returns the integer value of the port (0-3)
func (p Port) Int() int {
	return int(p)
}

// IsValid returns true if the port is valid (A, B, C, or D)
func (p Port) IsValid() bool {
	return p >= PortA && p <= PortD
}

// ParsePort converts a string ("A", "B", "C", or "D") to a BuildHatPort
func ParsePort(s string) (Port, error) {
	switch s {
	case "A":
		return PortA, nil
	case "B":
		return PortB, nil
	case "C":
		return PortC, nil
	case "D":
		return PortD, nil
	default:
		return -1, fmt.Errorf("invalid port: %s (must be A, B, C, or D)", s)
	}
}

// ParsePortNumber converts a rune ('0', '1', '2', or '3') to a BuildHatPort
func ParsePortNumber(r rune) (Port, error) {
	port := Port(r - '0')
	if port < PortA || port > PortD {
		return -1, fmt.Errorf("invalid port number: %c (must be 0-3)", r)
	}
	return port, nil
}

// AllPorts returns a slice of all valid ports
func AllPorts() []Port {
	return []Port{PortA, PortB, PortC, PortD}
}
