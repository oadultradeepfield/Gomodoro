package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCIIDigits contains Claude Code logo style digits using $ characters (7 lines tall)
var ASCIIDigits = map[rune][]string{
	'0': {
		" $$$$$ ",
		"$$   $$",
		"$$   $$",
		"$$   $$",
		"$$   $$",
		"$$   $$",
		" $$$$$ ",
	},
	'1': {
		"   $$  ",
		"  $$$  ",
		"   $$  ",
		"   $$  ",
		"   $$  ",
		"   $$  ",
		" $$$$$ ",
	},
	'2': {
		" $$$$$ ",
		"$$   $$",
		"     $$",
		"  $$$$ ",
		" $$    ",
		"$$     ",
		"$$$$$$$",
	},
	'3': {
		" $$$$$ ",
		"$$   $$",
		"     $$",
		"  $$$$ ",
		"     $$",
		"$$   $$",
		" $$$$$ ",
	},
	'4': {
		"$$   $$",
		"$$   $$",
		"$$   $$",
		"$$$$$$$",
		"     $$",
		"     $$",
		"     $$",
	},
	'5': {
		"$$$$$$$",
		"$$     ",
		"$$     ",
		"$$$$$$ ",
		"     $$",
		"$$   $$",
		" $$$$$ ",
	},
	'6': {
		" $$$$$ ",
		"$$     ",
		"$$     ",
		"$$$$$$ ",
		"$$   $$",
		"$$   $$",
		" $$$$$ ",
	},
	'7': {
		"$$$$$$$",
		"     $$",
		"    $$ ",
		"   $$  ",
		"  $$   ",
		"  $$   ",
		"  $$   ",
	},
	'8': {
		" $$$$$ ",
		"$$   $$",
		"$$   $$",
		" $$$$$ ",
		"$$   $$",
		"$$   $$",
		" $$$$$ ",
	},
	'9': {
		" $$$$$ ",
		"$$   $$",
		"$$   $$",
		" $$$$$$",
		"     $$",
		"     $$",
		" $$$$$ ",
	},
	':': {
		"   ",
		"   ",
		" $ ",
		"   ",
		" $ ",
		"   ",
		"   ",
	},
}

const DigitHeight = 7

// RenderASCIITime converts a time string (e.g., "25:00") to ASCII art
func RenderASCIITime(timeStr string) string {
	lines := make([]string, DigitHeight)
	first := true
	for _, char := range timeStr {
		digit, exists := ASCIIDigits[char]
		if !exists {
			continue
		}
		for i := 0; i < DigitHeight; i++ {
			if !first {
				lines[i] += " " // Add one column spacing between digits
			}
			lines[i] += digit[i]
		}
		first = false
	}
	return strings.Join(lines, "\n")
}

// RenderStyledASCIITime renders ASCII time with a lipgloss style applied
func RenderStyledASCIITime(timeStr string, style lipgloss.Style) string {
	return style.Render(RenderASCIITime(timeStr))
}
