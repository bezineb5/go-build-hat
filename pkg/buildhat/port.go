package buildhat

import "fmt"

// BuildHatPort represents a physical port on the BuildHat (A, B, C, or D)
type BuildHatPort int

const (
	// NumPorts is the number of ports on the BuildHat
	NumPorts = 4

	PortA BuildHatPort = 0
	PortB BuildHatPort = 1
	PortC BuildHatPort = 2
	PortD BuildHatPort = 3
)

// String returns the string representation of the port ("A", "B", "C", or "D")
func (p BuildHatPort) String() string {
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
func (p BuildHatPort) Int() int {
	return int(p)
}

// IsValid returns true if the port is valid (A, B, C, or D)
func (p BuildHatPort) IsValid() bool {
	return p >= PortA && p <= PortD
}

// ParsePort converts a string ("A", "B", "C", or "D") to a BuildHatPort
func ParsePort(s string) (BuildHatPort, error) {
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

// AllPorts returns a slice of all valid ports
func AllPorts() []BuildHatPort {
	return []BuildHatPort{PortA, PortB, PortC, PortD}
}
