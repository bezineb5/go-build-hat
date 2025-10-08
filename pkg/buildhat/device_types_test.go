package buildhat

import "testing"

func TestDeviceCategory_String(t *testing.T) {
	tests := []struct {
		category DeviceCategory
		expected string
	}{
		{DeviceCategoryUnknown, "Unknown"},
		{DeviceCategoryDisconnected, "Disconnected"},
		{DeviceCategoryMotor, "Motor"},
		{DeviceCategorySensor, "Sensor"},
		{DeviceCategoryPassiveMotor, "PassiveMotor"},
		{DeviceCategoryLight, "Light"},
		{DeviceCategory(999), "Unknown"}, // Invalid category
	}

	for _, tt := range tests {
		result := tt.category.String()
		if result != tt.expected {
			t.Errorf("DeviceCategory(%d).String() = %q, want %q", tt.category, result, tt.expected)
		}
	}
}

func TestGetDeviceSpec_Disconnected(t *testing.T) {
	spec := getDeviceSpec(-1)

	if spec.ID != -1 {
		t.Errorf("Expected ID -1, got %d", spec.ID)
	}
	if spec.Name != "Disconnected" {
		t.Errorf("Expected name 'Disconnected', got %q", spec.Name)
	}
	if spec.Category != DeviceCategoryDisconnected {
		t.Errorf("Expected category DeviceCategoryDisconnected, got %v", spec.Category)
	}
}

func TestGetDeviceSpec_Unknown(t *testing.T) {
	spec := getDeviceSpec(999) // Unknown device ID

	if spec.ID != 999 {
		t.Errorf("Expected ID 999, got %d", spec.ID)
	}
	if spec.Name != "Unknown" {
		t.Errorf("Expected name 'Unknown', got %q", spec.Name)
	}
	if spec.Category != DeviceCategoryUnknown {
		t.Errorf("Expected category DeviceCategoryUnknown, got %v", spec.Category)
	}
}

func TestGetDeviceSpec_PassiveMotors(t *testing.T) {
	tests := []struct {
		id   int
		name string
	}{
		{1, "PassiveMotor"},
		{2, "PassiveMotor"},
	}

	for _, tt := range tests {
		spec := getDeviceSpec(tt.id)
		if spec.ID != tt.id {
			t.Errorf("Expected ID %d, got %d", tt.id, spec.ID)
		}
		if spec.Name != tt.name {
			t.Errorf("Expected name %q, got %q", tt.name, spec.Name)
		}
		if spec.Category != DeviceCategoryPassiveMotor {
			t.Errorf("Expected category DeviceCategoryPassiveMotor, got %v", spec.Category)
		}
	}
}

func TestGetDeviceSpec_Light(t *testing.T) {
	spec := getDeviceSpec(8)

	if spec.ID != 8 {
		t.Errorf("Expected ID 8, got %d", spec.ID)
	}
	if spec.Name != "Light" {
		t.Errorf("Expected name 'Light', got %q", spec.Name)
	}
	if spec.Category != DeviceCategoryLight {
		t.Errorf("Expected category DeviceCategoryLight, got %v", spec.Category)
	}
}

func TestGetDeviceSpec_Sensors(t *testing.T) {
	tests := []struct {
		id   int
		name string
	}{
		{34, "TiltSensor"},
		{35, "MotionSensor"},
		{37, "ColorDistanceSensor"},
		{61, "ColorSensor"},
		{62, "DistanceSensor"},
		{63, "ForceSensor"},
		{64, "3x3 Color Light Matrix"},
	}

	for _, tt := range tests {
		spec := getDeviceSpec(tt.id)
		if spec.ID != tt.id {
			t.Errorf("ID %d: Expected ID %d, got %d", tt.id, tt.id, spec.ID)
		}
		if spec.Name != tt.name {
			t.Errorf("ID %d: Expected name %q, got %q", tt.id, tt.name, spec.Name)
		}
		if spec.Category != DeviceCategorySensor {
			t.Errorf("ID %d: Expected category DeviceCategorySensor, got %v", tt.id, spec.Category)
		}
	}
}

func TestGetDeviceSpec_ActiveMotors(t *testing.T) {
	tests := []struct {
		id   int
		name string
	}{
		{38, "Medium Linear Motor"},
		{46, "Large Motor"},
		{47, "XL Motor"},
		{48, "Medium Angular Motor (Cyan)"},
		{49, "Large Angular Motor (Cyan)"},
		{65, "Small Angular Motor"},
		{75, "Medium Angular Motor (Grey)"},
		{76, "Large Angular Motor (Grey)"},
	}

	for _, tt := range tests {
		spec := getDeviceSpec(tt.id)
		if spec.ID != tt.id {
			t.Errorf("ID %d: Expected ID %d, got %d", tt.id, tt.id, spec.ID)
		}
		if spec.Name != tt.name {
			t.Errorf("ID %d: Expected name %q, got %q", tt.id, tt.name, spec.Name)
		}
		if spec.Category != DeviceCategoryMotor {
			t.Errorf("ID %d: Expected category DeviceCategoryMotor, got %v", tt.id, spec.Category)
		}
	}
}

