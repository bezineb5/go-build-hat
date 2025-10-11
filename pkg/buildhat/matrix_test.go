package buildhat

import (
	"fmt"
	"strings"
	"testing"
)

func TestMatrix_SetPixel(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortA)

	err := matrix.SetPixel(1, 1, MatrixTurquoise, 8)
	if err != nil {
		t.Fatalf("SetPixel failed: %v", err)
	}

	// Verify EXACT command: "port 0 ; write1 c2 0 0 0 0 85 0 0 0 0\r"
	// SetPixel(1, 1, 5, 8) sets pixels[1][1] = {Color:5, Brightness:8}
	// Matrix data: 9 pixels starting from [0][0], [0][1], [0][2], [1][0], [1][1]...
	// Pixel [1][1] is at index 4 (0-indexed)
	// Byte value: (brightness << 4) | color = (8 << 4) | 5 = 128 + 5 = 133 = 0x85
	// All other pixels are 0
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	expectedCmd := "port 0 ; write1 c2 0 0 0 0 85 0 0 0 0\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestMatrix_SetPixel_InvalidCoordinates(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortB)

	testCases := []struct {
		x, y       int
		shouldFail bool
	}{
		{-1, 1, true}, // x too low
		{3, 1, true},  // x too high
		{1, -1, true}, // y too low
		{1, 3, true},  // y too high
		{0, 0, false}, // valid
		{2, 2, false}, // valid
	}

	for _, tc := range testCases {
		err := matrix.SetPixel(tc.x, tc.y, MatrixColor(5), 5)
		if tc.shouldFail && err == nil {
			t.Errorf("SetPixel(%d, %d) should have failed but didn't", tc.x, tc.y)
		}
		if !tc.shouldFail && err != nil {
			t.Errorf("SetPixel(%d, %d) failed: %v", tc.x, tc.y, err)
		}
	}
}

func TestMatrix_SetPixel_InvalidColorBrightness(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortC)

	testCases := []struct {
		color, brightness int
		shouldFail        bool
	}{
		{-1, 5, true},   // color too low
		{11, 5, true},   // color too high
		{5, -1, true},   // brightness too low
		{5, 11, true},   // brightness too high
		{0, 0, false},   // valid
		{10, 10, false}, // valid
		{5, 5, false},   // valid
	}

	for _, tc := range testCases {
		err := matrix.SetPixel(0, 0, MatrixColor(tc.color), tc.brightness)
		if tc.shouldFail && err == nil {
			t.Errorf("SetPixel with color=%d, brightness=%d should have failed but didn't", tc.color, tc.brightness)
		}
		if !tc.shouldFail && err != nil {
			t.Errorf("SetPixel with color=%d, brightness=%d failed: %v", tc.color, tc.brightness, err)
		}
	}
}

func TestMatrix_SetAll(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortD)

	err := matrix.SetAll(MatrixColor(7), 9)
	if err != nil {
		t.Fatalf("SetAll failed: %v", err)
	}

	// Verify EXACT command: "port 3 ; write1 c2 97 97 97 97 97 97 97 97 97\r"
	// SetAll(7, 9) sets all 9 pixels to {Color:7, Brightness:9}
	// Byte value: (9 << 4) | 7 = 144 + 7 = 151 = 0x97
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	expectedCmd := "port 3 ; write1 c2 97 97 97 97 97 97 97 97 97\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestMatrix_SetAll_InvalidValues(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortA)

	testCases := []struct {
		color, brightness int
		shouldFail        bool
	}{
		{-1, 5, true},
		{11, 5, true},
		{5, -1, true},
		{5, 11, true},
		{0, 0, false},
		{10, 10, false},
	}

	for _, tc := range testCases {
		err := matrix.SetAll(MatrixColor(tc.color), tc.brightness)
		if tc.shouldFail && err == nil {
			t.Errorf("SetAll(%d, %d) should have failed but didn't", tc.color, tc.brightness)
		}
		if !tc.shouldFail && err != nil {
			t.Errorf("SetAll(%d, %d) failed: %v", tc.color, tc.brightness, err)
		}
	}
}

