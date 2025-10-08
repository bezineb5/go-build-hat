package buildhat

import (
	"strings"
	"testing"
	"time"
)

func TestMotor_SetDefaultSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)

	// Test valid speed
	err := motor.SetDefaultSpeed(50)
	if err != nil {
		t.Fatalf("SetDefaultSpeed failed: %v", err)
	}

	if motor.defaultSpeed != 50 {
		t.Errorf("Expected default speed 50, got %d", motor.defaultSpeed)
	}

	// Test invalid speed (too high)
	err = motor.SetDefaultSpeed(150)
	if err == nil {
		t.Error("Expected error for speed > 100")
	}

	// Test invalid speed (too low)
	err = motor.SetDefaultSpeed(-150)
	if err == nil {
		t.Error("Expected error for speed < -100")
	}
}

func TestMotor_SetSpeedUnitRPM(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)

	motor.SetSpeedUnitRPM(true)
	if !motor.rpm {
		t.Error("Expected RPM mode to be enabled")
	}

	motor.SetSpeedUnitRPM(false)
	if motor.rpm {
		t.Error("Expected RPM mode to be disabled")
	}
}

func TestMotor_RunForRotations(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor position data
	mockPort.SimulateSensorResponse("0", 0, "0 0 0") // speed, position, aposition

	motor := brick.Motor(PortA)

	// Run for 2 rotations at speed 50
	err := motor.RunForRotations(2.0, 50)
	if err != nil {
		t.Fatalf("RunForRotations failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) < 2 {
		t.Fatalf("Expected at least 2 commands, got %d", len(writeHistory))
	}

	// Should contain select and ramp commands
	found := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "ramp") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ramp command to be sent")
	}
}

func TestMotor_RunForDegrees(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor position data
	mockPort.SimulateSensorResponse("0", 0, "0 0 0")

	motor := brick.Motor(PortA)

	// Run for 360 degrees at speed 50
	err := motor.RunForDegrees(360, 50)
	if err != nil {
		t.Fatalf("RunForDegrees failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	foundRamp := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "ramp") {
			foundRamp = true
		}
	}

	if !foundRamp {
		t.Error("Expected ramp command")
	}
}

func TestMotor_RunForDuration(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	// Run for 0.5 seconds to make test faster
	err := motor.RunForDuration(500*time.Millisecond, 50)

	if err != nil {
		t.Fatalf("RunForDuration failed: %v", err)
	}

	// Note: In tests with mock, we use a small delay for speed
	// Real hardware will take the actual specified duration

	// Verify pulse command was sent with correct speed value
	writeHistory := mockPort.GetWriteHistory()
	foundPulse := false
	foundCorrectSpeed := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "pulse") {
			foundPulse = true
			// Verify speed is 50, not 2.5 (should be "set pulse 50" not "set pulse 2.5")
			if strings.Contains(cmd, "pulse 50") || strings.Contains(cmd, "pulse 5.0") {
				foundCorrectSpeed = true
			}
		}
	}

	if !foundPulse {
		t.Error("Expected pulse command")
	}
	if !foundCorrectSpeed {
		t.Errorf("Expected 'set pulse 50', got commands: %v", writeHistory)
	}
}

func TestMotor_RunToPosition(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor data multiple times since getData() will request fresh data
	// Speed, position, absolute_position format: speed pos apos
	mockPort.SimulateSensorResponse("0", 0, "0 0 0")

	motor := brick.Motor(PortA)

	// Queue more sensor data since getData() sends select command and waits for response
	time.Sleep(20 * time.Millisecond) // Let first data be cached
	go func() {
		// Continuously provide sensor data while test runs
		for i := 0; i < 10; i++ {
			time.Sleep(50 * time.Millisecond)
			mockPort.SimulateSensorResponse("0", 0, "0 0 0")
		}
	}()

	// Test basic position move
	err := motor.RunToPosition(90, 50, DirectionShortest)
	if err != nil {
		t.Fatalf("RunToPosition failed: %v", err)
	}

	// Verify commands were sent
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("Expected commands to be sent")
	}

	// Should contain port selection and ramp/movement commands
	hasPortCmd := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "port 0") {
			hasPortCmd = true
			break
		}
	}
	if !hasPortCmd {
		t.Error("Expected port command to be sent")
	}
}

func TestMotor_RunToPosition_InvalidDirection(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)

	// Test with invalid direction (999 is not a valid MotorDirection)
	err := motor.RunToPosition(90, 50, MotorDirection(999))
	if err == nil {
		t.Error("Expected error for invalid direction")
	}
}

func TestMotor_RunToPosition_InvalidAngle(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)

	// Test with invalid angle (too large)
	err := motor.RunToPosition(200, 50, DirectionShortest)
	if err == nil {
		t.Error("Expected error for angle > 180")
	}

	// Test with invalid angle (too small)
	err = motor.RunToPosition(-200, 50, DirectionShortest)
	if err == nil {
		t.Error("Expected error for angle < -180")
	}
}

