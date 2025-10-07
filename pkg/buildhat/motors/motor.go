package motors

import (
	"github.com/bezineb5/go-build-hat/pkg/buildhat/sensors"
)

// Motor represents the interface for a motor
type Motor interface {
	sensors.Sensor

	// SetSpeed sets the speed of the motor
	// speed is between -100 and +100
	SetSpeed(speed int) error

	// Stop stops the motor
	Stop() error

	// Start starts the motor
	Start() error

	// StartWithSpeed starts with the specified speed
	// speed is between -100 and +100
	StartWithSpeed(speed int) error

	// GetSpeed gets the speed
	// speed is between -100 and +100
	GetSpeed() int

	// SetBias sets the bias of the motor
	// bias must be between 0 and 1
	SetBias(bias float64) error

	// SetPowerLimit sets the power consumption limit
	// plimit must be between 0 and 1
	SetPowerLimit(plimit float64) error

	// Speed gets the speed of the motor
	// speed is between -100 and +100
	Speed() int

	// GetMotorName gets the name of the motor
	GetMotorName() string

	// Float floats the motor and stops all constraints on it
	Float() error
}
