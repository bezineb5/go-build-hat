# Go BuildHat Library

[![Tests](https://github.com/bezineb5/go-build-hat/workflows/Tests/badge.svg)](https://github.com/bezineb5/go-build-hat/actions/workflows/test.yml)
[![CI](https://github.com/bezineb5/go-build-hat/workflows/CI/badge.svg)](https://github.com/bezineb5/go-build-hat/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/Coverage-GitHub%20Native-blue.svg)](https://github.com/bezineb5/go-build-hat/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bezineb5/go-build-hat)](https://goreportcard.com/report/github.com/bezineb5/go-build-hat)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A Go library for controlling LEGO SPIKE Prime and Technic motors and sensors through the Raspberry Pi BuildHat.

This library is a Go port of the Python BuildHat library, providing a clean and idiomatic Go interface for communicating with LEGO SPIKE Prime and Technic devices through the Raspberry Pi BuildHat.

## Features

- **Generic I/O Interface**: Uses `io.Reader` and `io.Writer` interfaces, decoupled from any specific serial library
- **Comprehensive Motor Control**: Full motor control with speed, position, rotation, PWM, and PID parameters
- **Rich Sensor Support**: Color, distance, force, tilt, motion sensors, lights, and 3x3 LED matrix
- **Python-Compatible API**: High-level motor and sensor interfaces matching the Python BuildHat library
- **Thread-Safe**: All operations are thread-safe with proper mutex protection
- **Asynchronous I/O**: Future-based sensor data retrieval with timeout support
- **Firmware Management**: Automatic firmware version checking and updating
- **Structured Logging**: Built-in structured logging with Go's standard `slog` package
- **Well Tested**: 64%+ test coverage with comprehensive unit and integration tests

## Installation

```bash
go get github.com/bezineb5/go-build-hat
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "log/slog"
    "os"
    "time"

    "github.com/bezineb5/go-build-hat/pkg/buildhat"
    "go.bug.st/serial"
)

func main() {
    // Create your serial port connection
    port, err := serial.Open("/dev/serial0", &serial.Mode{
        BaudRate: 115200,
        DataBits: 8,
        Parity:   serial.NoParity,
        StopBits: serial.OneStopBit,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer port.Close()
    
    // Create logger
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    // Create brick instance
    brick := buildhat.NewBrick(port, port, logger)
    defer brick.Close()

    // Initialize the BuildHat
    if err := brick.Initialize(); err != nil {
        log.Fatal(err)
    }

    // Create motor on port A
    motor := brick.Motor("A")
    
    // Run motor for 2 seconds at 50% speed
    if err := motor.RunForSeconds(2, 50); err != nil {
        log.Fatal(err)
    }
    
    // Or run for specific rotations
    if err := motor.RunForRotations(2.5, 75); err != nil {
        log.Fatal(err)
    }
    
    // Create a color sensor on port B
    colorSensor := brick.ColorSensor("B")
    color, err := colorSensor.GetColor()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Detected color: %s", color)
}
```

## Serial Port Integration

The library uses `io.Reader` and `io.Writer` interfaces, so you can use any serial communication library. Here are some popular options:

### Using go-serial

```go
import "go.bug.st/serial"

func createSerialPort() (io.ReadWriteCloser, error) {
    port, err := serial.Open("/dev/ttyAMA0", &serial.Mode{
        BaudRate: 115200,
        DataBits: 8,
        Parity:   serial.NoParity,
        StopBits: serial.OneStopBit,
    })
    return port, err
}
```

### Using tarm/serial

```go
import "github.com/tarm/serial"

func createSerialPort() (io.ReadWriteCloser, error) {
    config := &serial.Config{
        Name: "/dev/ttyAMA0",
        Baud: 115200,
    }
    return serial.OpenPort(config)
}
```

## API Reference

### Brick

The main `Brick` struct provides core functionality for device initialization and management.

#### Constructor

```go
func NewBrick(reader io.Reader, writer io.Writer, logger *slog.Logger) *Brick
```

#### Initialization

```go
func (b *Brick) Initialize() error
```

#### Device Access

```go
func (b *Brick) Motor(port string) *Motor                          // "A", "B", "C", "D"
func (b *Brick) PassiveMotor(port string) *PassiveMotor
func (b *Brick) ColorSensor(port string) *ColorSensor
func (b *Brick) DistanceSensor(port string) *DistanceSensor
func (b *Brick) ForceSensor(port string) *ForceSensor
func (b *Brick) ButtonSensor(port string) *ButtonSensor
func (b *Brick) ColorDistanceSensor(port string) *ColorDistanceSensor
func (b *Brick) TiltSensor(port string) *TiltSensor
func (b *Brick) MotionSensor(port string) *MotionSensor
func (b *Brick) Light(port string) *Light
func (b *Brick) Matrix(port string) *Matrix
```

#### Device Information

```go
func (b *Brick) ListDevices() []DeviceInfo
func (b *Brick) GetConnectedDevices() []DeviceInfo
```

#### Firmware Management

```go
func (b *Brick) GetHardwareVersion() (string, error)
func (b *Brick) CheckFirmwareVersion() (bool, error)
func (b *Brick) CheckAndUpdateFirmware() error
```

### Motor

High-level motor interface with Python-compatible API.

```go
// Configuration
func (m *Motor) SetDefaultSpeed(speed int) error           // -100 to 100
func (m *Motor) SetSpeedUnitRPM(rpm bool)
func (m *Motor) SetPowerLimit(limit float64) error         // 0.0 to 1.0
func (m *Motor) SetPWMParams(pwmThresh, minPWM float64) error
func (m *Motor) SetRelease(release bool)

// Movement
func (m *Motor) RunForSeconds(seconds float64, speed int) error
func (m *Motor) RunForDegrees(degrees, speed int) error
func (m *Motor) RunForRotations(rotations float64, speed int) error
func (m *Motor) RunToPosition(degrees, speed int, direction MotorDirection) error
func (m *Motor) Start(speed int) error
func (m *Motor) Stop() error

// Low-level control
func (m *Motor) PWM(value float64) error                   // -1.0 to 1.0
func (m *Motor) Coast() error
func (m *Motor) Float() error

// Status
func (m *Motor) GetPosition() (int, error)
func (m *Motor) GetAbsolutePosition() (int, error)
func (m *Motor) GetSpeed() (int, error)

// Calibration
func (m *Motor) PresetPosition() error
```

#### Motor Direction

```go
const (
    DirectionShortest      MotorDirection = iota
    DirectionClockwise
    DirectionAnticlockwise
)
```

### Sensors

#### ColorSensor

```go
func (c *ColorSensor) GetColor() (string, error)           // "red", "blue", "green", etc.
func (c *ColorSensor) GetReflectedLight() (int, error)     // 0-100
func (c *ColorSensor) GetAmbientLight() (int, error)       // 0-100
```

#### DistanceSensor

```go
func (d *DistanceSensor) GetDistance() (int, error)        // mm
```

#### ForceSensor

```go
func (f *ForceSensor) GetForce() (int, error)              // Newtons * 10
```

#### ButtonSensor

```go
func (b *ButtonSensor) IsPressed() (bool, error)
```

#### ColorDistanceSensor

```go
func (c *ColorDistanceSensor) GetColor() (string, error)
func (c *ColorDistanceSensor) GetDistance() (int, error)
func (c *ColorDistanceSensor) GetReflectedLight() (int, error)
func (c *ColorDistanceSensor) GetRGB() (r, g, b uint8, err error)
```

#### TiltSensor

```go
func (t *TiltSensor) GetTilt() (x, y, z int, err error)
func (t *TiltSensor) GetDirection() (string, error)       // "up", "down", "left", "right", "level"
```

#### MotionSensor

```go
func (m *MotionSensor) GetDistance() (int, error)
func (m *MotionSensor) GetMovementCount() (int, error)
```

#### Light

```go
func (l *Light) On() error
func (l *Light) Off() error
func (l *Light) SetBrightness(brightness int) error       // 0-100
func (l *Light) GetBrightness() (int, error)
```

#### Matrix (3x3 LED Matrix)

```go
func (m *Matrix) SetPixel(x, y, brightness int) error     // brightness: 0-10
func (m *Matrix) SetAll(brightness int) error
func (m *Matrix) SetRow(row, brightness int) error
func (m *Matrix) SetColumn(col, brightness int) error
func (m *Matrix) Clear() error
```

### Device Types

```go
type DeviceCategory int

const (
    DeviceCategoryUnknown DeviceCategory = iota
    DeviceCategoryDisconnected
    DeviceCategoryMotor
    DeviceCategorySensor
    DeviceCategoryPassiveMotor
    DeviceCategoryLight
)
```

## Supported Devices

### Active Motors (with position feedback)
- Medium Linear Motor (ID: 38)
- Large Motor (ID: 46)
- XL Motor (ID: 47)
- Medium Angular Motor - Cyan (ID: 48)
- Large Angular Motor - Cyan (ID: 49)
- Small Angular Motor (ID: 65)
- Medium Angular Motor - Grey (ID: 75)
- Large Angular Motor - Grey (ID: 76)

### Passive Motors (no position feedback)
- Train Motor (ID: 1, 2)

### Sensors
- Tilt Sensor (ID: 34)
- Motion Sensor (ID: 35)
- Color and Distance Sensor (ID: 37)
- Color Sensor (ID: 61)
- Distance Sensor (ID: 62)
- Force Sensor (ID: 63)
- 3x3 Color Light Matrix (ID: 64)

### Other
- Light (ID: 8)

## Error Handling

All methods return errors that should be checked:

```go
motor := brick.Motor("A")
if err := motor.RunForSeconds(2, 50); err != nil {
    log.Printf("Failed to run motor: %v", err)
}

sensor := brick.ColorSensor("B")
color, err := sensor.GetColor()
if err != nil {
    log.Printf("Failed to read sensor: %v", err)
}
```

Common errors include:
- Timeout errors when waiting for sensor data
- Invalid parameter ranges (e.g., speed > 100, brightness > 10)
- Device not connected or wrong device type
- Serial communication errors

## Thread Safety

The library is thread-safe and can be used from multiple goroutines:
- All internal state is protected by mutexes
- Sensor data retrieval uses futures with timeout protection
- Background reader goroutine handles asynchronous I/O

```go
// Safe to use from multiple goroutines
go func() {
    motor1 := brick.Motor("A")
    motor1.RunForSeconds(2, 50)
}()

go func() {
    motor2 := brick.Motor("B")
    motor2.RunForSeconds(2, -50)
}()
```

## Logging

The library uses Go's standard structured logging with `slog`. You can configure the log level and handler:

```go
import (
    "log/slog"
    "os"
)

// Text handler with debug level
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// JSON handler for structured output
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
```

## Examples

See the `examples/` directory for complete working examples.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

This library is a Go port of the [Python BuildHat library](https://github.com/RaspberryPiFoundation/python-build-hat) by the Raspberry Pi Foundation.