func TestMotor_StartStop(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	// Start motor
	err := motor.Start(50)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Verify motor is in free run mode
	if motor.runMode != MotorRunModeFree {
		t.Errorf("Expected run mode FREE, got %d", motor.runMode)
	}

	// Verify the actual speed value sent (should be 50, not 2.5!)
	writeHistory := mockPort.GetWriteHistory()
	foundCorrectSpeed := false
	for _, cmd := range writeHistory {
		// Check for "set 50" (with possible decimal), not "set 2.5"
		if strings.Contains(cmd, "set 5") && strings.Contains(cmd, "pid") {
			// More precise: should contain "set 50" not "set 2.5"
			if strings.Contains(cmd, "set 50") || strings.Contains(cmd, "set 5.0") {
				foundCorrectSpeed = true
			}
		}
	}
	if !foundCorrectSpeed {
		t.Errorf("Expected 'set 50' command, got commands: %v", writeHistory)
	}

	// Stop motor
	err = motor.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Verify motor is stopped
	if motor.runMode != MotorRunModeNone {
		t.Errorf("Expected run mode NONE, got %d", motor.runMode)
	}

	// Verify coast command was sent (get fresh history after Stop())
	writeHistory = mockPort.GetWriteHistory()
	foundCoast := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "coast") {
			foundCoast = true
		}
	}

	if !foundCoast {
		t.Error("Expected coast command")
	}
}

func TestMotor_Start_AlreadyRunning(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)
	motor.runMode = MotorRunModeFree
	motor.currentSpeed = 50

	// Try to start at same speed - should return immediately
	err := motor.Start(50)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
}

func TestMotor_Start_InvalidSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)

	// Test invalid speed
	err := motor.Start(150)
	if err == nil {
		t.Error("Expected error for speed > 100")
	}
}

func TestMotor_GetPosition(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor data: speed=10, position=720, aposition=0
	mockPort.SimulateSensorResponse("0", 0, "10 720 0")

	motor := brick.Motor(PortA)

	position, err := motor.GetPosition()
	if err != nil {
		t.Fatalf("GetPosition failed: %v", err)
	}

	if position != 720 {
		t.Errorf("Expected position 720, got %d", position)
	}
}

func TestMotor_GetAbsolutePosition(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor data: speed=10, position=720, aposition=90
	mockPort.SimulateSensorResponse("0", 0, "10 720 90")

	motor := brick.Motor(PortA)

	apos, err := motor.GetAbsolutePosition()
	if err != nil {
		t.Fatalf("GetAbsolutePosition failed: %v", err)
	}

	if apos != 90 {
		t.Errorf("Expected absolute position 90, got %d", apos)
	}
}

func TestMotor_GetSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor data: speed=25, position=0, aposition=0
	mockPort.SimulateSensorResponse("0", 0, "25 0 0")

	motor := brick.Motor(PortA)

	speed, err := motor.GetSpeed()
	if err != nil {
		t.Fatalf("GetSpeed failed: %v", err)
	}

	if speed != 25 {
		t.Errorf("Expected speed 25, got %d", speed)
	}
}

func TestMotor_SetPowerLimit(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	// Set valid power limit
	err := motor.SetPowerLimit(0.5)
	if err != nil {
		t.Fatalf("SetPowerLimit failed: %v", err)
	}

	// Verify command was sent
	writeHistory := mockPort.GetWriteHistory()
	found := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "port_plimit") && strings.Contains(cmd, "0.5") {
			found = true
		}
	}

	if !found {
		t.Error("Expected port_plimit command with value 0.5")
	}

	// Test invalid limit (too high)
	err = motor.SetPowerLimit(1.5)
	if err == nil {
		t.Error("Expected error for limit > 1")
	}

	// Test invalid limit (too low)
	err = motor.SetPowerLimit(-0.1)
	if err == nil {
		t.Error("Expected error for limit < 0")
	}
}

func TestMotor_SetPWMParams(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	// Set valid PWM params
	err := motor.SetPWMParams(0.7, 0.02)
	if err != nil {
		t.Fatalf("SetPWMParams failed: %v", err)
	}

	// Verify command was sent
	writeHistory := mockPort.GetWriteHistory()
	found := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "pwmparams") {
			found = true
		}
	}

	if !found {
		t.Error("Expected pwmparams command")
	}
}

func TestMotor_PWM(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	// Set valid PWM value
	err := motor.PWM(0.5)
	if err != nil {
		t.Fatalf("PWM failed: %v", err)
	}

	// Verify command was sent
	writeHistory := mockPort.GetWriteHistory()
	found := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "pwm") && strings.Contains(cmd, "set") {
			found = true
		}
	}

	if !found {
		t.Error("Expected pwm set command")
	}

	// Test invalid PWM value
	err = motor.PWM(1.5)
	if err == nil {
		t.Error("Expected error for PWM value > 1")
	}
}

func TestMotor_PresetPosition(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	err := motor.PresetPosition()
	if err != nil {
		t.Fatalf("PresetPosition failed: %v", err)
	}

	// Verify preset command was sent
	writeHistory := mockPort.GetWriteHistory()
	found := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "preset") {
			found = true
		}
	}

	if !found {
		t.Error("Expected preset command")
	}
}

func TestMotor_CoastAndFloat(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)

	// Test Coast
	err := motor.Coast()
	if err != nil {
		t.Fatalf("Coast failed: %v", err)
	}

	// Test Float (should be same as Coast)
	err = motor.Float()
	if err != nil {
		t.Fatalf("Float failed: %v", err)
	}

	// Verify coast commands were sent
	writeHistory := mockPort.GetWriteHistory()
	coastCount := 0
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "coast") {
			coastCount++
		}
	}

	if coastCount < 2 {
		t.Errorf("Expected at least 2 coast commands, got %d", coastCount)
	}
}

func TestMotor_SetRelease(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor(PortA)

	// Default should be true
	if !motor.release {
		t.Error("Expected default release to be true")
	}

	motor.SetRelease(false)
	if motor.release {
		t.Error("Expected release to be false")
	}

	motor.SetRelease(true)
	if !motor.release {
		t.Error("Expected release to be true")
	}
}