func TestGetDeviceName(t *testing.T) {
	tests := []struct {
		id   int
		name string
	}{
		{-1, "Disconnected"},
		{1, "PassiveMotor"},
		{8, "Light"},
		{34, "TiltSensor"},
		{38, "Medium Linear Motor"},
		{61, "ColorSensor"},
		{999, "Unknown"},
	}

	for _, tt := range tests {
		result := getDeviceName(tt.id)
		if result != tt.name {
			t.Errorf("getDeviceName(%d) = %q, want %q", tt.id, result, tt.name)
		}
	}
}

func TestGetDeviceCategory(t *testing.T) {
	tests := []struct {
		id       int
		category DeviceCategory
	}{
		{-1, DeviceCategoryDisconnected},
		{1, DeviceCategoryPassiveMotor},
		{2, DeviceCategoryPassiveMotor},
		{8, DeviceCategoryLight},
		{34, DeviceCategorySensor},
		{35, DeviceCategorySensor},
		{37, DeviceCategorySensor},
		{38, DeviceCategoryMotor},
		{46, DeviceCategoryMotor},
		{61, DeviceCategorySensor},
		{999, DeviceCategoryUnknown},
	}

	for _, tt := range tests {
		result := getDeviceCategory(tt.id)
		if result != tt.category {
			t.Errorf("getDeviceCategory(%d) = %v, want %v", tt.id, result, tt.category)
		}
	}
}

func TestDeviceRegistry_AllDevicesHaveValidData(t *testing.T) {
	// Verify that all devices in the registry have valid data
	for id, spec := range deviceRegistry {
		if spec.ID != id {
			t.Errorf("Device %d: registry key (%d) doesn't match spec.ID (%d)", id, id, spec.ID)
		}
		if spec.Name == "" {
			t.Errorf("Device %d: has empty name", id)
		}
		if spec.Category == DeviceCategoryUnknown || spec.Category == DeviceCategoryDisconnected {
			t.Errorf("Device %d: has invalid category %v", id, spec.Category)
		}
	}
}

func TestDeviceRegistry_ExpectedCount(t *testing.T) {
	// We expect 19 devices in the registry (2 passive motors + 1 light + 7 sensors + 8 motors)
	expectedCount := 18
	actualCount := len(deviceRegistry)

	if actualCount != expectedCount {
		t.Errorf("Expected %d devices in registry, got %d", expectedCount, actualCount)
	}
}

func TestDeviceRegistry_CategoryCounts(t *testing.T) {
	categoryCounts := make(map[DeviceCategory]int)

	for _, spec := range deviceRegistry {
		categoryCounts[spec.Category]++
	}

	expectedCounts := map[DeviceCategory]int{
		DeviceCategoryPassiveMotor: 2,
		DeviceCategoryLight:        1,
		DeviceCategorySensor:       7,
		DeviceCategoryMotor:        8,
	}

	for category, expectedCount := range expectedCounts {
		actualCount := categoryCounts[category]
		if actualCount != expectedCount {
			t.Errorf("Expected %d devices in category %v, got %d", expectedCount, category, actualCount)
		}
	}
}

func TestDeviceSpec_KnownMotorIDs(t *testing.T) {
	// Test that all known motor IDs are correctly categorized
	motorIDs := []int{38, 46, 47, 48, 49, 65, 75, 76}

	for _, id := range motorIDs {
		spec := getDeviceSpec(id)
		if spec.Category != DeviceCategoryMotor {
			t.Errorf("Motor ID %d has wrong category: %v", id, spec.Category)
		}
	}
}

func TestDeviceSpec_KnownSensorIDs(t *testing.T) {
	// Test that all known sensor IDs are correctly categorized
	sensorIDs := []int{34, 35, 37, 61, 62, 63, 64}

	for _, id := range sensorIDs {
		spec := getDeviceSpec(id)
		if spec.Category != DeviceCategorySensor {
			t.Errorf("Sensor ID %d has wrong category: %v", id, spec.Category)
		}
	}
}

func TestDeviceSpec_Consistency(t *testing.T) {
	// Test that getDeviceName and getDeviceCategory are consistent with getDeviceSpec
	testIDs := []int{-1, 1, 8, 34, 38, 61, 999}

	for _, id := range testIDs {
		spec := getDeviceSpec(id)
		name := getDeviceName(id)
		category := getDeviceCategory(id)

		if name != spec.Name {
			t.Errorf("ID %d: getDeviceName() = %q, but spec.Name = %q", id, name, spec.Name)
		}
		if category != spec.Category {
			t.Errorf("ID %d: getDeviceCategory() = %v, but spec.Category = %v", id, category, spec.Category)
		}
	}
}