func TestMatrix_SetRow(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortB)
	mockPort := brick.GetMockPort()

	// Test setting row 0
	err := matrix.SetRow(0, MatrixColor(6), 8)
	if err != nil {
		t.Fatalf("SetRow(0) failed: %v", err)
	}

	// Verify EXACT command: "port 1 ; write1 c2 86 86 86 0 0 0 0 0 0\r"
	// SetRow(0, 6, 8) sets pixels[0][0], [0][1], [0][2] to {Color:6, Brightness:8}
	// Byte value: (8 << 4) | 6 = 128 + 6 = 134 = 0x86
	commands := mockPort.GetWriteHistory()
	expectedCmd := "port 1 ; write1 c2 86 86 86 0 0 0 0 0 0\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestMatrix_SetRow_InvalidRow(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortC)

	err := matrix.SetRow(-1, MatrixColor(5), 5)
	if err == nil {
		t.Error("SetRow(-1) should have failed but didn't")
	}

	err = matrix.SetRow(3, MatrixColor(5), 5)
	if err == nil {
		t.Error("SetRow(3) should have failed but didn't")
	}
}

func TestMatrix_SetColumn(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortD)
	mockPort := brick.GetMockPort()

	// Test setting column 0
	err := matrix.SetColumn(0, MatrixColor(4), 7)
	if err != nil {
		t.Fatalf("SetColumn(0) failed: %v", err)
	}

	// Verify EXACT command: "port 3 ; write1 c2 74 0 0 74 0 0 74 0 0\r"
	// SetColumn(0, 4, 7) sets pixels[0][0], [1][0], [2][0] to {Color:4, Brightness:7}
	// Matrix layout: [0][0], [0][1], [0][2], [1][0], [1][1], [1][2], [2][0], [2][1], [2][2]
	// So positions 0, 3, 6 should be set
	// Byte value: (7 << 4) | 4 = 112 + 4 = 116 = 0x74
	commands := mockPort.GetWriteHistory()
	expectedCmd := "port 3 ; write1 c2 74 0 0 74 0 0 74 0 0\r"
	lastCmd := commands[len(commands)-1]
	if lastCmd != expectedCmd {
		t.Errorf("Expected exact command '%s', got: %s", expectedCmd, lastCmd)
	}
}

func TestMatrix_SetColumn_InvalidColumn(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortA)

	err := matrix.SetColumn(-1, MatrixColor(5), 5)
	if err == nil {
		t.Error("SetColumn(-1) should have failed but didn't")
	}

	err = matrix.SetColumn(3, MatrixColor(5), 5)
	if err == nil {
		t.Error("SetColumn(3) should have failed but didn't")
	}
}

func TestMatrix_Clear(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	matrix := brick.Matrix(PortB)

	// Set some pixels first
	if err := matrix.SetPixel(0, 0, MatrixColor(10), 10); err != nil {
		t.Fatalf("SetPixel(0, 0) failed: %v", err)
	}
	if err := matrix.SetPixel(1, 1, MatrixColor(10), 10); err != nil {
		t.Fatalf("SetPixel(1, 1) failed: %v", err)
	}

	// Clear the matrix
	err := matrix.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify all pixels are 0
	for x := range 3 {
		for y := range 3 {
			if matrix.pixels[x][y].Color != 0 || matrix.pixels[x][y].Brightness != 0 {
				t.Errorf("Pixel [%d][%d] not cleared: color=%d, brightness=%d",
					x, y, matrix.pixels[x][y].Color, matrix.pixels[x][y].Brightness)
			}
		}
	}
}

func TestMatrix_AllPorts(t *testing.T) {
	brick := TestBrick(t)
	defer CleanupTestBrick(brick)

	tests := []struct {
		port         Port
		expectedPort int // The numeric port value in the command
	}{
		{PortA, 0},
		{PortB, 1},
		{PortC, 2},
		{PortD, 3},
	}

	for _, tc := range tests {
		mockPort := brick.GetMockPort()
		mockPort.ClearWriteHistory()

		matrix := brick.Matrix(tc.port)

		// Set a pixel to trigger a write command
		err := matrix.SetPixel(0, 0, MatrixWhite, 5)
		if err != nil {
			t.Errorf("Port %s: SetPixel failed: %v", tc.port, err)
			continue
		}

		// Verify the command uses the correct port number
		history := mockPort.GetWriteHistory()
		if len(history) == 0 {
			t.Errorf("Port %s: no commands sent", tc.port)
			continue
		}

		// Check that the command contains the correct port number
		expectedPrefix := fmt.Sprintf("port %d ; ", tc.expectedPort)
		if !strings.HasPrefix(history[0], expectedPrefix) {
			t.Errorf("Port %s: expected command to start with '%s', got: %s", tc.port, expectedPrefix, history[0])
		}
	}
}
