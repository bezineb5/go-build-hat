package models

// LedColor represents the LED colors used with any of the color elements
type LedColor int

const (
	// LedOff LED off
	LedOff LedColor = iota
	// LedBlack black color
	LedBlack
	// LedBrown brown color
	LedBrown
	// LedMagenta magenta color
	LedMagenta
	// LedBlue blue color
	LedBlue
	// LedCyan cyan color
	LedCyan
	// LedPaleGreen pale green color
	LedPaleGreen
	// LedGreen green color
	LedGreen
	// LedYellow yellow color
	LedYellow
	// LedOrange orange color
	LedOrange
	// LedRed red color
	LedRed
	// LedWhite white color
	LedWhite
)

const (
	ledOffName     = "Off"
	ledGreenName   = "Green"
	ledOrangeName  = "Orange"
	ledUnknownName = "Unknown"
)

// String returns the string representation of the LED color
func (lc LedColor) String() string {
	switch lc {
	case LedOff:
		return ledOffName
	case LedBlack:
		return "Black"
	case LedBrown:
		return "Brown"
	case LedMagenta:
		return "Magenta"
	case LedBlue:
		return "Blue"
	case LedCyan:
		return "Cyan"
	case LedPaleGreen:
		return "Pale Green"
	case LedGreen:
		return ledGreenName
	case LedYellow:
		return "Yellow"
	case LedOrange:
		return ledOrangeName
	case LedRed:
		return "Red"
	case LedWhite:
		return "White"
	default:
		return ledUnknownName
	}
}
