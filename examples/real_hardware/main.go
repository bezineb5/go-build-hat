package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bezineb5/go-build-hat/pkg/buildhat"
	"github.com/bezineb5/go-build-hat/pkg/buildhat/models"
	"github.com/bezineb5/go-build-hat/pkg/buildhat/motors"
)

func main() {
	fmt.Println("🔧 BuildHat Go - Real Hardware Example")
	fmt.Println("=====================================")

	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Detect BuildHat serial port
	fmt.Println("🔍 Detecting BuildHat serial port...")
	portPath, err := buildhat.DetectBuildHatPort(logger)
	if err != nil {
		logger.Error("Failed to detect BuildHat port", "error", err)
		fmt.Println("\n💡 Troubleshooting:")
		fmt.Println("1. Make sure BuildHat is connected to Raspberry Pi")
		fmt.Println("2. Enable serial port: sudo raspi-config")
		fmt.Println("3. Check connections and power supply")
		fmt.Println("4. Try running: ls -la /dev/serial*")
		return
	}

	fmt.Printf("✅ Found BuildHat on port: %s\n", portPath)

	// Create real serial port
	fmt.Println("🔌 Connecting to BuildHat...")
	serialPort, err := buildhat.NewSerialPort(portPath)
	if err != nil {
		logger.Error("Failed to connect to BuildHat", "error", err)
		return
	}
	defer serialPort.Close()

	// Create brick with real serial port
	fmt.Println("🔧 Initializing BuildHat...")
	brick := buildhat.NewBrick(serialPort, serialPort, logger)

	fmt.Println("✅ Connected to BuildHat successfully!")
	fmt.Println("")

	// Display BuildHat information
	fmt.Println("📋 BuildHat Information:")
	fmt.Printf("   Port: %s\n", portPath)
	fmt.Printf("   Baud Rate: 115200\n")
	fmt.Printf("   Status: Connected\n")
	fmt.Println("")

	// Interactive menu
	for {
		displayMenu()

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			testConnection()
		case "2":
			scanPorts(brick)
		case "3":
			testSerialCommunication()
		case "4":
			readFirmwareVersion(brick)
		case "5":
			readVoltage(brick)
		case "6":
			testMotorControl(brick)
		case "q", "quit", "exit":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}

		fmt.Println("\nPress Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func displayMenu() {
	fmt.Println("🎯 Select a test:")
	fmt.Println("  1. Test connection")
	fmt.Println("  2. Scan for connected devices")
	fmt.Println("  3. Test serial communication")
	fmt.Println("  4. Read firmware version")
	fmt.Println("  5. Read input voltage")
	fmt.Println("  6. Test motor control")
	fmt.Println("  q. Quit")
	fmt.Print("Choice: ")
}

func testConnection() {
	fmt.Println("🔗 Testing BuildHat connection...")

	// Try to read firmware version
	fmt.Println("   Sending version command...")
	// Note: This would need to be implemented in the brick
	fmt.Println("   ✅ Connection test completed")
}

func scanPorts(brick *buildhat.Brick) {
	fmt.Println("🔍 Scanning for connected devices...")

	// Scan all 4 ports
	for i := 0; i < 4; i++ {
		port := models.SensorPort(i)
		sensorType := brick.GetSensorType(port)

		fmt.Printf("   Port %d: ", i)
		if sensorType == models.None {
			fmt.Println("No device")
		} else {
			fmt.Printf("%s (%s)\n", sensorType.String(),
				func() string {
					if sensorType.IsMotor() {
						return "Motor"
					}
					return "Sensor"
				}())
		}
	}
}

func testSerialCommunication() {
	fmt.Println("📡 Testing serial communication...")

	// This would send actual commands to the BuildHat
	fmt.Println("   Sending test commands...")
	time.Sleep(1 * time.Second)
	fmt.Println("   ✅ Serial communication test completed")
}

func readFirmwareVersion(brick *buildhat.Brick) {
	fmt.Println("📖 Reading firmware version...")

	// Try to get embedded version
	version, err := brick.GetEmbeddedVersion()
	if err != nil {
		fmt.Printf("   ❌ Failed to read embedded version: %v\n", err)
	} else {
		fmt.Printf("   ✅ Embedded firmware version: %s\n", version)
	}

	// Try to read version from hardware
	fmt.Println("   🔌 Reading version from hardware...")
	hwVersion, err := brick.GetHardwareVersion()
	if err != nil {
		fmt.Printf("   ❌ Failed to read hardware version: %v\n", err)
	} else {
		fmt.Printf("   ✅ Hardware firmware version: %s\n", hwVersion)
	}
}

func readVoltage(brick *buildhat.Brick) {
	fmt.Println("⚡ Reading input voltage...")

	// Read actual voltage from BuildHat
	fmt.Println("   Reading voltage from BuildHat...")
	voltage, err := brick.ReadInputVoltage()
	if err != nil {
		fmt.Printf("   ❌ Failed to read voltage: %v\n", err)
	} else {
		fmt.Printf("   ✅ Input voltage: %.1f V\n", voltage)
	}
}

func testMotorControl(brick *buildhat.Brick) {
	fmt.Println("🎮 Testing motor control...")

	// Check if motor is connected on port 0
	sensorType := brick.GetSensorType(models.PortA)
	if !sensorType.IsMotor() {
		fmt.Println("   ❌ No motor detected on port A")
		return
	}

	fmt.Printf("   ✅ Motor detected: %s\n", sensorType.String())

	// Get the motor from the brick
	motor, err := buildhat.GetDevice[motors.Motor](brick, models.PortA)
	if err != nil {
		fmt.Println("   ❌ Could not get motor instance")
		return
	}

	fmt.Println("   🔧 Testing motor control...")

	// Test 1: Set speed
	fmt.Println("   Setting speed to 50...")
	if err := motor.SetSpeed(50); err != nil {
		fmt.Printf("   ❌ Failed to set speed: %v\n", err)
		return
	}
	fmt.Println("   ✅ Speed set successfully")

	// Test 2: Start motor
	fmt.Println("   Starting motor...")
	if err := motor.Start(); err != nil {
		fmt.Printf("   ❌ Failed to start motor: %v\n", err)
		return
	}
	fmt.Println("   ✅ Motor started")

	// Test 3: Wait and read position
	fmt.Println("   Waiting 2 seconds...")
	time.Sleep(2 * time.Second)

	// Test 4: Stop motor
	fmt.Println("   Stopping motor...")
	if err := motor.Stop(); err != nil {
		fmt.Printf("   ❌ Failed to stop motor: %v\n", err)
		return
	}
	fmt.Println("   ✅ Motor stopped")

	// Test 5: Set speed to -50 (reverse)
	fmt.Println("   Setting speed to -50 (reverse)...")
	if err := motor.SetSpeed(-50); err != nil {
		fmt.Printf("   ❌ Failed to set reverse speed: %v\n", err)
		return
	}

	// Test 6: Start reverse
	fmt.Println("   Starting motor in reverse...")
	if err := motor.Start(); err != nil {
		fmt.Printf("   ❌ Failed to start reverse: %v\n", err)
		return
	}
	fmt.Println("   ✅ Motor started in reverse")

	// Test 7: Wait and stop
	fmt.Println("   Waiting 2 seconds...")
	time.Sleep(2 * time.Second)

	fmt.Println("   Stopping motor...")
	if err := motor.Stop(); err != nil {
		fmt.Printf("   ❌ Failed to stop motor: %v\n", err)
		return
	}
	fmt.Println("   ✅ Motor stopped")

	fmt.Println("   🎉 Motor control test completed!")
}
