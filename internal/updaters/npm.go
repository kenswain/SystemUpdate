package updaters

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"system-update/internal/output"
	"system-update/internal/runner"
)

// NPM manages npm package updates
type NPM struct {
	runner *runner.Runner
}

// NewNPM creates a new NPM updater
func NewNPM(r *runner.Runner) *NPM {
	return &NPM{runner: r}
}

// IsInstalled checks if npm is installed on the system
func (n *NPM) IsInstalled() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

// parseNpmUpdate extracts updated packages from npm update output
// Example output lines:
// added 5 packages, removed 2 packages, changed 8 packages in 5s
// Or individual package lines like:
// + typescript@5.2.0
// + eslint@8.50.0
func parseNpmUpdate(output string) ([]string, int) {
	if output == "" {
		return nil, 0
	}

	packages := []string{}

	// Pattern to match npm package update lines like "+ typescript@5.2.0"
	packagePattern := regexp.MustCompile(`^\+\s+(\S+)@`)

	// Pattern to match summary line like "changed 8 packages"
	changedPattern := regexp.MustCompile(`changed\s+(\d+)\s+packages?`)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Try to match individual package updates
		if matches := packagePattern.FindStringSubmatch(line); len(matches) > 1 {
			packages = append(packages, matches[1])
		}
	}

	// If we found specific packages, return them
	if len(packages) > 0 {
		return packages, len(packages)
	}

	// Otherwise, try to extract count from summary line
	if matches := changedPattern.FindStringSubmatch(output); len(matches) > 1 {
		var count int
		fmt.Sscanf(matches[1], "%d", &count)
		if count > 0 {
			return nil, count // Return count but no specific package names
		}
	}

	// Check for "up to date" message
	if strings.Contains(output, "up to date") ||
		strings.Contains(output, "already at latest") {
		return nil, 0
	}

	return nil, 0
}

// Update runs the complete npm update sequence
func (n *NPM) Update() {
	if !n.IsInstalled() {
		output.Warning("npm is not installed; skipping npm package updates")
		n.runner.RecordStep("npm Update", fmt.Errorf("npm not found in PATH"), nil)
		return
	}

	// Step 1: Update npm itself
	output.Progress("Updating npm CLI...")
	result := n.runner.ExecuteCommand("npm install", "npm", "install", "-g", "npm")

	// Parse npm install output to see if npm was actually updated
	packages, count := parseNpmUpdate(result.Output)
	n.runner.RecordStep("npm CLI Update", result.Error, &runner.StepDetails{
		Packages:     packages,
		PackageCount: count,
		Duration:     result.Duration,
	})

	// Step 2: Update global packages
	output.Progress("Updating all globally installed npm packages...")
	result = n.runner.ExecuteCommand("npm update", "npm", "update", "-g")

	// Parse npm update output to extract package information
	packages, count = parseNpmUpdate(result.Output)
	n.runner.RecordStep("npm Global Update", result.Error, &runner.StepDetails{
		Packages:     packages,
		PackageCount: count,
		Duration:     result.Duration,
	})
}
