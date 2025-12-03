package output

import (
	"fmt"
	"os"
	"strings"
	"time"
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

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d == 0 {
		return ""
	}

	// Round to tenths of a second for display
	seconds := d.Seconds()
	if seconds < 1 {
		return fmt.Sprintf("%.1fs", seconds)
	} else if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	} else {
		minutes := int(seconds / 60)
		remainingSeconds := seconds - float64(minutes*60)
		return fmt.Sprintf("%dm %.1fs", minutes, remainingSeconds)
	}
}

// formatPackageList formats a list of packages for display
func formatPackageList(packages []string, maxDisplay int) string {
	if len(packages) == 0 {
		return ""
	}

	if len(packages) <= maxDisplay {
		return strings.Join(packages, ", ")
	}

	// Show first few packages and indicate there are more
	displayed := strings.Join(packages[:maxDisplay], ", ")
	remaining := len(packages) - maxDisplay
	return fmt.Sprintf("%s, and %d more", displayed, remaining)
}

// PrintSummary prints the final summary of all operations
func PrintSummary(steps []UpdateStep) {
	Header("=== Update Summary ===")

	successCount := 0
	failCount := 0
	var totalDuration time.Duration

	// First pass: display basic status with timing and details
	for _, step := range steps {
		totalDuration += step.Duration

		if step.Success {
			// Format the status line with duration
			statusLine := step.Name
			if step.Duration > 0 {
				statusLine += fmt.Sprintf(" (%s)", formatDuration(step.Duration))
			}
			Success(statusLine)

			// Add detailed information on subsequent lines if available
			if step.PackageCount > 0 {
				if len(step.Packages) > 0 {
					// Show specific package names
					packageList := formatPackageList(step.Packages, 8)
					fmt.Printf("  %s %d package(s) updated: %s\n",
						Colorize(ColorBlue, "•"),
						step.PackageCount,
						packageList)
				} else {
					// Just show count
					fmt.Printf("  %s %d package(s) updated\n",
						Colorize(ColorBlue, "•"),
						step.PackageCount)
				}
			}

			if step.DiskSpaceFreed != "" {
				fmt.Printf("  %s Freed %s of disk space\n",
					Colorize(ColorBlue, "•"),
					step.DiskSpaceFreed)
			}

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

	// Print final summary with total time
	fmt.Println()
	if failCount == 0 {
		Success(fmt.Sprintf("All %d operation(s) completed successfully!", successCount))
	} else {
		Warning(fmt.Sprintf("%d succeeded, %d failed", successCount, failCount))
	}

	if totalDuration > 0 {
		fmt.Printf("Total time: %s\n", formatDuration(totalDuration))
	}
}

// UpdateStep represents a single operation and its result
type UpdateStep struct {
	Name           string
	Success        bool
	Error          error
	Packages       []string
	PackageCount   int
	DiskSpaceFreed string
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
}
