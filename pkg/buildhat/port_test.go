package buildhat

import "testing"

func TestBuildHatPort_String(t *testing.T) {
	tests := []struct {
		port     Port
		expected string
	}{
		{PortA, "A"},
		{PortB, "B"},
		{PortC, "C"},
		{PortD, "D"},
		{Port(-1), "Invalid(-1)"},
		{Port(99), "Invalid(99)"},
	}

	for _, tt := range tests {
		result := tt.port.String()
		if result != tt.expected {
			t.Errorf("Port %d: String() = %q, want %q", tt.port, result, tt.expected)
		}
	}
}

func TestBuildHatPort_Int(t *testing.T) {
	tests := []struct {
		port     Port
		expected int
	}{
		{PortA, 0},
		{PortB, 1},
		{PortC, 2},
		{PortD, 3},
	}

	for _, tt := range tests {
		result := tt.port.Int()
		if result != tt.expected {
			t.Errorf("Port %s: Int() = %d, want %d", tt.port, result, tt.expected)
		}
	}
}

func TestBuildHatPort_IsValid(t *testing.T) {
	tests := []struct {
		port  Port
		valid bool
	}{
		{PortA, true},
		{PortB, true},
		{PortC, true},
		{PortD, true},
		{Port(-1), false},
		{Port(NumPorts), false},      // One past the end
		{Port(NumPorts + 10), false}, // Way out of range
	}

	for _, tt := range tests {
		result := tt.port.IsValid()
		if result != tt.valid {
			t.Errorf("Port %d: IsValid() = %v, want %v", tt.port, result, tt.valid)
		}
	}
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		input       string
		expected    Port
		shouldError bool
	}{
		{"A", PortA, false},
		{"B", PortB, false},
		{"C", PortC, false},
		{"D", PortD, false},
		{"a", Port(-1), true},
		{"E", Port(-1), true},
		{"", Port(-1), true},
		{"AB", Port(-1), true},
		{"1", Port(-1), true},
	}

	for _, tt := range tests {
		result, err := ParsePort(tt.input)
		if tt.shouldError {
			if err == nil {
				t.Errorf("ParsePort(%q): expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParsePort(%q): unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParsePort(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		}
	}
}

func TestAllPorts(t *testing.T) {
	ports := AllPorts()

	if len(ports) != NumPorts {
		t.Errorf("AllPorts() returned %d ports, want %d", len(ports), NumPorts)
	}

	expected := []Port{PortA, PortB, PortC, PortD}
	for i, port := range ports {
		if port != expected[i] {
			t.Errorf("AllPorts()[%d] = %v, want %v", i, port, expected[i])
		}
	}
}

func TestPortRoundTrip(t *testing.T) {
	// Test that String() and ParsePort() are inverses
	for _, port := range AllPorts() {
		s := port.String()
		parsed, err := ParsePort(s)
		if err != nil {
			t.Errorf("Port %v: ParsePort(%q) failed: %v", port, s, err)
		}
		if parsed != port {
			t.Errorf("Port %v: String() = %q, ParsePort(%q) = %v", port, s, s, parsed)
		}
	}
}
