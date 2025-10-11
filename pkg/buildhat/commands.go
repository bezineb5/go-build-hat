package buildhat

import (
	"fmt"
	"strings"
)

// Command represents a BuildHAT serial command
type Command interface {
	// CommandString returns the command string to send to the BuildHAT
	CommandString() string
}

// ======== Simple Commands (no parameters) ========

// HelpCommand requests help/command synopsis
type HelpCommand struct{}

func (c HelpCommand) CommandString() string { return "help" }

// Help creates a help command
func Help() Command { return HelpCommand{} }

// VersionCommand requests firmware version
type VersionCommand struct{}

func (c VersionCommand) CommandString() string { return "version" }

// Version creates a version command
func Version() Command { return VersionCommand{} }

// ListCommand requests list of connected devices
type ListCommand struct{}

func (c ListCommand) CommandString() string { return "list" }

// List creates a list command
func List() Command { return ListCommand{} }

// VinCommand requests input voltage
type VinCommand struct{}

func (c VinCommand) CommandString() string { return "vin" }

// Vin creates a vin command
func Vin() Command { return VinCommand{} }

// ClearFaultsCommand clears motor power faults
type ClearFaultsCommand struct{}

func (c ClearFaultsCommand) CommandString() string { return "clear_faults" }

// ClearFaults creates a clear_faults command
func ClearFaults() Command { return ClearFaultsCommand{} }

// CoastCommand switches motor driver to coast mode
type CoastCommand struct{}

func (c CoastCommand) CommandString() string { return "coast" }

// Coast creates a coast command
func Coast() Command { return CoastCommand{} }

// PWMCommand switches controller to direct PWM mode
type PWMCommand struct{}

func (c PWMCommand) CommandString() string { return "pwm" }

// PWM creates a pwm command
func PWM() Command { return PWMCommand{} }

// OffCommand turns off motor (pwm; set 0)
type OffCommand struct{}

func (c OffCommand) CommandString() string { return "off" }

// Off creates an off command
func Off() Command { return OffCommand{} }

// OnCommand turns on motor full power (pwm; set 1)
type OnCommand struct{}

func (c OnCommand) CommandString() string { return "on" }

// On creates an on command
func On() Command { return OnCommand{} }

// SignatureCommand requests signature (not normally needed)
type SignatureCommand struct{}

func (c SignatureCommand) CommandString() string { return "signature" }

// Signature creates a signature command
func Signature() Command { return SignatureCommand{} }

// ======== Simple Parameterized Commands ========

// PortCommand sets the current port
type PortCommand struct {
	port Port
}

func (c *PortCommand) CommandString() string {
	return fmt.Sprintf("port %d", c.port.Int())
}

// SelectPort creates a port command
func SelectPort(port Port) Command {
	return &PortCommand{port: port}
}

// EchoCommand enables/disables echo
type EchoCommand struct {
	enable bool
}

func (c *EchoCommand) CommandString() string {
	val := 0
	if c.enable {
		val = 1
	}
	return fmt.Sprintf("echo %d", val)
}

// Echo creates an echo command
func Echo(enable bool) Command {
	return &EchoCommand{enable: enable}
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
	mode LEDMode
}

func (c *LEDModeCommand) CommandString() string {
	return fmt.Sprintf("ledmode %d", c.mode)
}

// LEDModeCmd creates an ledmode command
func LEDModeCmd(mode LEDMode) Command {
	return &LEDModeCommand{mode: mode}
}

// PLimitCommand sets global power limit
type PLimitCommand struct {
	limit float64
}

func (c *PLimitCommand) CommandString() string {
	return fmt.Sprintf("plimit %g", c.limit)
}

// PLimit creates a plimit command
func PLimit(limit float64) Command {
	return &PLimitCommand{limit: limit}
}

// BiasCommand sets bias for motor drive
type BiasCommand struct {
	bias float64
}

func (c *BiasCommand) CommandString() string {
	return fmt.Sprintf("bias %g", c.bias)
}

// Bias creates a bias command
func Bias(bias float64) Command {
	return &BiasCommand{bias: bias}
}

// DebugCommand sets debug mode
type DebugCommand struct {
	debugCode int
}

func (c *DebugCommand) CommandString() string {
	return fmt.Sprintf("debug %d", c.debugCode)
}

// Debug creates a debug command
func Debug(code int) Command {
	return &DebugCommand{debugCode: code}
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
	value float64
}

func (s *ConstantSetpoint) String() string {
	return fmt.Sprintf("%g", s.value)
}

// SquareWaveSetpoint represents a square wave setpoint
type SquareWaveSetpoint struct {
	min    float64
	max    float64
	period float64
	phase  float64
}

