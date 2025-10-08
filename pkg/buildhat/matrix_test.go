package buildhat

import (
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

	// Verify command was sent
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	lastCmd := commands[len(commands)-1]
	if !strings.Contains(lastCmd, "port 0") || !strings.Contains(lastCmd, "write1") {
		t.Errorf("Expected write1 command, got: %s", lastCmd)
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

	// Verify command was sent
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) == 0 {
		t.Fatal("No commands were sent")
	}

	lastCmd := commands[len(commands)-1]
	if !strings.Contains(lastCmd, "port 3") || !strings.Contains(lastCmd, "write1") {
		t.Errorf("Expected write1 command, got: %s", lastCmd)
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

	for row := 0; row < 3; row++ {
		err := matrix.SetRow(row, MatrixColor(6), 8)
		if err != nil {
			t.Fatalf("SetRow(%d) failed: %v", row, err)
		}
	}

	// Verify commands were sent
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) < 3 {
		t.Fatalf("Expected at least 3 commands, got %d", len(commands))
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

	for col := 0; col < 3; col++ {
		err := matrix.SetColumn(col, MatrixColor(4), 7)
		if err != nil {
			t.Fatalf("SetColumn(%d) failed: %v", col, err)
		}
	}

	// Verify commands were sent
	mockPort := brick.GetMockPort()
	commands := mockPort.GetWriteHistory()
	if len(commands) < 3 {
		t.Fatalf("Expected at least 3 commands, got %d", len(commands))
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

	ports := []struct {
		port     Port
		expected int
	}{
		{PortA, 0},
		{PortB, 1},
		{PortC, 2},
		{PortD, 3},
	}

	for _, tc := range ports {
		matrix := brick.Matrix(tc.port)
		if matrix.port != tc.expected {
			t.Errorf("Port %s: expected port number %d, got %d", tc.port, tc.expected, matrix.port)
		}
	}
}
