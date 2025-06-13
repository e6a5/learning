package ansi

import "testing"

func TestPrintAtCoordinates(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		text     rune
		expected string
	}{
		{1, 2, 'X', ESC + "[1;2HX"},
		{1, 2, 'Y', ESC + "[1;2HY"},
		{1, 3, 'Z', ESC + "[1;3HZ"},
		{2, 1, 'A', ESC + "[2;1HA"},
		{2, 2, 'B', ESC + "[2;2HB"},
		{2, 3, 'C', ESC + "[2;3HC"},
		{3, 1, 'D', ESC + "[3;1HD"},
		{3, 2, 'E', ESC + "[3;2HE"},
	}

	for _, test := range tests {
		result := PrintAtCoordinates(test.x, test.y, test.text)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestClearScreen(t *testing.T) {
	expected := ESC + "[2J"
	result := ClearScreen()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestClearLine(t *testing.T) {
	expected := ESC + "[2K"
	result := ClearLine()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestMoveCursor(t *testing.T) {
	expected := ESC + "[1;2H"
	result := MoveCursor(1, 2)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestColorize(t *testing.T) {
	expected := ESC + "[31mHello, World!" + ESC + "[0m"
	result := Colorize("Hello, World!", 31)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
