package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/e6a5/learning/experiment/ternimal-with-go/ansi"
)

func run(args []string) (string, error) {
	x, y, char, color, err := parseArgs(args)
	if err != nil {
		return "", err
	}

	if err := validateArgs(x, y); err != nil {
		return "", err
	}
	runes := []rune(char)
	if len(runes) != 1 {
		return "", fmt.Errorf("char must be exactly one character, got %d", len(runes))
	}

	colorCode, err := colorNameToCode(color)
	if err != nil {
		return "", err
	}
	result := ansi.PrintAtCoordinatesWithColor(x, y, runes[0], colorCode)
	return result, nil
}

func parseArgs(args []string) (int, int, string, string, error) {
	fs := flag.NewFlagSet("draw-at", flag.ContinueOnError)
	x := fs.Int("x", 0, "x coordinate")
	y := fs.Int("y", 0, "y coordinate")
	char := fs.String("char", "", "character to print")
	color := fs.String("color", "", "color to print")

	if err := fs.Parse(args); err != nil {
		return 0, 0, "", "", err
	}

	return *x, *y, *char, *color, nil
}

func validateArgs(x, y int) error {
	if x < 0 || y < 0 {
		return fmt.Errorf("x and y must be positive")
	}
	return nil
}

func colorNameToCode(colorName string) (int, error) {
	// using map to store the color name and code
	colorMap := map[string]int{
		"red":     31,
		"green":   32,
		"yellow":  33,
		"blue":    34,
		"magenta": 35,
		"cyan":    36,
		"white":   37,
	}

	if code, ok := colorMap[colorName]; ok {
		return code, nil
	}
	return 0, nil
}

func main() {
	result, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
