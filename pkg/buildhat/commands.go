package buildhat

import (
	"fmt"
	"strings"
)

// Command represents a BuildHAT serial command
type Command interface {
	// String returns the command string to send to the BuildHAT
	String() string
}

// ======== Simple Commands (no parameters) ========

// HelpCommand requests help/command synopsis
type HelpCommand struct{}

func (c HelpCommand) String() string { return "help" }

// Help creates a help command
func Help() Command { return HelpCommand{} }

// VersionCommand requests firmware version
type VersionCommand struct{}

func (c VersionCommand) String() string { return "version" }

// Version creates a version command
func Version() Command { return VersionCommand{} }

// ListCommand requests list of connected devices
type ListCommand struct{}

func (c ListCommand) String() string { return "list" }

// List creates a list command
func List() Command { return ListCommand{} }

// VinCommand requests input voltage
type VinCommand struct{}

func (c VinCommand) String() string { return "vin" }

// Vin creates a vin command
func Vin() Command { return VinCommand{} }

// ClearFaultsCommand clears motor power faults
type ClearFaultsCommand struct{}

func (c ClearFaultsCommand) String() string { return "clear_faults" }

// ClearFaults creates a clear_faults command
func ClearFaults() Command { return ClearFaultsCommand{} }

// CoastCommand switches motor driver to coast mode
type CoastCommand struct{}

func (c CoastCommand) String() string { return "coast" }

// Coast creates a coast command
func Coast() Command { return CoastCommand{} }

// PWMCommand switches controller to direct PWM mode
type PWMCommand struct{}

func (c PWMCommand) String() string { return "pwm" }

// PWM creates a pwm command
func PWM() Command { return PWMCommand{} }

// OffCommand turns off motor (pwm; set 0)
type OffCommand struct{}

func (c OffCommand) String() string { return "off" }

// Off creates an off command
func Off() Command { return OffCommand{} }

// OnCommand turns on motor full power (pwm; set 1)
type OnCommand struct{}

func (c OnCommand) String() string { return "on" }

// On creates an on command
func On() Command { return OnCommand{} }

// SignatureCommand requests signature (not normally needed)
type SignatureCommand struct{}

func (c SignatureCommand) String() string { return "signature" }

// Signature creates a signature command
func Signature() Command { return SignatureCommand{} }

// ======== Simple Parameterized Commands ========

// PortCommand sets the current port
type PortCommand struct {
	PortNum int
}

func (c PortCommand) String() string {
	return fmt.Sprintf("port %d", c.PortNum)
}

// SelectPort creates a port command
func SelectPort(port int) Command {
	return PortCommand{PortNum: port}
}

// EchoCommand enables/disables echo
type EchoCommand struct {
	Enable bool
}

func (c EchoCommand) String() string {
	val := 0
	if c.Enable {
		val = 1
	}
	return fmt.Sprintf("echo %d", val)
}

// Echo creates an echo command
func Echo(enable bool) Command {
	return EchoCommand{Enable: enable}
}

// LEDMode represents the LED behavior mode
type LEDMode int

const (
	LEDModeAuto           LEDMode = -1 // LEDs depend on voltage
	LEDModeOff            LEDMode = 0  // LEDs off
	LEDModeOrange         LEDMode = 1  // Orange LED
	LEDModeGreen          LEDMode = 2  // Green LED
	LEDModeOrangeAndGreen LEDMode = 3  // Both LEDs
)

// LEDModeCommand sets LED behavior
type LEDModeCommand struct {
	Mode LEDMode
}

func (c LEDModeCommand) String() string {
	return fmt.Sprintf("ledmode %d", c.Mode)
}

// LEDModeCmd creates an ledmode command
func LEDModeCmd(mode LEDMode) Command {
	return LEDModeCommand{Mode: mode}
}

// PLimitCommand sets global power limit
type PLimitCommand struct {
	Limit float64
}

func (c PLimitCommand) String() string {
	return fmt.Sprintf("plimit %g", c.Limit)
}