func (s *SquareWaveSetpoint) String() string {
	return fmt.Sprintf("square %g %g %g %g", s.min, s.max, s.period, s.phase)
}

// SineWaveSetpoint represents a sine wave setpoint
type SineWaveSetpoint struct {
	min    float64
	max    float64
	period float64
	phase  float64
}

func (s *SineWaveSetpoint) String() string {
	return fmt.Sprintf("sine %g %g %g %g", s.min, s.max, s.period, s.phase)
}

// TriangleWaveSetpoint represents a triangle wave setpoint
type TriangleWaveSetpoint struct {
	min    float64
	max    float64
	period float64
	phase  float64
}

func (s *TriangleWaveSetpoint) String() string {
	return fmt.Sprintf("triangle %g %g %g %g", s.min, s.max, s.period, s.phase)
}

// PulseSetpoint represents a pulse setpoint
type PulseSetpoint struct {
	duringValue float64
	afterValue  float64
	duration    float64
}

func (s *PulseSetpoint) String() string {
	// Special case: format 0.0 as "0.0" not "0"
	afterStr := fmt.Sprintf("%g", s.afterValue)
	if s.afterValue == 0 {
		afterStr = "0.0"
	}
	return fmt.Sprintf("pulse %f %s %f 0", s.duringValue, afterStr, s.duration)
}

// RampSetpoint represents a ramp setpoint
type RampSetpoint struct {
	startValue float64
	endValue   float64
	duration   float64
}

func (s *RampSetpoint) String() string {
	return fmt.Sprintf("ramp %f %f %f 0", s.startValue, s.endValue, s.duration)
}

func (c SetCommand) CommandString() string {
	return fmt.Sprintf("set %s", c.setpoint.String())
}

// SetConstant creates a set command with constant value
func SetConstant(value float64) Command {
	return SetCommand{setpoint: &ConstantSetpoint{value: value}}
}

// SetConstantFormatted creates a set command with a specific float format
func SetConstantFormatted(value float64, format string) Command {
	return SetCommand{setpoint: &formattedConstant{value: value, format: format}}
}

type formattedConstant struct {
	value  float64
	format string
}

func (f *formattedConstant) String() string {
	return fmt.Sprintf(f.format, f.value)
}

// SetSquareWave creates a set command with square wave
func SetSquareWave(minVal, maxVal, period, phase float64) Command {
	return SetCommand{setpoint: &SquareWaveSetpoint{min: minVal, max: maxVal, period: period, phase: phase}}
}

// SetSineWave creates a set command with sine wave
func SetSineWave(minVal, maxVal, period, phase float64) Command {
	return SetCommand{setpoint: &SineWaveSetpoint{min: minVal, max: maxVal, period: period, phase: phase}}
}

// SetTriangleWave creates a set command with triangle wave
func SetTriangleWave(minVal, maxVal, period, phase float64) Command {
	return SetCommand{setpoint: &TriangleWaveSetpoint{min: minVal, max: maxVal, period: period, phase: phase}}
}

// SetPulse creates a set command with pulse
func SetPulse(duringValue, afterValue, duration float64) Command {
	return SetCommand{setpoint: &PulseSetpoint{duringValue: duringValue, afterValue: afterValue, duration: duration}}
}

// SetRamp creates a set command with ramp
func SetRamp(startValue, endValue, duration float64) Command {
	return SetCommand{setpoint: &RampSetpoint{startValue: startValue, endValue: endValue, duration: duration}}
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
	pvPort   int        // port to fetch process variable from
	pvMode   int        // mode to fetch process variable from
	pvOffset int        // byte offset into mode
	pvFormat DataFormat // format of process variable
	pvScale  float64    // multiplicative scale factor
	pvUnwrap int        // 0=no unwrapping, otherwise modulo for phase unwrap
	kp       float64    // proportional gain
	ki       float64    // integral gain
	kd       float64    // differential gain
	windup   float64    // integral windup limit
	bias     float64    // bias value (undocumented but used in practice)
}

func (c *PIDCommand) CommandString() string {
	return fmt.Sprintf("pid %d %d %d %s %g %d %g %g %g %g %g",
		c.pvPort, c.pvMode, c.pvOffset, c.pvFormat,
		c.pvScale, c.pvUnwrap, c.kp, c.ki, c.kd, c.windup, c.bias)
}

// PID creates a PID command
func PID(pvPort, pvMode, pvOffset int, pvFormat DataFormat, pvScale float64,
	pvUnwrap int, kp, ki, kd, windup, bias float64) Command {
	return &PIDCommand{
		pvPort:   pvPort,
		pvMode:   pvMode,
		pvOffset: pvOffset,
		pvFormat: pvFormat,
		pvScale:  pvScale,
		pvUnwrap: pvUnwrap,
		kp:       kp,
		ki:       ki,
		kd:       kd,
		windup:   windup,
		bias:     bias,
	}
}

