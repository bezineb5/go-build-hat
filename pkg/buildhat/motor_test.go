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
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Run for 2 rotations at speed 50
	err := motor.RunForRotations(2.0, 50)
	if err != nil {
		t.Fatalf("RunForRotations failed: %v", err)
	}

	// Verify EXACT command sequence
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("No commands were sent")
	}

	// Should have a ramp command with EXACT format:
	// "port 0 ; select 0 ; selrate 10 ; pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01 ; set ramp 0.000000 2.000000 <duration> 0\r"
	expectedPrefix := "port 0 ; select 0 ; selrate 10 ; pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01 ; set ramp 0.000000 2.000000"
	found := false
	for _, cmd := range writeHistory {
		if strings.HasPrefix(cmd, expectedPrefix) && strings.HasSuffix(cmd, " 0\r") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected exact ramp command with prefix '%s', got: %v", expectedPrefix, writeHistory)
	}

	// Should also have a coast command at the end
	expectedCoast := "port 0 ; coast\r"
	foundCoast := false
	for _, cmd := range writeHistory {
		if cmd == expectedCoast {
			foundCoast = true
			break
		}
	}
	if !foundCoast {
		t.Errorf("Expected coast command '%s', got: %v", expectedCoast, writeHistory)
	}
}

func TestMotor_RunForDegrees(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor position data
	mockPort.SimulateSensorResponse("0", 0, "0 0 0")

	motor := brick.Motor(PortA)
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Run for 360 degrees at speed 50
	err := motor.RunForDegrees(360, 50)
	if err != nil {
		t.Fatalf("RunForDegrees failed: %v", err)
	}

	// Verify EXACT command sequence
	writeHistory := mockPort.GetWriteHistory()
	if len(writeHistory) == 0 {
		t.Fatal("No commands were sent")
	}

	// Should have exact ramp command:
	// "port 0 ; select 0 ; selrate 10 ; pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01 ; set ramp 0.000000 1.000000 <duration> 0\r"
	// 360 degrees = 1 rotation
	expectedPrefix := "port 0 ; select 0 ; selrate 10 ; pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01 ; set ramp 0.000000 1.000000"
	found := false
	for _, cmd := range writeHistory {
		if strings.HasPrefix(cmd, expectedPrefix) && strings.HasSuffix(cmd, " 0\r") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected exact ramp command with prefix '%s', got: %v", expectedPrefix, writeHistory)
	}

	// Should have coast command
	expectedCoast := "port 0 ; coast\r"
	foundCoast := false
	for _, cmd := range writeHistory {
		if cmd == expectedCoast {
			foundCoast = true
		}
	}
	if !foundCoast {
		t.Errorf("Expected coast command '%s', got: %v", expectedCoast, writeHistory)
	}
}

func TestMotor_RunForDuration(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Run for 0.5 seconds to make test faster
	err := motor.RunForDuration(500*time.Millisecond, 50)

	if err != nil {
		t.Fatalf("RunForDuration failed: %v", err)
	}

	// Verify EXACT pulse command was sent:
	// "port 0 ; select 0 ; selrate 10 ; pid 0 0 0 s1 1 0 0.003 0.01 0 100 0.01 ; set pulse 50.000000 0.0 0.500000 0\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedPulse := "port 0 ; select 0 ; selrate 10 ; pid 0 0 0 s1 1 0 0.003 0.01 0 100 0.01 ; set pulse 50.000000 0.0 0.500000 0\r"
	foundPulse := false
	for _, cmd := range writeHistory {
		if cmd == expectedPulse {
			foundPulse = true
			break
		}
	}
	if !foundPulse {
		t.Errorf("Expected exact pulse command '%s', got: %v", expectedPulse, writeHistory)
	}

	// Should have coast command
	expectedCoast := "port 0 ; coast\r"
	foundCoast := false
	for _, cmd := range writeHistory {
		if cmd == expectedCoast {
			foundCoast = true
		}
	}
	if !foundCoast {
		t.Errorf("Expected coast command '%s', got: %v", expectedCoast, writeHistory)
	}
}