// PLimit creates a plimit command
func PLimit(limit float64) Command {
	return PLimitCommand{Limit: limit}
}

// BiasCommand sets bias for motor drive
type BiasCommand struct {
	Bias float64
}

func (c BiasCommand) String() string {
	return fmt.Sprintf("bias %g", c.Bias)
}

// Bias creates a bias command
func Bias(bias float64) Command {
	return BiasCommand{Bias: bias}
}

// DebugCommand sets debug mode
type DebugCommand struct {
	DebugCode int
}

func (c DebugCommand) String() string {
	return fmt.Sprintf("debug %d", c.DebugCode)
}

// Debug creates a debug command
func Debug(code int) Command {
	return DebugCommand{DebugCode: code}
}

// ======== Set Command (Setpoint Control) ========

// SetCommand configures the setpoint for a controller
type SetCommand struct {
	setpoint setpoint
}

type setpoint interface {
	String() string
}

// ConstantSetpoint represents a constant setpoint value
type ConstantSetpoint struct {
	Value float64
}

func (s ConstantSetpoint) String() string {
	return fmt.Sprintf("%g", s.Value)
}

// SquareWaveSetpoint represents a square wave setpoint
type SquareWaveSetpoint struct {
	Min    float64
	Max    float64
	Period float64
	Phase  float64
}

func (s SquareWaveSetpoint) String() string {
	return fmt.Sprintf("square %g %g %g %g", s.Min, s.Max, s.Period, s.Phase)
}

// SineWaveSetpoint represents a sine wave setpoint
type SineWaveSetpoint struct {
	Min    float64
	Max    float64
	Period float64
	Phase  float64
}

func (s SineWaveSetpoint) String() string {
	return fmt.Sprintf("sine %g %g %g %g", s.Min, s.Max, s.Period, s.Phase)
}

// TriangleWaveSetpoint represents a triangle wave setpoint
type TriangleWaveSetpoint struct {
	Min    float64
	Max    float64
	Period float64
	Phase  float64
}

func (s TriangleWaveSetpoint) String() string {
	return fmt.Sprintf("triangle %g %g %g %g", s.Min, s.Max, s.Period, s.Phase)
}

// PulseSetpoint represents a pulse setpoint
type PulseSetpoint struct {
	DuringValue float64
	AfterValue  float64
	Duration    float64
}

func (s PulseSetpoint) String() string {
	// Special case: format 0.0 as "0.0" not "0"
	afterStr := fmt.Sprintf("%g", s.AfterValue)
	if s.AfterValue == 0 {
		afterStr = "0.0"
	}
	return fmt.Sprintf("pulse %f %s %f 0", s.DuringValue, afterStr, s.Duration)
}

// RampSetpoint represents a ramp setpoint
type RampSetpoint struct {
	StartValue float64
	EndValue   float64
	Duration   float64
}

func (s RampSetpoint) String() string {
	return fmt.Sprintf("ramp %f %f %f 0", s.StartValue, s.EndValue, s.Duration)
}

func (c SetCommand) String() string {
	return fmt.Sprintf("set %s", c.setpoint.String())
}

// SetConstant creates a set command with constant value
func SetConstant(value float64) Command {
	return SetCommand{setpoint: ConstantSetpoint{Value: value}}
}

// SetConstantFormatted creates a set command with a specific float format
func SetConstantFormatted(value float64, format string) Command {
	return SetCommand{setpoint: &formattedConstant{Value: value, Format: format}}
}

type formattedConstant struct {
	Value  float64
	Format string
}

func (f *formattedConstant) String() string {
	return fmt.Sprintf(f.Format, f.Value)
}

// SetSquareWave creates a set command with square wave
func SetSquareWave(minVal, maxVal, period, phase float64) Command {
	return SetCommand{setpoint: SquareWaveSetpoint{Min: minVal, Max: maxVal, Period: period, Phase: phase}}
}

