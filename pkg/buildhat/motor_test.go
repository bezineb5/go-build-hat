package buildhat

import (
	"strings"
	"testing"
	"time"
)

func TestMotor_SetDefaultSpeed(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor("A")

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

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	// Queue motor position data
	mockPort.SimulateSensorResponse("0", 0, "0 0 0") // speed, position, aposition

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	// Queue motor position data
	mockPort.SimulateSensorResponse("0", 0, "0 0 0")

	motor := brick.Motor("A")

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

func TestMotor_RunForSeconds(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

	// Run for 0.5 seconds to make test faster
	start := time.Now()
	err := motor.RunForSeconds(0.5, 50)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("RunForSeconds failed: %v", err)
	}

	// Should take approximately 0.5 seconds
	if elapsed < 400*time.Millisecond || elapsed > 800*time.Millisecond {
		t.Errorf("Expected ~500ms, got %v", elapsed)
	}

	// Verify pulse command was sent
	writeHistory := mockPort.GetWriteHistory()
	foundPulse := false
	for _, cmd := range writeHistory {
		if strings.Contains(cmd, "pulse") {
			foundPulse = true
		}
	}

	if !foundPulse {
		t.Error("Expected pulse command")
	}
}

func TestMotor_RunToPosition(t *testing.T) {
	t.Skip("Skipping - requires motor position data which has async timing issues in mock environment")

	// Note: RunToPosition calls getData() internally which requires motor data to be available
	// In real hardware, the motor continuously sends position data
	// Testing this properly would require a more sophisticated mock that simulates continuous data
}

func TestMotor_RunToPosition_InvalidDirection(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor("A")

	// Test with invalid direction (999 is not a valid MotorDirection)
	err := motor.RunToPosition(90, 50, MotorDirection(999))
	if err == nil {
		t.Error("Expected error for invalid direction")
	}
}

func TestMotor_RunToPosition_InvalidAngle(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

	// Start motor
	err := motor.Start(50)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Verify motor is in free run mode
	if motor.runMode != MotorRunModeFree {
		t.Errorf("Expected run mode FREE, got %d", motor.runMode)
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

	// Verify coast command was sent
	writeHistory := mockPort.GetWriteHistory()
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

	brick.SetupMockScanner()

	motor := brick.Motor("A")
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

	motor := brick.Motor("A")

	// Test invalid speed
	err := motor.Start(150)
	if err == nil {
		t.Error("Expected error for speed > 100")
	}
}

func TestMotor_GetPosition(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	// Queue motor data: speed=10, position=720, aposition=0
	mockPort.SimulateSensorResponse("0", 0, "10 720 0")

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	// Queue motor data: speed=10, position=720, aposition=90
	mockPort.SimulateSensorResponse("0", 0, "10 720 90")

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	// Queue motor data: speed=25, position=0, aposition=0
	mockPort.SimulateSensorResponse("0", 0, "25 0 0")

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

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

	brick.SetupMockScanner()
	mockPort := brick.GetMockPort()

	motor := brick.Motor("A")

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

	motor := brick.Motor("A")

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