// PIDDiffCommand switches controller to PID differential mode (for velocity control)
type PIDDiffCommand struct {
	pvPort   int        // port to fetch process variable from
	pvMode   int        // mode to fetch process variable from
	pvOffset int        // byte offset into mode
	pvFormat DataFormat // format of process variable
	pvScale  float64    // multiplicative scale factor
	pvUnwrap int        // 0=no unwrapping, otherwise modulo for phase unwrap
	kp       float64    // proportional gain
	ki       float64    // integral gain
	kd       float64    // differential gain
	windup   float64    // integral windup limit
	bias     float64    // bias value (undocumented but used in practice)
}

func (c *PIDDiffCommand) CommandString() string {
	return fmt.Sprintf("pid_diff %d %d %d %s %g %d %g %g %g %g %g",
		c.pvPort, c.pvMode, c.pvOffset, c.pvFormat,
		c.pvScale, c.pvUnwrap, c.kp, c.ki, c.kd, c.windup, c.bias)
}

// PIDDiff creates a PID differential command
func PIDDiff(pvPort, pvMode, pvOffset int, pvFormat DataFormat, pvScale float64,
	pvUnwrap int, kp, ki, kd, windup, bias float64) Command {
	return &PIDDiffCommand{
		pvPort:   pvPort,
		pvMode:   pvMode,
		pvOffset: pvOffset,
		pvFormat: pvFormat,
		pvScale:  pvScale,
		pvUnwrap: pvUnwrap,
		kp:       kp,
		ki:       ki,
		kd:       kd,
		windup:   windup,
		bias:     bias,
	}
}

// ======== Select Commands ========

// SelectCommand selects a mode on the current port
type SelectCommand struct {
	mode   *int        // nil means deselect
	offset *int        // nil for raw hex output
	format *DataFormat // nil for raw hex output
}

func (c *SelectCommand) CommandString() string {
	if c.mode == nil {
		return "select"
	}
	if c.offset == nil || c.format == nil {
		return fmt.Sprintf("select %d", *c.mode)
	}
	return fmt.Sprintf("select %d %d %s", *c.mode, *c.offset, *c.format)
}

// SelectDeselect creates a select command that deselects any mode
func SelectDeselect() Command {
	return &SelectCommand{}
}

// Select creates a select command for raw hex output
func Select(mode int) Command {
	return &SelectCommand{mode: &mode}
}

// SelectFormatted creates a select command with offset and format
func SelectFormatted(mode, offset int, format DataFormat) Command {
	return &SelectCommand{mode: &mode, offset: &offset, format: &format}
}

// SelectOnceCommand is like select but outputs once
type SelectOnceCommand struct {
	mode   *int        // nil means deselect
	offset *int        // nil for raw hex output
	format *DataFormat // nil for raw hex output
}

func (c *SelectOnceCommand) CommandString() string {
	if c.mode == nil {
		return "selonce"
	}
	if c.offset == nil || c.format == nil {
		return fmt.Sprintf("selonce %d", *c.mode)
	}
	return fmt.Sprintf("selonce %d %d %s", *c.mode, *c.offset, *c.format)
}

// SelectOnceDeselect creates a selonce command that deselects any mode
func SelectOnceDeselect() Command {
	return &SelectOnceCommand{}
}

// SelectOnce creates a selonce command for raw hex output
func SelectOnce(mode int) Command {
	return &SelectOnceCommand{mode: &mode}
}

// SelectOnceFormatted creates a selonce command with offset and format
func SelectOnceFormatted(mode, offset int, format DataFormat) Command {
	return &SelectOnceCommand{mode: &mode, offset: &offset, format: &format}
}

// ======== Combi Command ========

// ModeDataset represents a mode and dataset offset pair for combi mode
type ModeDataset struct {
	mode   int
	offset int
}

// NewModeDataset creates a new ModeDataset
func NewModeDataset(mode, offset int) ModeDataset {
	return ModeDataset{mode: mode, offset: offset}
}

// CombiCommand configures a combi mode
type CombiCommand struct {
	index    int
	modeList []ModeDataset // nil/empty means deconfigure
}

func (c *CombiCommand) CommandString() string {
	if len(c.modeList) == 0 {
		return fmt.Sprintf("combi %d", c.index)
	}

	parts := []string{fmt.Sprintf("combi %d", c.index)}
	for _, md := range c.modeList {
		parts = append(parts, fmt.Sprintf("%d %d", md.mode, md.offset))
	}
	return strings.Join(parts, " ")
}

