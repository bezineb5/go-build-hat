package models

// LedMode represents the LED mode for the LEDs on the Build HAT
type LedMode int

const (
	// VoltageDependant LEDs lit depend on the voltage on the input power jack (default)
	VoltageDependant LedMode = -1
	// Off LEDs off
	Off LedMode = 0
	// Orange Orange LEDs
	Orange LedMode = 1
	// Green Green LEDs
	Green LedMode = 2
	// Both Orange and green together
	Both LedMode = 3
)

// String returns the string representation of the LED mode
func (lm LedMode) String() string {
	switch lm {
	case VoltageDependant:
		return "Voltage Dependent"
	case Off:
		return "Off"
	case Orange:
		return "Orange"
	case Green:
		return "Green"
	case Both:
		return "Both"
	default:
		return ledUnknownName
	}
}
