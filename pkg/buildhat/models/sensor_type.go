package models

// SensorType represents all types of supported sensors
type SensorType byte

const (
	// Passive sensors
	None SensorType = iota
	SystemMediumMotor
	SystemTrainMotor
	SystemTurntableMotor
	GeneralPwm
	ButtonOrTouchSensor
	TechnicLargeMotor
	TechnicXLMotor
	SimpleLights
	FutureLights1
	FutureLights2
	SystemFutureActuator

	// Active sensors (starting from 0x22)
	WeDoTiltSensor                     SensorType = 0x22
	WeDoDistanceSensor                 SensorType = 0x23
	ColourAndDistanceSensor            SensorType = 0x25
	MediumLinearMotor                  SensorType = 0x26
	TechnicLargeMotorId                SensorType = 0x2E
	TechnicXLMotorId                   SensorType = 0x2F
	SpikePrimeMediumMotor              SensorType = 0x30
	SpikePrimeLargeMotor               SensorType = 0x31
	SpikePrimeColorSensor              SensorType = 0x3D
	SpikePrimeUltrasonicDistanceSensor SensorType = 0x3E
	SpikePrimeForceSensor              SensorType = 0x3F
	SpikeEssential3x3ColorLightMatrix  SensorType = 0x40
	SpikeEssentialSmallAngularMotor    SensorType = 0x41
	TechnicMediumAngularMotor          SensorType = 0x4B
	TechnicMotor                       SensorType = 0x4C
)

// IsActiveSensor checks if the sensor is an active one
func (st SensorType) IsActiveSensor() bool {
	return st >= ColourAndDistanceSensor
}

// IsMotor checks if the sensor is a motor
func (st SensorType) IsMotor() bool {
	switch st {
	case MediumLinearMotor,
		SpikeEssentialSmallAngularMotor,
		SpikePrimeLargeMotor,
		SpikePrimeMediumMotor,
		SystemMediumMotor,
		SystemTrainMotor,
		SystemTurntableMotor,
		TechnicLargeMotor,
		TechnicLargeMotorId,
		TechnicXLMotor,
		TechnicXLMotorId,
		TechnicMediumAngularMotor,
		TechnicMotor:
		return true
	default:
		return false
	}
}

// CanSetPowerLevel checks if the sensor can have its power level controlled
func (st SensorType) CanSetPowerLevel() bool {
	switch st {
	case // Motors
		MediumLinearMotor,
		SpikeEssentialSmallAngularMotor,
		SpikePrimeLargeMotor,
		SpikePrimeMediumMotor,
		SystemMediumMotor,
		SystemTrainMotor,
		SystemTurntableMotor,
		TechnicLargeMotor,
		TechnicLargeMotorId,
		TechnicXLMotor,
		TechnicXLMotorId,
		TechnicMediumAngularMotor,
		TechnicMotor,
		// Lights
		SimpleLights,
		FutureLights1,
		FutureLights2:
		return true
	default:
		return false
	}
}

// String returns the string representation of the sensor type
func (st SensorType) String() string {
	switch st {
	case None:
		return "None"
	case SystemMediumMotor:
		return "System Medium Motor"
	case SystemTrainMotor:
		return "System Train Motor"
	case SystemTurntableMotor:
		return "System Turntable Motor"
	case GeneralPwm:
		return "General PWM/Third Party"
	case ButtonOrTouchSensor:
		return "Button/Touch Sensor"
	case TechnicLargeMotor:
		return "Technic Large Motor"
	case TechnicXLMotor:
		return "Technic XL Motor"
	case SimpleLights:
		return "Simple Lights"
	case FutureLights1:
		return "Future Lights 1"
	case FutureLights2:
		return "Future Lights 2"
	case SystemFutureActuator:
		return "System Future Actuator"
	case WeDoTiltSensor:
		return "WeDo Tilt Sensor"
	case WeDoDistanceSensor:
		return "WeDo Distance Sensor"
	case ColourAndDistanceSensor:
		return "Colour and Distance Sensor"
	case MediumLinearMotor:
		return "Medium Linear Motor"
	case TechnicLargeMotorId:
		return "Technic Large Motor (ID)"
	case TechnicXLMotorId:
		return "Technic XL Motor (ID)"
	case SpikePrimeMediumMotor:
		return "SPIKE Prime Medium Motor"
	case SpikePrimeLargeMotor:
		return "SPIKE Prime Large Motor"
	case SpikePrimeColorSensor:
		return "SPIKE Prime Colour Sensor"
	case SpikePrimeUltrasonicDistanceSensor:
		return "SPIKE Prime Ultrasonic Distance Sensor"
	case SpikePrimeForceSensor:
		return "SPIKE Prime Force Sensor"
	case SpikeEssential3x3ColorLightMatrix:
		return "SPIKE Essential 3x3 Colour Light Matrix"
	case SpikeEssentialSmallAngularMotor:
		return "SPIKE Essential Small Angular Motor"
	case TechnicMediumAngularMotor:
		return "Technic Medium Angular Motor"
	case TechnicMotor:
		return "Technic Motor"
	default:
		return "Unknown"
	}
}
