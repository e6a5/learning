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
