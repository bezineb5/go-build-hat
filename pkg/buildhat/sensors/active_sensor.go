package sensors

import (
	"sync"
	"time"

	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

// ActiveSensor represents an active sensor that can read data continuously
type ActiveSensor struct {
	*BaseSensor

	// Internal fields for sensor data
	valuesAsString          []string
	hasValueAsStringUpdated bool
	baudRate                int
	hardwareVersion         uint32
	softwareVersion         uint32
	modeDetails             []models.ModeDetail
	combiModes              []models.CombiModes
	combiReadingModes       []int
	speedPid                models.RecommendedPid
	positionPid             models.RecommendedPid

	// Trigger mechanism for data updates
	triggerFlag bool

	// Thread safety
	mu sync.RWMutex
}

// NewActiveSensor creates a new active sensor
func NewActiveSensor(brick BrickInterface, port models.SensorPort, sensorType models.SensorType) *ActiveSensor {
	return &ActiveSensor{
		BaseSensor:        NewBaseSensor(brick, port, sensorType),
		modeDetails:       make([]models.ModeDetail, 0),
		combiModes:        make([]models.CombiModes, 0),
		combiReadingModes: make([]int, 0),
	}
}

// GetSensorName gets the name of the sensor
func (s *ActiveSensor) GetSensorName() string {
	return "Active sensor"
}

// ValuesAsString gets the raw values as strings
func (s *ActiveSensor) ValuesAsString() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.valuesAsString
}

// SetValuesAsString sets the raw values as strings
func (s *ActiveSensor) SetValuesAsString(values []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if values have changed
	changed := false
	if len(s.valuesAsString) != len(values) {
		changed = true
	} else {
		for i, v := range values {
			if i >= len(s.valuesAsString) || s.valuesAsString[i] != v {
				changed = true
				break
			}
		}
	}

	if changed {
		s.valuesAsString = make([]string, len(values))
		copy(s.valuesAsString, values)
	}

	s.hasValueAsStringUpdated = true
}

// HasValueAsStringUpdated gets whether values have been updated
func (s *ActiveSensor) HasValueAsStringUpdated() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hasValueAsStringUpdated
}

// BaudRate gets the baud rate
func (s *ActiveSensor) BaudRate() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.baudRate
}

// SetBaudRate sets the baud rate
func (s *ActiveSensor) SetBaudRate(baudRate int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.baudRate = baudRate
}

// HardwareVersion gets the hardware version
func (s *ActiveSensor) HardwareVersion() uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hardwareVersion
}

// SetHardwareVersion sets the hardware version
func (s *ActiveSensor) SetHardwareVersion(version uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hardwareVersion = version
}

// SoftwareVersion gets the software version
func (s *ActiveSensor) SoftwareVersion() uint32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.softwareVersion
}

// SetSoftwareVersion sets the software version
func (s *ActiveSensor) SetSoftwareVersion(version uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.softwareVersion = version
}

// CombiModes gets the combi modes
func (s *ActiveSensor) CombiModes() []models.CombiModes {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.CombiModes, len(s.combiModes))
	copy(result, s.combiModes)
	return result
}

// SetCombiModes sets the combi modes
func (s *ActiveSensor) SetCombiModes(modes []models.CombiModes) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.combiModes = make([]models.CombiModes, len(modes))
	copy(s.combiModes, modes)
}

// ModeDetails gets the mode details
func (s *ActiveSensor) ModeDetails() []models.ModeDetail {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.ModeDetail, len(s.modeDetails))
	copy(result, s.modeDetails)
	return result
}

// SetModeDetails sets the mode details
func (s *ActiveSensor) SetModeDetails(details []models.ModeDetail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modeDetails = make([]models.ModeDetail, len(details))
	copy(s.modeDetails, details)
}

// NumberOfModes gets the number of modes
func (s *ActiveSensor) NumberOfModes() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.modeDetails)
}

// SpeedPid gets the speed PID settings
func (s *ActiveSensor) SpeedPid() models.RecommendedPid {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.speedPid
}

// SetSpeedPid sets the speed PID settings
func (s *ActiveSensor) SetSpeedPid(pid models.RecommendedPid) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.speedPid = pid
}

// PositionPid gets the position PID settings
func (s *ActiveSensor) PositionPid() models.RecommendedPid {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.positionPid
}

// SetPositionPid sets the position PID settings
func (s *ActiveSensor) SetPositionPid(pid models.RecommendedPid) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.positionPid = pid
}

// SelectCombiModesAndRead selects combi modes and reads data
func (s *ActiveSensor) SelectCombiModesAndRead(modes []int, readOnce bool) error {
	return s.GetBrick().SelectCombiModesAndRead(s.GetPort(), modes, readOnce)
}

// SelectModeAndRead selects a mode and reads data
func (s *ActiveSensor) SelectModeAndRead(mode int, readOnce bool) error {
	return s.GetBrick().SelectModeAndRead(s.GetPort(), mode, readOnce)
}

// StopReading stops continuous reading
func (s *ActiveSensor) StopReading() error {
	return s.GetBrick().StopContinuousReadingSensor(s.GetPort())
}

// SwitchOn switches the sensor on
func (s *ActiveSensor) SwitchOn() error {
	return s.GetBrick().SwitchSensorOn(s.GetPort())
}

// SwitchOff switches the sensor off
func (s *ActiveSensor) SwitchOff() error {
	return s.GetBrick().SwitchSensorOff(s.GetPort())
}

// WriteBytes writes bytes directly to the sensor
func (s *ActiveSensor) WriteBytes(data []byte, singleHeader bool) error {
	return s.GetBrick().WriteBytesToSensor(s.GetPort(), data, singleHeader)
}

// SetupModeAndRead sets up mode and reads data with timeout
func (s *ActiveSensor) SetupModeAndRead(mode int, trigger *bool, readOnce bool) bool {
	const timeoutMeasuresSeconds = 4

	*trigger = false
	timeout := time.Now().Add(time.Duration(timeoutMeasuresSeconds) * time.Second)

	err := s.SelectModeAndRead(mode, readOnce)
	if err != nil {
		return false
	}

	for !*trigger && time.Now().Before(timeout) {
		time.Sleep(10 * time.Millisecond)
	}

	return *trigger
}

// UpdateFromSensorData updates sensor values from raw sensor data
func (s *ActiveSensor) UpdateFromSensorData(data []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store raw values
	s.valuesAsString = make([]string, len(data))
	copy(s.valuesAsString, data)
	s.hasValueAsStringUpdated = true

	return nil
}

// GetTriggerFlag returns a pointer to the trigger flag for this sensor
func (s *ActiveSensor) GetTriggerFlag() *bool {
	return &s.triggerFlag
}
