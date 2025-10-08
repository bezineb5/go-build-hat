package buildhat

// DeviceCategory represents the category of a device
type DeviceCategory int

const (
	DeviceCategoryUnknown DeviceCategory = iota
	DeviceCategoryDisconnected
	DeviceCategoryMotor
	DeviceCategorySensor
	DeviceCategoryPassiveMotor
	DeviceCategoryLight
)

// String returns the string representation of the device category
func (d DeviceCategory) String() string {
	switch d {
	case DeviceCategoryDisconnected:
		return "Disconnected"
	case DeviceCategoryMotor:
		return "Motor"
	case DeviceCategorySensor:
		return "Sensor"
	case DeviceCategoryPassiveMotor:
		return "PassiveMotor"
	case DeviceCategoryLight:
		return "Light"
	default:
		return "Unknown"
	}
}

// DeviceSpec contains the specification for a specific device type
type DeviceSpec struct {
	ID       int
	Name     string
	Category DeviceCategory
}

// Known device types from the LEGO Powered Up specification
var deviceRegistry = map[int]DeviceSpec{
	// Passive Motors
	1: {ID: 1, Name: "PassiveMotor", Category: DeviceCategoryPassiveMotor},
	2: {ID: 2, Name: "PassiveMotor", Category: DeviceCategoryPassiveMotor},

	// Lights
	8: {ID: 8, Name: "Light", Category: DeviceCategoryLight},

	// Sensors
	34: {ID: 34, Name: "TiltSensor", Category: DeviceCategorySensor},
	35: {ID: 35, Name: "MotionSensor", Category: DeviceCategorySensor},
	37: {ID: 37, Name: "ColorDistanceSensor", Category: DeviceCategorySensor},
	61: {ID: 61, Name: "ColorSensor", Category: DeviceCategorySensor},
	62: {ID: 62, Name: "DistanceSensor", Category: DeviceCategorySensor},
	63: {ID: 63, Name: "ForceSensor", Category: DeviceCategorySensor},
	64: {ID: 64, Name: "3x3 Color Light Matrix", Category: DeviceCategorySensor},

	// Active Motors
	38: {ID: 38, Name: "Medium Linear Motor", Category: DeviceCategoryMotor},
	46: {ID: 46, Name: "Large Motor", Category: DeviceCategoryMotor},
	47: {ID: 47, Name: "XL Motor", Category: DeviceCategoryMotor},
	48: {ID: 48, Name: "Medium Angular Motor (Cyan)", Category: DeviceCategoryMotor},
	49: {ID: 49, Name: "Large Angular Motor (Cyan)", Category: DeviceCategoryMotor},
	65: {ID: 65, Name: "Small Angular Motor", Category: DeviceCategoryMotor},
	75: {ID: 75, Name: "Medium Angular Motor (Grey)", Category: DeviceCategoryMotor},
	76: {ID: 76, Name: "Large Angular Motor (Grey)", Category: DeviceCategoryMotor},
}

// getDeviceSpec returns the device specification for a type ID
func getDeviceSpec(typeID int) DeviceSpec {
	if typeID == -1 {
		return DeviceSpec{
			ID:       -1,
			Name:     "Disconnected",
			Category: DeviceCategoryDisconnected,
		}
	}

	if spec, exists := deviceRegistry[typeID]; exists {
		return spec
	}

	return DeviceSpec{
		ID:       typeID,
		Name:     "Unknown",
		Category: DeviceCategoryUnknown,
	}
}

// getDeviceName returns the name for a device type ID
func getDeviceName(typeID int) string {
	return getDeviceSpec(typeID).Name
}

// getDeviceType returns the device type for a type ID
func getDeviceType(typeID int) string {
	return getDeviceSpec(typeID).Category.String()
}
