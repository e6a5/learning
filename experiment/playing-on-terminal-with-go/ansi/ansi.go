package ansi

import "fmt"

const (
	ESC = "\033"
)

func PrintAtCoordinates(x, y int, text rune) string {
	return fmt.Sprintf("%s[%d;%dH%c", ESC, x, y, text)
}

func ClearScreen() string {
	return fmt.Sprintf("%s[2J", ESC)
}

func ClearLine() string {
	return fmt.Sprintf("%s[2K", ESC)
}

func MoveCursor(x, y int) string {
	return fmt.Sprintf("%s[%d;%dH", ESC, x, y)
}

func Colorize(text string, color int) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", ESC, color, text, ESC)
}