// CombiDeconfigure creates a combi command that deconfigures a combi mode
func CombiDeconfigure(index int) Command {
	return &CombiCommand{index: index}
}

// Combi creates a combi command that configures a combi mode
func Combi(index int, modeList ...ModeDataset) Command {
	return &CombiCommand{index: index, modeList: modeList}
}

// ======== Write Commands ========

// Write1Command writes bytes with 1-byte header
type Write1Command struct {
	bytes []byte
}

func (c *Write1Command) CommandString() string {
	hexParts := make([]string, len(c.bytes))
	for i, b := range c.bytes {
		hexParts[i] = fmt.Sprintf("%x", b)
	}
	return fmt.Sprintf("write1 %s", strings.Join(hexParts, " "))
}

// Write1 creates a write1 command
func Write1(bytes ...byte) Command {
	return &Write1Command{bytes: bytes}
}

// Write2Command writes bytes with 2-byte header
type Write2Command struct {
	bytes []byte
}

func (c *Write2Command) CommandString() string {
	hexParts := make([]string, len(c.bytes))
	for i, b := range c.bytes {
		hexParts[i] = fmt.Sprintf("%x", b)
	}
	return fmt.Sprintf("write2 %s", strings.Join(hexParts, " "))
}

// Write2 creates a write2 command
func Write2(bytes ...byte) Command {
	return &Write2Command{bytes: bytes}
}

// ======== Compound Commands ========

// CompoundCommand allows multiple commands on one line
type CompoundCommand struct {
	commands []Command
}

func (c *CompoundCommand) CommandString() string {
	parts := make([]string, len(c.commands))
	for i, cmd := range c.commands {
		parts[i] = cmd.CommandString()
	}
	return strings.Join(parts, " ; ")
}

// Compound creates a compound command from multiple commands
func Compound(commands ...Command) Command {
	return &CompoundCommand{commands: commands}
}

// ======== Extended Commands (not in protocol.md but used in practice) ========

// SelRateCommand sets the selection rate (frequency of sensor readings)
type SelRateCommand struct {
	rate int
}

func (c *SelRateCommand) CommandString() string {
	return fmt.Sprintf("selrate %d", c.rate)
}

// SelRate creates a selrate command
func SelRate(rate int) Command {
	return &SelRateCommand{rate: rate}
}

// PresetCommand presets the motor position
type PresetCommand struct{}

func (c PresetCommand) CommandString() string { return "preset" }

// Preset creates a preset command
func Preset() Command { return PresetCommand{} }

// PortPLimitCommand sets power limit for specific port
type PortPLimitCommand struct {
	limit float64
}

func (c *PortPLimitCommand) CommandString() string {
	return fmt.Sprintf("port_plimit %.2f", c.limit)
}

// PortPLimit creates a port_plimit command
func PortPLimit(limit float64) Command {
	return &PortPLimitCommand{limit: limit}
}

// PWMParamsCommand sets PWM threshold and minimum PWM parameters
type PWMParamsCommand struct {
	pwmThresh float64
	minPWM    float64
}

func (c *PWMParamsCommand) CommandString() string {
	return fmt.Sprintf("pwmparams %.2f %.2f", c.pwmThresh, c.minPWM)
}

// PWMParams creates a pwmparams command
func PWMParams(pwmThresh, minPWM float64) Command {
	return &PWMParamsCommand{pwmThresh: pwmThresh, minPWM: minPWM}
}

// ======== Firmware Management Commands ========

// ClearCommand clears the command buffer
type ClearCommand struct{}

func (c ClearCommand) CommandString() string { return "clear" }

// Clear creates a clear command
func Clear() Command { return ClearCommand{} }

// LoadCommand loads firmware
type LoadCommand struct {
	length   int
	checksum int
}

func (c *LoadCommand) CommandString() string {
	return fmt.Sprintf("load %d %d", c.length, c.checksum)
}

// Load creates a load command
func Load(length, checksum int) Command {
	return &LoadCommand{length: length, checksum: checksum}
}

// SignatureLoadCommand loads signature
type SignatureLoadCommand struct {
	length int
}

func (c *SignatureLoadCommand) CommandString() string {
	return fmt.Sprintf("signature %d", c.length)
}

// SignatureLoad creates a signature load command
func SignatureLoad(length int) Command {
	return &SignatureLoadCommand{length: length}
}

// RebootCommand reboots the BuildHAT
type RebootCommand struct{}

func (c RebootCommand) CommandString() string { return "reboot" }

// Reboot creates a reboot command
func Reboot() Command { return RebootCommand{} }
