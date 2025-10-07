package models

import "reflect"

// ModeDetail represents mode details
type ModeDetail struct {
	// Number gets the mode number
	Number int
	// Name gets the name of the mode
	Name string
	// Unit gets the unit of the mode
	Unit string
	// NumberOfDataItems gets the number of data items
	NumberOfDataItems int
	// DataType gets the data type
	DataType reflect.Type
	// NumberOfCharsToDisplay gets the number of chars to display the value
	NumberOfCharsToDisplay int
	// NumberOfData gets the number of data in the mode
	NumberOfData int
	// DecimalPrecision gets the decimal precision (for float, 0 otherwise)
	DecimalPrecision int
	// MinimumMaximumValues gets the minimum and maximum values for the mode
	MinimumMaximumValues []MinimumMaximumValues
}
