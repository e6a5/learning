package drawing

import (
	"testing"

	"github.com/e6a5/learning/experiment/ternimal-with-go/ansi"
)

func TestDrawLine(t *testing.T) {
	tests := []struct {
		name           string
		x1, y1, x2, y2 int
		char           rune
		expected       string
	}{
		{
			name:     "horizontal line",
			x1:       1,
			y1:       1,
			x2:       5,
			y2:       1,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[1;2HX" + ansi.ESC + "[1;3HX" + ansi.ESC + "[1;4HX" + ansi.ESC + "[1;5HX",
		},
		{
			name:     "vertical line",
			x1:       1,
			y1:       1,
			x2:       1,
			y2:       5,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[2;1HX" + ansi.ESC + "[3;1HX" + ansi.ESC + "[4;1HX" + ansi.ESC + "[5;1HX",
		},
		{
			name:     "diagonal line",
			x1:       1,
			y1:       1,
			x2:       5,
			y2:       5,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[2;2HX" + ansi.ESC + "[3;3HX" + ansi.ESC + "[4;4HX" + ansi.ESC + "[5;5HX",
		},
		// TODO: add more tests for other diagonal lines like top-left to bottom-right, bottom-right to top-left, top-right to bottom-left, bottom-left to top-right
		{
			name:     "top-left to bottom-right diagonal line",
			x1:       1,
			y1:       5,
			x2:       3,
			y2:       3,
			char:     'X',
			expected: ansi.ESC + "[5;1HX" + ansi.ESC + "[4;2HX" + ansi.ESC + "[3;3HX",
		},
		{
			name:     "bottom-right to top-left diagonal line",
			x1:       3,
			y1:       3,
			x2:       1,
			y2:       5,
			char:     'X',
			expected: ansi.ESC + "[3;3HX" + ansi.ESC + "[4;2HX" + ansi.ESC + "[5;1HX",
		},
		{
			name:     "top-right to bottom-left diagonal line",
			x1:       3,
			y1:       1,
			x2:       1,
			y2:       3,
			char:     'X',
			expected: ansi.ESC + "[1;3HX" + ansi.ESC + "[2;2HX" + ansi.ESC + "[3;1HX",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := DrawLine(test.x1, test.y1, test.x2, test.y2, test.char)
			if result != test.expected {
				t.Errorf("DrawLine() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestDrawRect(t *testing.T) {
	tests := []struct {
		name          string
		x, y          int
		width, height int
		char          rune
		expected      string
	}{
		{
			name:     "rectangle",
			x:        1,
			y:        1,
			width:    3,
			height:   2,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[1;2HX" + ansi.ESC + "[1;3HX" + ansi.ESC + "[2;1HX" + ansi.ESC + "[2;2HX" + ansi.ESC + "[2;3HX",
		},
		{
			name:     "rectangle with width 1",
			x:        1,
			y:        1,
			width:    1,
			height:   3,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[2;1HX" + ansi.ESC + "[3;1HX",
		},
		{
			name:     "rectangle with height 1",
			x:        1,
			y:        1,
			width:    3,
			height:   1,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[1;2HX" + ansi.ESC + "[1;3HX",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := DrawRect(test.x, test.y, test.width, test.height, test.char)
			if result != test.expected {
				t.Errorf("DrawRect() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestDrawBox(t *testing.T) {
	tests := []struct {
		name          string
		x, y          int
		width, height int
		char          rune
		expected      string
	}{
		{
			name:     "box with width 1",
			x:        1,
			y:        1,
			width:    1,
			height:   3,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[2;1HX" + ansi.ESC + "[3;1HX",
		},
		{
			name:     "box with height 1",
			x:        1,
			y:        1,
			width:    3,
			height:   1,
			char:     'X',
			expected: ansi.ESC + "[1;1HX" + ansi.ESC + "[1;2HX" + ansi.ESC + "[1;3HX",
		},
		{
			name:     "box with width 1 and height 1",
			x:        1,
			y:        1,
			width:    1,
			height:   1,
			char:     'X',
			expected: ansi.ESC + "[1;1HX",
		},
		{
			name:     "box with width 4 and height 3",
			x:        2,
			y:        2,
			width:    4,
			height:   3,
			char:     'X',
			expected: ansi.ESC + "[2;2HX" + ansi.ESC + "[2;3HX" + ansi.ESC + "[2;4HX" + ansi.ESC + "[2;5HX" + ansi.ESC + "[3;2HX" + ansi.ESC + "[3;5HX" + ansi.ESC + "[4;2HX" + ansi.ESC + "[4;3HX" + ansi.ESC + "[4;4HX" + ansi.ESC + "[4;5HX",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := DrawBox(test.x, test.y, test.width, test.height, test.char)
			if result != test.expected {
				t.Errorf("DrawBox() = %q, want %q", result, test.expected)
			}
		})
	}
}
