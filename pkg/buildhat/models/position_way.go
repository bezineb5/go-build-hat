package models

// PositionWay represents the way to go when running a motor to a position
type PositionWay int

const (
	// Shortest shortest way
	Shortest PositionWay = iota
	// Clockwise clockwise way
	Clockwise
	// AntiClockwise anti clockwise way
	AntiClockwise
)

// String returns the string representation of the position way
func (pw PositionWay) String() string {
	switch pw {
	case Shortest:
		return "Shortest"
	case Clockwise:
		return "Clockwise"
	case AntiClockwise:
		return "Anti Clockwise"
	default:
		return ledUnknownName
	}
}
