package main

import (
	"testing"

	"github.com/e6a5/learning/experiment/ternimal-with-go/ansi"
)

func TestRun(t *testing.T) {
	//Valid case: ["--x=5", "--y=10", "--char=X"] → what should expected be?
	//Error case: ["--x=-1", "--y=10", "--char=X"] → wantErr should be true
	tests := []struct {
		name     string
		args     []string
		expected string
		wantErr  bool
	}{
		{
			name:     "print at coordinates",
			args:     []string{"--x=5", "--y=10", "--char=X"},
			expected: ansi.ESC + "[5;10HX",
			wantErr:  false,
		},
		{
			name:     "print at coordinates with color",
			args:     []string{"--x=5", "--y=10", "--char=X", "--color=red"},
			expected: ansi.ESC + "[5;10H" + ansi.ESC + "[31mX" + ansi.ESC + "[0m",
			wantErr:  false,
		},
		{
			name:     "error case",
			args:     []string{"--x=-1", "--y=10", "--char=X"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := run(test.args)
			if (err != nil) != test.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, test.wantErr)
			}
			if result != test.expected {
				t.Errorf("run() result = %v, expected %v", result, test.expected)
			}
		})
	}
}
