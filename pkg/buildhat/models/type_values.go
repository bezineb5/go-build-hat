package models

// TypeValues represents the type of values for each mode
type TypeValues int

const (
	// Raw raw values
	Raw TypeValues = iota
	// Percent percent values
	Percent
	// Signal signal values
	Signal
)

// String returns the string representation of the type values
func (tv TypeValues) String() string {
	switch tv {
	case Raw:
		return "Raw"
	case Percent:
		return "Percent"
	case Signal:
		return "Signal"
	default:
		return "Unknown"
	}
}
