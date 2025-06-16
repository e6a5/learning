package drawing

import "github.com/e6a5/learning/experiment/ternimal-with-go/ansi"

func DrawLine(x1, y1, x2, y2 int, char rune) string {
	if y1 == y2 {
		return DrawHorizontalLine(x1, x2, y1, char)
	}
	if x1 == x2 {
		return DrawVerticalLine(x1, y1, y2, char)
	}

	return DrawDiagonalLine(x1, y1, x2, y2, char)
}

func DrawHorizontalLine(x1, x2, y int, char rune) string {
	result := ""
	for x := x1; x <= x2; x++ {
		result += ansi.PrintAtCoordinates(x, y, char)
	}
	return result
}

func DrawVerticalLine(x, y1, y2 int, char rune) string {
	result := ""
	for y := y1; y <= y2; y++ {
		result += ansi.PrintAtCoordinates(x, y, char)
	}
	return result
}

func DrawDiagonalLine(x1, y1, x2, y2 int, char rune) string {
	result := ""

	// Calculate step directions
	xStep := 1
	if x1 > x2 {
		xStep = -1
	}

	yStep := 1
	if y1 > y2 {
		yStep = -1
	}

	// Draw line with calculated steps
	for x, y := x1, y1; x != x2+xStep && y != y2+yStep; x, y = x+xStep, y+yStep {
		result += ansi.PrintAtCoordinates(x, y, char)
	}

	return result
}
