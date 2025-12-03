package updaters

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"system-update/internal/output"
	"system-update/internal/runner"
)

// Homebrew manages Homebrew package updates
type Homebrew struct {
	runner *runner.Runner
}

// NewHomebrew creates a new Homebrew updater
func NewHomebrew(r *runner.Runner) *Homebrew {
	return &Homebrew{runner: r}
}

// IsInstalled checks if Homebrew is installed on the system
func (h *Homebrew) IsInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

// parseBrewUpgrade extracts upgraded packages from brew upgrade output
// Example output lines:
// ==> Upgrading 3 outdated packages:
// git 2.42.0 -> 2.43.0
// node 20.5.0 -> 20.6.0
func parseBrewUpgrade(output string) ([]string, int) {
	if output == "" {
		return nil, 0
	}

	packages := []string{}
	lines := strings.Split(output, "\n")

	// Pattern to match upgrade lines like "git 2.42.0 -> 2.43.0"
	upgradePattern := regexp.MustCompile(`^(\S+)\s+[\d.]+\s+->\s+[\d.]+`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := upgradePattern.FindStringSubmatch(line); len(matches) > 1 {
			packages = append(packages, matches[1])
		}
	}

	// Also check for "Already up-to-date" message
	if strings.Contains(output, "Already up-to-date") ||
		strings.Contains(output, "already installed") {
		return nil, 0
	}

	return packages, len(packages)
}

// parseBrewCleanup extracts disk space freed from brew cleanup output
// Example output: "Pruned 5 symbolic links and removed 10 files from /usr/local/Cellar."
// Or: "This operation has freed approximately 2.3GB of disk space."
func parseBrewCleanup(output string) string {
	if output == "" {
		return ""
	}

	// Pattern to match disk space like "2.3GB" or "450MB"
	spacePattern := regexp.MustCompile(`freed approximately ([0-9.]+\s*[KMGT]B)`)
	if matches := spacePattern.FindStringSubmatch(output); len(matches) > 1 {
		return matches[1]
	}

	// Alternative pattern: "Removing: ... (150MB)"
	altPattern := regexp.MustCompile(`\(([0-9.]+\s*[KMGT]B)\)`)
	totalSize := 0.0

	for _, match := range altPattern.FindAllStringSubmatch(output, -1) {
		if len(match) > 1 {
			// Simple accumulation - this is approximate
			var size float64
			var u string
			fmt.Sscanf(match[1], "%f%s", &size, &u)
			if u == "GB" {
				size *= 1024
			}
			totalSize += size
		}
	}

	if totalSize > 1024 {
		return fmt.Sprintf("%.1fGB", totalSize/1024)
	} else if totalSize > 0 {
		return fmt.Sprintf("%.1fMB", totalSize)
	}

	return ""
}

// Update runs the complete Homebrew update sequence
func (h *Homebrew) Update() {
	if !h.IsInstalled() {
		output.Warning("Homebrew is not installed; skipping Homebrew updates")
		h.runner.RecordStep("Homebrew Update", fmt.Errorf("brew not found in PATH"), nil)
		return
	}

	// Step 1: Update Homebrew itself
	output.Progress("Updating Homebrew...")
	result := h.runner.ExecuteCommand("brew update", "brew", "update")
	h.runner.RecordStep("Homebrew Update", result.Error, &runner.StepDetails{
		Duration: result.Duration,
	})

	// Step 2: Upgrade packages (continue even if update failed)
	output.Progress("Upgrading all Homebrew packages...")
	result = h.runner.ExecuteCommand("brew upgrade", "brew", "upgrade")

	// Parse upgrade output to extract package information
	packages, count := parseBrewUpgrade(result.Output)
	h.runner.RecordStep("Homebrew Upgrade", result.Error, &runner.StepDetails{
		Packages:     packages,
		PackageCount: count,
		Duration:     result.Duration,
	})

	// Step 3: Cleanup old versions
	output.Progress("Cleaning up old Homebrew versions and caches...")
	result = h.runner.ExecuteCommand("brew cleanup", "brew", "cleanup")

	// Parse cleanup output to extract disk space freed
	diskSpaceFreed := parseBrewCleanup(result.Output)
	h.runner.RecordStep("Homebrew Cleanup", result.Error, &runner.StepDetails{
		DiskSpaceFreed: diskSpaceFreed,
		Duration:       result.Duration,
	})
}