// SetSineWave creates a set command with sine wave
func SetSineWave(minVal, maxVal, period, phase float64) Command {
	return SetCommand{setpoint: SineWaveSetpoint{Min: minVal, Max: maxVal, Period: period, Phase: phase}}
}

// SetTriangleWave creates a set command with triangle wave
func SetTriangleWave(minVal, maxVal, period, phase float64) Command {
	return SetCommand{setpoint: TriangleWaveSetpoint{Min: minVal, Max: maxVal, Period: period, Phase: phase}}
}

// SetPulse creates a set command with pulse
func SetPulse(duringValue, afterValue, duration float64) Command {
	return SetCommand{setpoint: PulseSetpoint{DuringValue: duringValue, AfterValue: afterValue, Duration: duration}}
}

// SetRamp creates a set command with ramp
func SetRamp(startValue, endValue, duration float64) Command {
	return SetCommand{setpoint: RampSetpoint{StartValue: startValue, EndValue: endValue, Duration: duration}}
}

// ======== PID Command ========

// DataFormat represents the format of process variable data
type DataFormat string

const (
	DataFormatU1 DataFormat = "u1" // unsigned byte
	DataFormatS1 DataFormat = "s1" // signed byte
	DataFormatU2 DataFormat = "u2" // unsigned short
	DataFormatS2 DataFormat = "s2" // signed short
	DataFormatU4 DataFormat = "u4" // unsigned int
	DataFormatS4 DataFormat = "s4" // signed int
	DataFormatF4 DataFormat = "f4" // float
)

// PIDCommand switches controller to PID mode
type PIDCommand struct {
	PVPort   int        // port to fetch process variable from
	PVMode   int        // mode to fetch process variable from
	PVOffset int        // byte offset into mode
	PVFormat DataFormat // format of process variable
	PVScale  float64    // multiplicative scale factor
	PVUnwrap int        // 0=no unwrapping, otherwise modulo for phase unwrap
	Kp       float64    // proportional gain
	Ki       float64    // integral gain
	Kd       float64    // differential gain
	Windup   float64    // integral windup limit
	Bias     float64    // bias value (undocumented but used in practice)
}

func (c PIDCommand) String() string {
	return fmt.Sprintf("pid %d %d %d %s %g %d %g %g %g %g %g",
		c.PVPort, c.PVMode, c.PVOffset, c.PVFormat,
		c.PVScale, c.PVUnwrap, c.Kp, c.Ki, c.Kd, c.Windup, c.Bias)
}

// PID creates a PID command
func PID(pvPort, pvMode, pvOffset int, pvFormat DataFormat, pvScale float64,
	pvUnwrap int, kp, ki, kd, windup, bias float64) Command {
	return PIDCommand{
		PVPort:   pvPort,
		PVMode:   pvMode,
		PVOffset: pvOffset,
		PVFormat: pvFormat,
		PVScale:  pvScale,
		PVUnwrap: pvUnwrap,
		Kp:       kp,
		Ki:       ki,
		Kd:       kd,
		Windup:   windup,
		Bias:     bias,
	}
}

// PIDDiffCommand switches controller to PID differential mode (for velocity control)
type PIDDiffCommand struct {
	PVPort   int        // port to fetch process variable from
	PVMode   int        // mode to fetch process variable from
	PVOffset int        // byte offset into mode
	PVFormat DataFormat // format of process variable
	PVScale  float64    // multiplicative scale factor
	PVUnwrap int        // 0=no unwrapping, otherwise modulo for phase unwrap
	Kp       float64    // proportional gain
	Ki       float64    // integral gain
	Kd       float64    // differential gain
	Windup   float64    // integral windup limit
	Bias     float64    // bias value (undocumented but used in practice)
}

func (c PIDDiffCommand) String() string {
	return fmt.Sprintf("pid_diff %d %d %d %s %g %d %g %g %g %g %g",
		c.PVPort, c.PVMode, c.PVOffset, c.PVFormat,
		c.PVScale, c.PVUnwrap, c.Kp, c.Ki, c.Kd, c.Windup, c.Bias)
}

