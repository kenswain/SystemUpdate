package main

import (
	"flag"
	"fmt"
	"os"
	"system-update/internal/output"
	"system-update/internal/runner"
	"system-update/internal/updaters"
)

const version = "1.0.0"

func main() {
	// Define command-line flags
	brewOnly := flag.Bool("brew-only", false, "Only update Homebrew packages")
	npmOnly := flag.Bool("npm-only", false, "Only update npm packages")
	dryRun := flag.Bool("dry-run", false, "Show what would be executed without running commands")
	verbose := flag.Bool("verbose", false, "Show detailed command output")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "system-update v%s - Update system packages (Homebrew and npm)\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: system-update [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  system-update                    # Update both Homebrew and npm\n")
		fmt.Fprintf(os.Stderr, "  system-update --brew-only        # Only update Homebrew\n")
		fmt.Fprintf(os.Stderr, "  system-update --npm-only         # Only update npm\n")
		fmt.Fprintf(os.Stderr, "  system-update --dry-run          # Preview what would be executed\n")
		fmt.Fprintf(os.Stderr, "  system-update --verbose          # Show detailed command output\n")
	}

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("system-update version %s\n", version)
		os.Exit(0)
	}

	// Validate flag combinations
	if *brewOnly && *npmOnly {
		fmt.Fprintf(os.Stderr, "Error: --brew-only and --npm-only cannot be used together\n")
		os.Exit(1)
	}

	// Display mode indicators
	if *dryRun {
		output.Info("=== DRY RUN MODE - No commands will be executed ===")
	}
	if *verbose {
		output.Info("=== VERBOSE MODE - Showing detailed output ===")
	}

	// Create runner with configuration
	config := runner.Config{
		DryRun:  *dryRun,
		Verbose: *verbose,
	}
	r := runner.New(config)

	// Determine which updaters to run
	runBrew := !*npmOnly
	runNpm := !*brewOnly

	// Execute updates
	output.Header("Starting system update...")

	if runBrew {
		brew := updaters.NewHomebrew(r)
		brew.Update()
	}

	if runNpm {
		npm := updaters.NewNPM(r)
		npm.Update()
	}

	// Print summary
	output.PrintSummary(r.GetSteps())

	// Exit with appropriate code
	if r.HasFailures() {
		os.Exit(1)
	}

	output.Success("System update complete!")
	os.Exit(0)
}
