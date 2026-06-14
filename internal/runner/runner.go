package runner

import (
	"bytes"
	"fmt"
	"os/exec"
	"system-update/internal/output"
	"time"
)

// allowedCommands is the allowlist of command names that can be executed
var allowedCommands = map[string]struct{}{
	"brew": {},
	"npm":  {},
}

// lookupCommand validates the command against the allowlist and resolves its
// absolute path via PATH. This keeps the allowlist guard against arbitrary
// commands while letting the binary resolve to wherever it is actually
// installed (e.g. /opt/homebrew/bin on Apple Silicon, /usr/local/bin on Intel,
// or an nvm/fnm/volta-managed location for npm).
func lookupCommand(command string) (string, error) {
	if _, ok := allowedCommands[command]; !ok {
		return "", fmt.Errorf("command %q is not in the allowlist", command)
	}
	return exec.LookPath(command)
}

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

// CommandResult holds the results of a command execution
type CommandResult struct {
	Output   string
	Error    error
	Duration time.Duration
}

// ExecuteCommand runs a command and returns its output, duration, and any error.
// The command must be in the allowlist of permitted commands.
func (r *Runner) ExecuteCommand(name, command string, args ...string) CommandResult {
	startTime := time.Now()

	cmdPath, err := lookupCommand(command)
	if err != nil {
		return CommandResult{
			Output:   "",
			Error:    err,
			Duration: 0,
		}
	}

	if r.config.DryRun {
		output.Info(fmt.Sprintf("[DRY RUN] Would execute: %s %v", cmdPath, args))
		return CommandResult{
			Output:   "",
			Error:    nil,
			Duration: 0,
		}
	}

	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command -- cmdPath is exec.LookPath of an allowlist-validated command name ("brew"/"npm"); no untrusted input reaches here
	cmd := exec.Command(cmdPath, args...) // #nosec G204 -- cmdPath resolved from allowedCommands allowlist
	var stdout, stderr bytes.Buffer

	if r.config.Verbose {
		// In verbose mode, we want to show output in real-time
		// but also capture it for parsing
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		duration := time.Since(startTime)

		// Print captured output to stderr for verbose mode
		if stdout.Len() > 0 {
			output.VerboseWriter{}.Write(stdout.Bytes())
		}
		if stderr.Len() > 0 {
			output.VerboseWriter{}.Write(stderr.Bytes())
		}

		if err != nil {
			return CommandResult{
				Output:   stdout.String(),
				Error:    fmt.Errorf("%v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String()),
				Duration: duration,
			}
		}

		return CommandResult{
			Output:   stdout.String(),
			Error:    nil,
			Duration: duration,
		}
	}

	// Capture output for parsing and error reporting
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		// Return detailed error with command output
		return CommandResult{
			Output:   stdout.String(),
			Error:    fmt.Errorf("%v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String()),
			Duration: duration,
		}
	}

	return CommandResult{
		Output:   stdout.String(),
		Error:    nil,
		Duration: duration,
	}
}

// StepDetails contains optional detailed information about a step
type StepDetails struct {
	Packages       []string
	PackageCount   int
	DiskSpaceFreed string
	Duration       time.Duration
}

// RecordStep adds a step result to the runner's history
func (r *Runner) RecordStep(name string, err error, details *StepDetails) {
	step := output.UpdateStep{
		Name:    name,
		Success: err == nil,
		Error:   err,
	}

	// Add optional detailed information if provided
	if details != nil {
		step.Packages = details.Packages
		step.PackageCount = details.PackageCount
		step.DiskSpaceFreed = details.DiskSpaceFreed
		step.Duration = details.Duration
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
