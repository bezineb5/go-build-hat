package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bezineb5/go-build-hat/pkg/buildhat"
	"go.bug.st/serial"
)

func main() {
	// Set up logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	fmt.Println("🚀 BuildHat Real Hardware Test")
	fmt.Println("==============================")

	// Initialize BuildHat
	port, err := serial.Open("/dev/serial0", &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		logger.Error("Failed to open serial port", "error", err)
		return
	}
	defer port.Close()

	brick := buildhat.NewBrick(port, port, logger)
	defer brick.Close()

	// Initialize the BuildHat
	if err := brick.Initialize(); err != nil {
		logger.Error("Failed to initialize BuildHat", "error", err)
		return
	}

	fmt.Println("✅ BuildHat initialized successfully!")
	fmt.Println()

	// Main menu loop
	reader := bufio.NewReader(os.Stdin)
	for {
		showMenu()
		fmt.Print("Enter your choice: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			testConnection(brick)
		case "2":
			scanDevices(brick, logger)
		case "3":
			testVersion(brick, logger)
		case "4":
			testVoltage(brick, logger)
		case "5":
			testMotorControl(brick, logger)
		case "q", "quit", "exit":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}

		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("📋 Available Tests:")
	fmt.Println("1. Test Connection")
	fmt.Println("2. Scan for Connected Devices")
	fmt.Println("3. Read Firmware Version")
	fmt.Println("4. Read Input Voltage")
	fmt.Println("5. Test Motor Control")
	fmt.Println("q. Quit")
	fmt.Println()
}

func testConnection(brick *buildhat.Brick) {
	fmt.Println("🔌 Testing connection...")

	// Try to get device info
	devices := brick.GetDeviceInfo()
	fmt.Printf("✅ Connection successful! Found %d ports\n", len(devices))

	for port, info := range devices {
		status := "❌ Disconnected"
		if info.Connected {
			status = "✅ Connected"
		}
		fmt.Printf("   Port %s: %s (%s)\n", port, info.Name, status)
	}
}

func scanDevices(brick *buildhat.Brick, logger *slog.Logger) {
	fmt.Println("🔍 Scanning for connected devices...")

	if err := brick.ScanDevices(); err != nil {
		logger.Error("Failed to scan devices", "error", err)
		return
	}

	// Wait a bit for scanning to complete
	time.Sleep(2 * time.Second)

	devices := brick.GetDeviceInfo()
	for port, info := range devices {
		if info.Connected {
			fmt.Printf("Port %s: %s (%s)\n", port, info.Name, info.DeviceType)
		} else {
			fmt.Printf("Port %s: No device\n", port)
		}
	}

	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func testVersion(brick *buildhat.Brick, logger *slog.Logger) {
	fmt.Println("📖 Reading firmware version...")

	// Get embedded firmware version
	embeddedVersion, err := brick.GetEmbeddedFirmwareVersion()
	if err != nil {
		logger.Error("Failed to get embedded firmware version", "error", err)
		fmt.Println("❌ Failed to get embedded firmware version")
	} else {
		fmt.Printf("📦 Embedded firmware version: %s\n", embeddedVersion)
	}

	// Get current hardware version
	currentVersion, err := brick.GetHardwareVersion()
	if err != nil {
		logger.Error("Failed to read hardware version", "error", err)
		fmt.Println("❌ Failed to read hardware version")
		return
	}
	fmt.Printf("🔧 Current firmware version: %s\n", currentVersion)

	// Check if versions match
	versionsMatch, err := brick.CheckFirmwareVersion()
	if err != nil {
		logger.Error("Failed to check firmware version", "error", err)
		fmt.Println("❌ Failed to check firmware version")
		return
	}

	if versionsMatch {
		fmt.Println("✅ Firmware versions match - no update needed")
	} else {
		fmt.Println("⚠️  Firmware versions differ - update may be needed")
	}
}

func testVoltage(brick *buildhat.Brick, logger *slog.Logger) {
	fmt.Println("⚡ Reading input voltage...")

	voltage, err := brick.GetVoltage()
	if err != nil {
		logger.Error("Failed to read voltage", "error", err)
		return
	}

	fmt.Printf("✅ Input voltage: %.2f V\n", voltage)
}

func testMotorControl(brick *buildhat.Brick, logger *slog.Logger) {
	fmt.Println("🎮 Testing motor control...")

	// Check if we have any motors connected
	devices := brick.GetDeviceInfo()
	var motorPort string

	for port, info := range devices {
		if info.Connected && info.DeviceType == "Motor" {
			motorPort = port
			break
		}
	}

	if motorPort == "" {
		fmt.Println("❌ No motor found. Please connect a motor to any port.")
		return
	}

	fmt.Printf("✅ Motor detected: %s on port %s\n", devices[motorPort].Name, motorPort)

	// Create motor instance
	motor := brick.Motor(motorPort)

	// Configure motor
	fmt.Println("\n⚙️  Configuring motor...")
	motor.SetDefaultSpeed(30)
	motor.SetPowerLimit(0.8)
	fmt.Println("✅ Motor configured (default speed: 30%, power limit: 80%)")

	// Test 1: Read motor position and speed
	fmt.Println("\n📊 Reading motor status...")
	if position, err := motor.GetPosition(); err == nil {
		fmt.Printf("   Position: %d degrees\n", position)
	}
	if speed, err := motor.GetSpeed(); err == nil {
		fmt.Printf("   Speed: %d\n", speed)
	}
	if apos, err := motor.GetAbsolutePosition(); err == nil {
		fmt.Printf("   Absolute Position: %d degrees\n", apos)
	}

	// Test 2: Preset position to 0
	fmt.Println("\n🔄 Resetting position to 0...")
	if err := motor.PresetPosition(); err != nil {
		logger.Error("Failed to preset position", "error", err)
	} else {
		fmt.Println("✅ Position reset")
	}
	time.Sleep(500 * time.Millisecond)

	// Test 3: Run for 1 rotation
	fmt.Println("\n🔄 Running motor for 1 rotation at 50% speed...")
	if err := motor.RunForRotations(1.0, 50); err != nil {
		logger.Error("Failed to run motor", "error", err)
	} else {
		fmt.Println("✅ Completed 1 rotation")
	}

	// Test 4: Run for 360 degrees
	fmt.Println("\n🔄 Running motor for 360 degrees at 40% speed...")
	if err := motor.RunForDegrees(360, 40); err != nil {
		logger.Error("Failed to run motor", "error", err)
	} else {
		fmt.Println("✅ Completed 360 degrees")
	}

	// Test 5: Run for 2 seconds
	fmt.Println("\n⏱️  Running motor for 2 seconds at 50% speed...")
	if err := motor.RunForSeconds(2, 50); err != nil {
		logger.Error("Failed to run motor", "error", err)
	} else {
		fmt.Println("✅ Completed 2 seconds")
	}

	// Test 6: Run for 2 seconds in reverse
	fmt.Println("\n⏱️  Running motor for 2 seconds at -50% speed (reverse)...")
	if err := motor.RunForSeconds(2, -50); err != nil {
		logger.Error("Failed to run motor in reverse", "error", err)
	} else {
		fmt.Println("✅ Completed reverse run")
	}

	// Test 7: Start/Stop motor
	fmt.Println("\n▶️  Starting motor in free-run mode at 30% speed...")
	if err := motor.Start(30); err != nil {
		logger.Error("Failed to start motor", "error", err)
	} else {
		fmt.Println("✅ Motor started")
		time.Sleep(2 * time.Second)

		fmt.Println("⏸️  Stopping motor...")
		if err := motor.Stop(); err != nil {
			logger.Error("Failed to stop motor", "error", err)
		} else {
			fmt.Println("✅ Motor stopped")
		}
	}

	// Test 8: Direct PWM control
	fmt.Println("\n⚡ Testing PWM control (0.5 = 50% power)...")
	if err := motor.PWM(0.5); err != nil {
		logger.Error("Failed to set PWM", "error", err)
	} else {
		fmt.Println("✅ PWM set to 50%")
		time.Sleep(1 * time.Second)
		motor.Coast()
		fmt.Println("✅ Motor coasted")
	}

	// Test 9: Run to position (if supported)
	fmt.Println("\n🎯 Running to position 90° (shortest path)...")
	if err := motor.RunToPosition(90, 40, buildhat.DirectionShortest); err != nil {
		logger.Warn("RunToPosition not available or failed", "error", err)
	} else {
		fmt.Println("✅ Reached position 90°")

		// Return to 0
		fmt.Println("🎯 Returning to position 0°...")
		if err := motor.RunToPosition(0, 40, buildhat.DirectionShortest); err == nil {
			fmt.Println("✅ Returned to position 0°")
		}
	}

	// Final status
	fmt.Println("\n📊 Final motor status:")
	if position, err := motor.GetPosition(); err == nil {
		fmt.Printf("   Position: %d degrees\n", position)
	}
	if speed, err := motor.GetSpeed(); err == nil {
		fmt.Printf("   Speed: %d\n", speed)
	}

	fmt.Println("\n🎉 Motor control test completed!")
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}
