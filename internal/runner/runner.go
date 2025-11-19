package runner

import (
	"bytes"
	"fmt"
	"os/exec"
	"system-update/internal/output"
)

// Config holds runtime configuration for the update runner
type Config struct {
	DryRun  bool
	Verbose bool
}

// Runner orchestrates the update process and collects results
type Runner struct {
	config Config
	steps  []output.UpdateStep
}

// New creates a new Runner with the given configuration
func New(config Config) *Runner {
	return &Runner{
		config: config,
		steps:  make([]output.UpdateStep, 0),
	}
}

// ExecuteCommand runs a command and returns its output and any error
func (r *Runner) ExecuteCommand(name, command string, args ...string) error {
	if r.config.DryRun {
		output.Info(fmt.Sprintf("[DRY RUN] Would execute: %s %v", command, args))
		return nil
	}

	cmd := exec.Command(command, args...)

	if r.config.Verbose {
		// Stream output in real-time
		cmd.Stdout = output.VerboseWriter{}
		cmd.Stderr = output.VerboseWriter{}
		return cmd.Run()
	}

	// Capture output for error reporting
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Return detailed error with command output
		return fmt.Errorf("%v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	return nil
}

// RecordStep adds a step result to the runner's history
func (r *Runner) RecordStep(name string, err error) {
	step := output.UpdateStep{
		Name:    name,
		Success: err == nil,
		Error:   err,
	}
	r.steps = append(r.steps, step)

	// Provide immediate feedback
	if err == nil {
		output.Success(name)
	} else {
		output.Fail(fmt.Sprintf("%s: %v", name, err))
	}
}

// GetSteps returns all recorded steps
func (r *Runner) GetSteps() []output.UpdateStep {
	return r.steps
}

// HasFailures returns true if any step failed
func (r *Runner) HasFailures() bool {
	for _, step := range r.steps {
		if !step.Success {
			return true
		}
	}
	return false
}
