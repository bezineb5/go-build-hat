# Go BuildHat Library

[![Tests](https://github.com/bezineb5/go-build-hat/workflows/Tests/badge.svg)](https://github.com/bezineb5/go-build-hat/actions/workflows/test.yml)
[![CI](https://github.com/bezineb5/go-build-hat/workflows/CI/badge.svg)](https://github.com/bezineb5/go-build-hat/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/Coverage-GitHub%20Native-blue.svg)](https://github.com/bezineb5/go-build-hat/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bezineb5/go-build-hat)](https://goreportcard.com/report/github.com/bezineb5/go-build-hat)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A Go library for controlling LEGO SPIKE Prime and Technic motors and sensors through the Raspberry Pi BuildHat.

This library is a Go port of the .NET BuildHat library, providing a clean interface for communicating with LEGO SPIKE Prime and Technic devices through the Raspberry Pi BuildHat.

## Features

- **Generic I/O Interface**: Uses `io.Reader` and `io.Writer` interfaces, allowing you to choose any serial communication library
- **Motor Control**: Full control over LEGO motors including speed, position, and power management
- **Sensor Support**: Support for various LEGO sensors including color, distance, force, and tilt sensors
- **Thread-Safe**: All operations are thread-safe with proper mutex protection
- **Context Support**: Full context support for cancellation and timeouts
- **Structured Logging**: Built-in structured logging with Go's standard `slog` package

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
    "github.com/bezineb5/go-build-hat/pkg/buildhat/models"
)

func main() {
    // Create your serial port connection
    // This example uses a mock, but you would use a real serial port
    serialPort := createSerialPort() // Your serial port implementation
    
    // Create logger
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    // Create brick instance
    brick := buildhat.NewBrick(serialPort, serialPort, logger)
    defer brick.Close()

    // Wait for initialization
    time.Sleep(2 * time.Second)

    // Set LED mode
    if err := brick.SetLedMode(models.Green); err != nil {
        log.Fatal(err)
    }

    // Set motor power
    if err := brick.SetMotorPower(models.PortA, 50); err != nil {
        log.Fatal(err)
    }

    // Move motor for 2 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := brick.MoveMotorForSeconds(models.PortA, 2.0, 75, true, ctx); err != nil {
        log.Fatal(err)
    }

    // Float the motor
    if err := brick.FloatMotor(models.PortA); err != nil {
        log.Fatal(err)
    }
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

The main `Brick` struct provides all the functionality for controlling motors and sensors.

#### Constructor

```go
func NewBrick(reader io.Reader, writer io.Writer, logger *slog.Logger) *Brick
```

#### Motor Control Methods

- `SetMotorPower(port SensorPort, powerPercent int) error` - Set motor power (-100 to 100)
- `SetMotorLimits(port SensorPort, powerLimit float64) error` - Set motor power limit (0 to 1)
- `SetMotorBias(port SensorPort, bias float64) error` - Set motor bias (0 to 1)
- `MoveMotorForSeconds(port SensorPort, seconds float64, speed int, blocking bool, ctx context.Context) error` - Move motor for specified time
- `FloatMotor(port SensorPort) error` - Float the motor (remove all constraints)

#### Sensor Control Methods

- `SelectModeAndRead(port SensorPort, mode int, readOnce bool) error` - Select sensor mode
- `SelectCombiModesAndRead(port SensorPort, modes []int, readOnce bool) error` - Select multiple sensor modes
- `StopContinuousReadingSensor(port SensorPort) error` - Stop continuous sensor reading
- `SwitchSensorOn(port SensorPort) error` - Turn sensor on
- `SwitchSensorOff(port SensorPort) error` - Turn sensor off

#### Utility Methods

- `SetLedMode(mode LedMode) error` - Set LED mode
- `GetLedMode() LedMode` - Get current LED mode
- `GetInputVoltage() float64` - Get input voltage
- `GetSensorType(port SensorPort) SensorType` - Get sensor type at port
- `SendRawCommand(command string) error` - Send raw command
- `ClearFaults() error` - Clear any faults
- `Close() error` - Close and cleanup

### Models

#### SensorPort

```go
const (
    PortA SensorPort = iota
    PortB
    PortC
    PortD
)
```

#### SensorType

Various sensor types including:
- `SystemMediumMotor`, `SystemTrainMotor`, `SystemTurntableMotor`
- `SpikePrimeMediumMotor`, `SpikePrimeLargeMotor`
- `SpikePrimeColorSensor`, `SpikePrimeUltrasonicDistanceSensor`
- `SpikePrimeForceSensor`, `WeDoTiltSensor`, `WeDoDistanceSensor`
- And many more...

#### LedMode

```go
const (
    VoltageDependant LedMode = -1
    Off             LedMode = 0
    Orange          LedMode = 1
    Green           LedMode = 2
    Both            LedMode = 3
)
```

## Supported Devices

### Motors
- LEGO SPIKE Prime Medium Motor
- LEGO SPIKE Prime Large Motor
- LEGO SPIKE Essential Small Angular Motor
- LEGO Technic Large Motor
- LEGO Technic XL Motor
- LEGO Technic Medium Angular Motor
- LEGO System Medium Motor
- LEGO System Train Motor
- LEGO System Turntable Motor

### Sensors
- LEGO SPIKE Prime Color Sensor
- LEGO SPIKE Prime Ultrasonic Distance Sensor
- LEGO SPIKE Prime Force Sensor
- LEGO SPIKE Essential 3x3 Color Light Matrix
- LEGO WeDo Tilt Sensor
- LEGO WeDo Distance Sensor
- LEGO Color and Distance Sensor
- Button/Touch Sensor
- Simple Lights

## Error Handling

All methods return errors that should be checked:

```go
if err := brick.SetMotorPower(models.PortA, 50); err != nil {
    log.Printf("Failed to set motor power: %v", err)
}
```

Common errors include:
- `"not a motor connected to port X"` - No motor connected to the specified port
- `"not an active motor connected to port X"` - Motor is not an active motor type
- `"mode can be changed only on active sensors"` - Trying to change mode on passive sensor
- Serial communication errors

## Thread Safety

The library is thread-safe and can be used from multiple goroutines. All internal state is protected by mutexes.

## Context Support

Methods that can take time (like `MoveMotorForSeconds`) accept a `context.Context` for cancellation and timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := brick.MoveMotorForSeconds(models.PortA, 5.0, 50, true, ctx); err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Motor movement timed out")
    } else {
        log.Printf("Motor movement failed: %v", err)
    }
}
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

This library is a Go port of the .NET BuildHat library by the .NET Foundation.

