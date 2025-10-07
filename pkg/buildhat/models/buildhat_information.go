package models

import "time"

// BuildHatInformation contains the board information
type BuildHatInformation struct {
	// Version gets or sets the version information
	Version string
	// Signature gets or sets the signature of the firmware
	Signature []byte
	// FirmwareDate gets or sets the firmware date
	FirmwareDate time.Time
}

// NewBuildHatInformation creates a new BuildHat information struct
func NewBuildHatInformation(version string, signature []byte, firmwareDate time.Time) *BuildHatInformation {
	return &BuildHatInformation{
		Version:      version,
		Signature:    signature,
		FirmwareDate: firmwareDate,
	}
}
