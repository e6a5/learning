package ansi

import "testing"

func TestPrintAtCoordinates(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		text     rune
		expected string
	}{
		{1, 2, 'X', ESC + "[2;1HX"},
		{1, 2, 'Y', ESC + "[2;1HY"},
		{1, 3, 'Z', ESC + "[3;1HZ"},
		{2, 1, 'A', ESC + "[1;2HA"},
		{2, 2, 'B', ESC + "[2;2HB"},
		{2, 3, 'C', ESC + "[3;2HC"},
		{3, 1, 'D', ESC + "[1;3HD"},
		{3, 2, 'E', ESC + "[2;3HE"},
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
	expected := ESC + "[2;1H"
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

func TestPrintAtCoordinatesWithColor(t *testing.T) {
	expected := ESC + "[10;5H" + ESC + "[31mX" + ESC + "[0m"
	result := PrintAtCoordinatesWithColor(5, 10, 'X', 31)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
