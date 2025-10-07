package models

// MinimumMaximumValues represents minimum and maximum values for a specific mode type
type MinimumMaximumValues struct {
	// TypeValues type of values
	TypeValues TypeValues
	// MinimumValue minimum value
	MinimumValue int
	// MaximumValue maximum value
	MaximumValue int
}
