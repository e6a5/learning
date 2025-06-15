package ansi

import "fmt"

const (
	ESC = "\033"
)

func PrintAtCoordinates(x, y int, text rune) string {
	return fmt.Sprintf("%s[%d;%dH%c", ESC, y, x, text)
}

func ClearScreen() string {
	return fmt.Sprintf("%s[2J", ESC)
}

func ClearLine() string {
	return fmt.Sprintf("%s[2K", ESC)
}

func MoveCursor(x, y int) string {
	return fmt.Sprintf("%s[%d;%dH", ESC, y, x)
}

func Colorize(text string, color int) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", ESC, color, text, ESC)
}

// if color code is 0, it will not be colored
func PrintAtCoordinatesWithColor(x, y int, char rune, colorCode int) string {
	positioned := MoveCursor(x, y)
	colored := ""
	if colorCode != 0 {
		colored = Colorize(string(char), colorCode)
	} else {
		colored = string(char)
	}
	return positioned + colored
}