// PIDDiff creates a PID differential command
func PIDDiff(pvPort, pvMode, pvOffset int, pvFormat DataFormat, pvScale float64,
	pvUnwrap int, kp, ki, kd, windup, bias float64) Command {
	return PIDDiffCommand{
		PVPort:   pvPort,
		PVMode:   pvMode,
		PVOffset: pvOffset,
		PVFormat: pvFormat,
		PVScale:  pvScale,
		PVUnwrap: pvUnwrap,
		Kp:       kp,
		Ki:       ki,
		Kd:       kd,
		Windup:   windup,
		Bias:     bias,
	}
}

// ======== Select Commands ========

// SelectCommand selects a mode on the current port
type SelectCommand struct {
	Mode   *int        // nil means deselect
	Offset *int        // nil for raw hex output
	Format *DataFormat // nil for raw hex output
}

func (c SelectCommand) String() string {
	if c.Mode == nil {
		return "select"
	}
	if c.Offset == nil || c.Format == nil {
		return fmt.Sprintf("select %d", *c.Mode)
	}
	return fmt.Sprintf("select %d %d %s", *c.Mode, *c.Offset, *c.Format)
}

// SelectDeselect creates a select command that deselects any mode
func SelectDeselect() Command {
	return SelectCommand{}
}

// Select creates a select command for raw hex output
func Select(mode int) Command {
	return SelectCommand{Mode: &mode}
}

// SelectFormatted creates a select command with offset and format
func SelectFormatted(mode, offset int, format DataFormat) Command {
	return SelectCommand{Mode: &mode, Offset: &offset, Format: &format}
}

// SelectOnceCommand is like select but outputs once
type SelectOnceCommand struct {
	Mode   *int        // nil means deselect
	Offset *int        // nil for raw hex output
	Format *DataFormat // nil for raw hex output
}

func (c SelectOnceCommand) String() string {
	if c.Mode == nil {
		return "selonce"
	}
	if c.Offset == nil || c.Format == nil {
		return fmt.Sprintf("selonce %d", *c.Mode)
	}
	return fmt.Sprintf("selonce %d %d %s", *c.Mode, *c.Offset, *c.Format)
}

// SelectOnceDeselect creates a selonce command that deselects any mode
func SelectOnceDeselect() Command {
	return SelectOnceCommand{}
}

// SelectOnce creates a selonce command for raw hex output
func SelectOnce(mode int) Command {
	return SelectOnceCommand{Mode: &mode}
}

// SelectOnceFormatted creates a selonce command with offset and format
func SelectOnceFormatted(mode, offset int, format DataFormat) Command {
	return SelectOnceCommand{Mode: &mode, Offset: &offset, Format: &format}
}

// ======== Combi Command ========

// ModeDataset represents a mode and dataset offset pair for combi mode
type ModeDataset struct {
	Mode   int
	Offset int
}

// CombiCommand configures a combi mode
type CombiCommand struct {
	Index    int
	ModeList []ModeDataset // nil/empty means deconfigure
}

func (c CombiCommand) String() string {
	if len(c.ModeList) == 0 {
		return fmt.Sprintf("combi %d", c.Index)
	}

	parts := []string{fmt.Sprintf("combi %d", c.Index)}
	for _, md := range c.ModeList {
		parts = append(parts, fmt.Sprintf("%d %d", md.Mode, md.Offset))
	}
	return strings.Join(parts, " ")
}

// CombiDeconfigure creates a combi command that deconfigures a combi mode
func CombiDeconfigure(index int) Command {
	return CombiCommand{Index: index}
}

// Combi creates a combi command that configures a combi mode
func Combi(index int, modeList ...ModeDataset) Command {
	return CombiCommand{Index: index, ModeList: modeList}
}

// ======== Write Commands ========

// Write1Command writes bytes with 1-byte header
type Write1Command struct {
	Bytes []byte
}

