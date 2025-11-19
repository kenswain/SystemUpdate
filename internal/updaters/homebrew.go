package updaters

import (
	"fmt"
	"os/exec"
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

// Update runs the complete Homebrew update sequence
func (h *Homebrew) Update() {
	if !h.IsInstalled() {
		output.Warning("Homebrew is not installed; skipping Homebrew updates")
		h.runner.RecordStep("Homebrew Update", fmt.Errorf("brew not found in PATH"))
		return
	}

	// Step 1: Update Homebrew itself
	output.Progress("Updating Homebrew...")
	err := h.runner.ExecuteCommand("brew update", "brew", "update")
	h.runner.RecordStep("Homebrew Update", err)

	// Step 2: Upgrade packages (continue even if update failed)
	output.Progress("Upgrading all Homebrew packages...")
	err = h.runner.ExecuteCommand("brew upgrade", "brew", "upgrade")
	h.runner.RecordStep("Homebrew Upgrade", err)

	// Step 3: Cleanup old versions
	output.Progress("Cleaning up old Homebrew versions and caches...")
	err = h.runner.ExecuteCommand("brew cleanup", "brew", "cleanup")
	h.runner.RecordStep("Homebrew Cleanup", err)
}
