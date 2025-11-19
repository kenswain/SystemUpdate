package updaters

import (
	"fmt"
	"os/exec"
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

// Update runs the complete npm update sequence
func (n *NPM) Update() {
	if !n.IsInstalled() {
		output.Warning("npm is not installed; skipping npm package updates")
		n.runner.RecordStep("npm Update", fmt.Errorf("npm not found in PATH"))
		return
	}

	// Step 1: Update npm itself
	output.Progress("Updating npm CLI...")
	err := n.runner.ExecuteCommand("npm install", "npm", "install", "-g", "npm")
	n.runner.RecordStep("npm CLI Update", err)

	// Step 2: Update global packages
	output.Progress("Updating all globally installed npm packages...")
	err = n.runner.ExecuteCommand("npm update", "npm", "update", "-g")
	n.runner.RecordStep("npm Global Update", err)
}
