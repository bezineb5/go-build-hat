package buildhat

import "testing"

func TestBuildHatPort_String(t *testing.T) {
	tests := []struct {
		port     BuildHatPort
		expected string
	}{
		{PortA, "A"},
		{PortB, "B"},
		{PortC, "C"},
		{PortD, "D"},
		{BuildHatPort(-1), "Invalid(-1)"},
		{BuildHatPort(99), "Invalid(99)"},
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
		port     BuildHatPort
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
		port  BuildHatPort
		valid bool
	}{
		{PortA, true},
		{PortB, true},
		{PortC, true},
		{PortD, true},
		{BuildHatPort(-1), false},
		{BuildHatPort(NumPorts), false},      // One past the end
		{BuildHatPort(NumPorts + 10), false}, // Way out of range
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
		expected    BuildHatPort
		shouldError bool
	}{
		{"A", PortA, false},
		{"B", PortB, false},
		{"C", PortC, false},
		{"D", PortD, false},
		{"a", BuildHatPort(-1), true},
		{"E", BuildHatPort(-1), true},
		{"", BuildHatPort(-1), true},
		{"AB", BuildHatPort(-1), true},
		{"1", BuildHatPort(-1), true},
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

	expected := []BuildHatPort{PortA, PortB, PortC, PortD}
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