func TestMotor_RunToPosition(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	// Queue motor data: speed=0, position=0, aposition=0
	mockPort.SimulateSensorResponse("0", 0, "0 0 0")

	motor := brick.Motor(PortA)
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Queue more sensor data since getData() sends select command and waits for response
	time.Sleep(20 * time.Millisecond) // Let first data be cached
	go func() {
		// Continuously provide sensor data while test runs
		for i := 0; i < 10; i++ {
			time.Sleep(50 * time.Millisecond)
			mockPort.SimulateSensorResponse("0", 0, "0 0 0")
		}
	}()

	// Test basic position move: 90 degrees at 50% speed from position 0
	// Expected calculation:
	// - Current pos: 0 degrees = 0.0 rotations
	// - Target: 90 degrees absolute = 0.25 rotations
	// - Duration: (0.25 - 0.0) / (50 * 0.05) = 0.1 seconds
	err := motor.RunToPosition(90, 50, DirectionShortest)
	if err != nil {
		t.Fatalf("RunToPosition failed: %v", err)
	}

	// Verify EXACT ramp command
	writeHistory := mockPort.GetWriteHistory()
	expectedRamp := "port 0 ; select 0 ; selrate 10 ; pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01 ; set ramp 0.000000 0.250000 0.100000 0\r"
	foundRamp := false
	for _, cmd := range writeHistory {
		if cmd == expectedRamp {
			foundRamp = true
			break
		}
	}
	if !foundRamp {
		t.Errorf("Expected exact ramp command '%s', got: %v", expectedRamp, writeHistory)
	}

	// Should have coast command at end
	expectedCoast := "port 0 ; coast\r"
	foundCoast := false
	for _, cmd := range writeHistory {
		if cmd == expectedCoast {
			foundCoast = true
		}
	}
	if !foundCoast {
		t.Errorf("Expected coast command '%s', got: %v", expectedCoast, writeHistory)
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
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Start motor
	err := motor.Start(50)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Verify motor is in free run mode
	if motor.runMode != MotorRunModeFree {
		t.Errorf("Expected run mode FREE, got %d", motor.runMode)
	}

	// Verify EXACT start command:
	// "port 0 ; select 0 ; selrate 10 ; pid 0 0 0 s1 1 0 0.003 0.01 0 100 0.01 ; set 50.000000\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedStart := "port 0 ; select 0 ; selrate 10 ; pid 0 0 0 s1 1 0 0.003 0.01 0 100 0.01 ; set 50.000000\r"
	foundStart := false
	for _, cmd := range writeHistory {
		if cmd == expectedStart {
			foundStart = true
			break
		}
	}
	if !foundStart {
		t.Errorf("Expected exact start command '%s', got: %v", expectedStart, writeHistory)
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

	// Verify EXACT coast command
	writeHistory = mockPort.GetWriteHistory()
	expectedCoast := "port 0 ; coast\r"
	foundCoast := false
	for _, cmd := range writeHistory {
		if cmd == expectedCoast {
			foundCoast = true
			break
		}
	}

	if !foundCoast {
		t.Errorf("Expected exact coast command '%s', got: %v", expectedCoast, writeHistory)
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
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Set valid power limit
	err := motor.SetPowerLimit(0.5)
	if err != nil {
		t.Fatalf("SetPowerLimit failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; port_plimit 0.50\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedCmd := "port 0 ; port_plimit 0.50\r"
	found := false
	for _, cmd := range writeHistory {
		if cmd == expectedCmd {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected exact command '%s', got: %v", expectedCmd, writeHistory)
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
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Set valid PWM params
	err := motor.SetPWMParams(0.7, 0.02)
	if err != nil {
		t.Fatalf("SetPWMParams failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; pwmparams 0.70 0.02\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedCmd := "port 0 ; pwmparams 0.70 0.02\r"
	found := false
	for _, cmd := range writeHistory {
		if cmd == expectedCmd {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected exact command '%s', got: %v", expectedCmd, writeHistory)
	}
}

func TestMotor_PWM(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Set valid PWM value
	err := motor.PWM(0.5)
	if err != nil {
		t.Fatalf("PWM failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; pwm ; set 0.50\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedCmd := "port 0 ; pwm ; set 0.50\r"
	found := false
	for _, cmd := range writeHistory {
		if cmd == expectedCmd {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected exact command '%s', got: %v", expectedCmd, writeHistory)
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
	mockPort.ClearWriteHistory() // Clear initialization commands

	err := motor.PresetPosition()
	if err != nil {
		t.Fatalf("PresetPosition failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; preset\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedCmd := "port 0 ; preset\r"
	found := false
	for _, cmd := range writeHistory {
		if cmd == expectedCmd {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected exact command '%s', got: %v", expectedCmd, writeHistory)
	}
}

func TestMotor_CoastAndFloat(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	mockPort := brick.GetMockPort()

	motor := brick.Motor(PortA)
	mockPort.ClearWriteHistory() // Clear initialization commands

	// Test Coast
	err := motor.Coast()
	if err != nil {
		t.Fatalf("Coast failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; coast\r"
	writeHistory := mockPort.GetWriteHistory()
	expectedCmd := "port 0 ; coast\r"
	coastCount := 0
	for _, cmd := range writeHistory {
		if cmd == expectedCmd {
			coastCount++
		}
	}
	if coastCount < 1 {
		t.Errorf("Expected at least 1 exact coast command '%s', got: %v", expectedCmd, writeHistory)
	}

	mockPort.ClearWriteHistory()

	// Test Float (should be same as Coast)
	err = motor.Float()
	if err != nil {
		t.Fatalf("Float failed: %v", err)
	}

	// Verify Float also sends exact coast command
	writeHistory = mockPort.GetWriteHistory()
	found := false
	for _, cmd := range writeHistory {
		if cmd == expectedCmd {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected exact command '%s' for Float, got: %v", expectedCmd, writeHistory)
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