func (c Write1Command) String() string {
	hexParts := make([]string, len(c.Bytes))
	for i, b := range c.Bytes {
		hexParts[i] = fmt.Sprintf("%x", b)
	}
	return fmt.Sprintf("write1 %s", strings.Join(hexParts, " "))
}

// Write1 creates a write1 command
func Write1(bytes ...byte) Command {
	return Write1Command{Bytes: bytes}
}

// Write2Command writes bytes with 2-byte header
type Write2Command struct {
	Bytes []byte
}

func (c Write2Command) String() string {
	hexParts := make([]string, len(c.Bytes))
	for i, b := range c.Bytes {
		hexParts[i] = fmt.Sprintf("%x", b)
	}
	return fmt.Sprintf("write2 %s", strings.Join(hexParts, " "))
}

// Write2 creates a write2 command
func Write2(bytes ...byte) Command {
	return Write2Command{Bytes: bytes}
}

// ======== Compound Commands ========

// CompoundCommand allows multiple commands on one line
type CompoundCommand struct {
	Commands []Command
}

func (c CompoundCommand) String() string {
	parts := make([]string, len(c.Commands))
	for i, cmd := range c.Commands {
		parts[i] = cmd.String()
	}
	return strings.Join(parts, " ; ")
}

// Compound creates a compound command from multiple commands
func Compound(commands ...Command) Command {
	return CompoundCommand{Commands: commands}
}

// ======== Extended Commands (not in protocol.md but used in practice) ========

// SelRateCommand sets the selection rate (frequency of sensor readings)
type SelRateCommand struct {
	Rate int
}

func (c SelRateCommand) String() string {
	return fmt.Sprintf("selrate %d", c.Rate)
}

// SelRate creates a selrate command
func SelRate(rate int) Command {
	return SelRateCommand{Rate: rate}
}

// PresetCommand presets the motor position
type PresetCommand struct{}

func (c PresetCommand) String() string { return "preset" }

// Preset creates a preset command
func Preset() Command { return PresetCommand{} }

// PortPLimitCommand sets power limit for specific port
type PortPLimitCommand struct {
	Limit float64
}

func (c PortPLimitCommand) String() string {
	return fmt.Sprintf("port_plimit %.2f", c.Limit)
}

// PortPLimit creates a port_plimit command
func PortPLimit(limit float64) Command {
	return PortPLimitCommand{Limit: limit}
}

// PWMParamsCommand sets PWM threshold and minimum PWM parameters
type PWMParamsCommand struct {
	PWMThresh float64
	MinPWM    float64
}

func (c PWMParamsCommand) String() string {
	return fmt.Sprintf("pwmparams %.2f %.2f", c.PWMThresh, c.MinPWM)
}

// PWMParams creates a pwmparams command
func PWMParams(pwmThresh, minPWM float64) Command {
	return PWMParamsCommand{PWMThresh: pwmThresh, MinPWM: minPWM}
}

// ======== Firmware Management Commands ========

// ClearCommand clears the command buffer
type ClearCommand struct{}

func (c ClearCommand) String() string { return "clear" }

// Clear creates a clear command
func Clear() Command { return ClearCommand{} }

// LoadCommand loads firmware
type LoadCommand struct {
	Length   int
	Checksum int
}

func (c LoadCommand) String() string {
	return fmt.Sprintf("load %d %d", c.Length, c.Checksum)
}

// Load creates a load command
func Load(length, checksum int) Command {
	return LoadCommand{Length: length, Checksum: checksum}
}

// SignatureLoadCommand loads signature
type SignatureLoadCommand struct {
	Length int
}

func (c SignatureLoadCommand) String() string {
	return fmt.Sprintf("signature %d", c.Length)
}

// SignatureLoad creates a signature load command
func SignatureLoad(length int) Command {
	return SignatureLoadCommand{Length: length}
}

// RebootCommand reboots the BuildHAT
type RebootCommand struct{}

func (c RebootCommand) String() string { return "reboot" }

// Reboot creates a reboot command
func Reboot() Command { return RebootCommand{} }
