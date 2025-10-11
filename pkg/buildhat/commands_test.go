package buildhat

import (
	"testing"
)

func TestSimpleCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{"Help", Help(), "help"},
		{"Version", Version(), "version"},
		{"List", List(), "list"},
		{"Vin", Vin(), "vin"},
		{"ClearFaults", ClearFaults(), "clear_faults"},
		{"Coast", Coast(), "coast"},
		{"PWM", PWM(), "pwm"},
		{"Off", Off(), "off"},
		{"On", On(), "on"},
		{"Signature", Signature(), "signature"},
		{"Clear", Clear(), "clear"},
		{"Preset", Preset(), "preset"},
		{"Reboot", Reboot(), "reboot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPortCommand(t *testing.T) {
	tests := []struct {
		port     int
		expected string
	}{
		{0, "port 0"},
		{1, "port 1"},
		{2, "port 2"},
		{3, "port 3"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cmd := SelectPort(Port(tt.port))
			result := cmd.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestEchoCommand(t *testing.T) {
	tests := []struct {
		enable   bool
		expected string
	}{
		{false, "echo 0"},
		{true, "echo 1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cmd := Echo(tt.enable)
			result := cmd.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLEDModeCommand(t *testing.T) {
	tests := []struct {
		mode     LEDMode
		expected string
	}{
		{LEDModeAuto, "ledmode -1"},
		{LEDModeOff, "ledmode 0"},
		{LEDModeOrange, "ledmode 1"},
		{LEDModeGreen, "ledmode 2"},
		{LEDModeOrangeAndGreen, "ledmode 3"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cmd := LEDModeCmd(tt.mode)
			result := cmd.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPLimitCommand(t *testing.T) {
	tests := []struct {
		limit    float64
		expected string
	}{
		{0.0, "plimit 0"},
		{0.5, "plimit 0.5"},
		{1.0, "plimit 1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cmd := PLimit(tt.limit)
			result := cmd.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBiasCommand(t *testing.T) {
	tests := []struct {
		bias     float64
		expected string
	}{
		{0.0, "bias 0"},
		{0.4, "bias 0.4"},
		{1.0, "bias 1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cmd := Bias(tt.bias)
			result := cmd.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSetConstantCommand(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{0.0, "set 0"},
		{0.5, "set 0.5"},
		{1.0, "set 1"},
		{-1.0, "set -1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cmd := SetConstant(tt.value)
			result := cmd.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSetWaveformCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{
			"SquareWave",
			SetSquareWave(0, 1, 10, 0),
			"set square 0 1 10 0",
		},
		{
			"SineWave",
			SetSineWave(-1, 1, 5, 0.5),
			"set sine -1 1 5 0.5",
		},
		{
			"TriangleWave",
			SetTriangleWave(0, 1, 10, 0),
			"set triangle 0 1 10 0",
		},
		{
			"Pulse",
			SetPulse(50, 0, 2.5),
			"set pulse 50.000000 0.0 2.500000 0",
		},
		{
			"Ramp",
			SetRamp(0, 360, 3),
			"set ramp 0.000000 360.000000 3.000000 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPIDCommand(t *testing.T) {
	cmd := PID(0, 0, 1, DataFormatS4, 0.0027777778, 0, 5, 0, 0.1, 3, 0.01)
	expected := "pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01"
	result := cmd.CommandString()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestPIDDiffCommand(t *testing.T) {
	cmd := PIDDiff(0, 0, 5, DataFormatS2, 0.0027777778, 1, 0, 2.5, 0, 0.4, 0.01)
	expected := "pid_diff 0 0 5 s2 0.0027777778 1 0 2.5 0 0.4 0.01"
	result := cmd.CommandString()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSelectCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{"Deselect", SelectDeselect(), "select"},
		{"SelectMode", Select(0), "select 0"},
		{"SelectFormatted", SelectFormatted(1, 0, DataFormatS4), "select 1 0 s4"},
		{"SelectOnceDeselect", SelectOnceDeselect(), "selonce"},
		{"SelectOnce", SelectOnce(2), "selonce 2"},
		{"SelectOnceFormatted", SelectOnceFormatted(3, 5, DataFormatS2), "selonce 3 5 s2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCombiCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{
			"Deconfigure",
			CombiDeconfigure(0),
			"combi 0",
		},
		{
			"SingleMode",
			Combi(0, ModeDataset{Mode: 1, Offset: 0}),
			"combi 0 1 0",
		},
		{
			"MultipleModes",
			Combi(0,
				ModeDataset{Mode: 1, Offset: 0},
				ModeDataset{Mode: 2, Offset: 0},
				ModeDataset{Mode: 3, Offset: 0},
			),
			"combi 0 1 0 2 0 3 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWriteCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{
			"Write1",
			Write1(0xc2, 0x12, 0x23),
			"write1 c2 12 23",
		},
		{
			"Write2",
			Write2(0xaa, 0xbb),
			"write2 aa bb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCompoundCommand(t *testing.T) {
	cmd := Compound(
		SelectPort(0),
		PLimit(1),
		SetConstant(-1),
	)
	expected := "port 0 ; plimit 1 ; set -1"
	result := cmd.CommandString()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestExtendedCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{"SelRate", SelRate(10), "selrate 10"},
		{"PortPLimit", PortPLimit(0.7), "port_plimit 0.70"},
		{"PWMParams", PWMParams(0.65, 0.01), "pwmparams 0.65 0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFirmwareCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected string
	}{
		{"Load", Load(1024, 12345), "load 1024 12345"},
		{"SignatureLoad", SignatureLoad(256), "signature 256"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.CommandString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestComplexMotorCommand(t *testing.T) {
	// Test a realistic motor command like what's used in motor.go
	cmd := Compound(
		SelectPort(0),
		Select(0),
		SelRate(10),
		PID(0, 0, 1, DataFormatS4, 0.0027777778, 0, 5, 0, 0.1, 3, 0.01),
		SetRamp(0, 360, 3),
	)

	expected := "port 0 ; select 0 ; selrate 10 ; pid 0 0 1 s4 0.0027777778 0 5 0 0.1 3 0.01 ; set ramp 0.000000 360.000000 3.000000 0"
	result := cmd.CommandString()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
