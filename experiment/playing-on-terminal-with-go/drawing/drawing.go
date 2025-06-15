package drawing

import "github.com/e6a5/learning/experiment/ternimal-with-go/ansi"

func DrawLine(x1, y1, x2, y2 int, char rune) string {
	if y1 == y2 {
		return DrawHorizontalLine(x1, x2, y1, char)
	}
	if x1 == x2 {
		return DrawVerticalLine(x1, y1, y2, char)
	}
	if x1 < x2 && y1 < y2 {
		return DrawDiagonalLine(x1, y1, x2, y2, char)
	}
	return ""
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
	for x, y := x1, y1; x <= x2 && y <= y2; x, y = x+1, y+1 {
		result += ansi.PrintAtCoordinates(x, y, char)
	}
	return result
}
