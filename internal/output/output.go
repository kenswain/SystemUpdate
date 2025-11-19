package output

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// Symbols for status indicators
const (
	SymbolSuccess = "✓"
	SymbolFail    = "✗"
	SymbolArrow   = "→"
	SymbolWarning = "⚠"
)

// IsColorEnabled checks if color output should be used
func IsColorEnabled() bool {
	// Disable colors if NO_COLOR env var is set or not a terminal
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return true
}

// Colorize wraps text with ANSI color codes if colors are enabled
func Colorize(color, text string) string {
	if !IsColorEnabled() {
		return text
	}
	return color + text + ColorReset
}

// Success prints a success message with green checkmark
func Success(message string) {
	fmt.Printf("%s %s\n", Colorize(ColorGreen, SymbolSuccess), message)
}

// Fail prints a failure message with red X
func Fail(message string) {
	fmt.Printf("%s %s\n", Colorize(ColorRed, SymbolFail), message)
}

// Warning prints a warning message with yellow symbol
func Warning(message string) {
	fmt.Printf("%s %s\n", Colorize(ColorYellow, SymbolWarning), message)
}

// Progress prints an in-progress message with blue arrow
func Progress(message string) {
	fmt.Printf("%s %s\n", Colorize(ColorBlue, SymbolArrow), message)
}

// Info prints an informational message in cyan
func Info(message string) {
	fmt.Printf("%s\n", Colorize(ColorCyan, message))
}

// Header prints a section header
func Header(message string) {
	fmt.Printf("\n%s\n", Colorize(ColorCyan, message))
}

// PrintSummary prints the final summary of all operations
func PrintSummary(steps []UpdateStep) {
	Header("=== Update Summary ===")

	successCount := 0
	failCount := 0

	for _, step := range steps {
		if step.Success {
			Success(step.Name)
			successCount++
		} else {
			if step.Error != nil {
				Fail(fmt.Sprintf("%s: %v", step.Name, step.Error))
			} else {
				Fail(step.Name + " (skipped)")
			}
			failCount++
		}
	}

	fmt.Println()
	if failCount == 0 {
		Success(fmt.Sprintf("All %d operation(s) completed successfully!", successCount))
	} else {
		Warning(fmt.Sprintf("%d succeeded, %d failed", successCount, failCount))
	}
}

// UpdateStep represents a single operation and its result
type UpdateStep struct {
	Name    string
	Success bool
	Error   error
}
