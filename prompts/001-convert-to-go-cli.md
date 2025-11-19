<objective>
Convert the bash system-update script into a professional, full-featured Go CLI application. The tool should update Homebrew packages and npm global packages with proper error handling, user-friendly output, and flexible configuration options via command-line flags.

This replaces a simple sequential bash script with a robust CLI tool that collects errors during execution and provides a comprehensive summary at the end, making it easier to understand what succeeded and what failed.
</objective>

<context>
The original bash script (`system-update`) performs three main operations in sequence:
1. Update Homebrew and upgrade all packages
2. Clean up old Homebrew versions
3. Update npm globally (if installed) and upgrade global npm packages

The new Go application should preserve this functionality while adding:
- Selective execution modes (brew-only, npm-only flags)
- Dry-run capability to preview changes without executing
- Better error handling and reporting
- Clear progress indicators and status messages

Build targets: macOS (darwin) with bash script semantics
</context>

<requirements>
Functional Requirements:
1. Implement core update operations matching bash script behavior:
   - `brew update` + `brew upgrade`
   - `brew cleanup`
   - `npm install -g npm` + `npm update -g` (if npm exists)

2. CLI Flags:
   - `--brew-only` - Skip npm updates
   - `--npm-only` - Skip brew updates
   - `--dry-run` - Show what would be executed without running commands
   - `--verbose` - Show detailed command output

3. Error Handling:
   - Continue with subsequent steps even if one command fails
   - Collect all errors that occur during execution
   - Display comprehensive summary at end showing:
     - What succeeded ✓
     - What failed ✗ (with error messages)
     - Overall status

4. User Experience:
   - Clear status messages before each operation begins
   - Progress indicators showing which step is running
   - Color-coded output (green for success, red for errors, yellow for warnings)
   - Helpful error messages if a tool isn't installed

Code Organization:
1. Project structure must use modular packages:
   - `cmd/system-update/main.go` - CLI entry point and flag parsing
   - `internal/updaters/homebrew.go` - Homebrew update logic
   - `internal/updaters/npm.go` - npm update logic
   - `internal/runner/runner.go` - Orchestration and error collection
   - `internal/output/output.go` - Formatted output and colors

2. Use Go best practices:
   - Proper error types and error handling patterns
   - Interfaces where appropriate (e.g., Updater interface for different package managers)
   - Clean separation of concerns

3. Include project files:
   - `go.mod` - Module definition
   - `Makefile` - Build target: `make build` creates `./system-update` binary
   - `README.md` - Usage documentation with examples

</requirements>

<implementation>
Architecture Pattern:
- Create an UpdateStep type that tracks operation name, success status, and any error
- Use a Runner type that collects UpdateStep results
- Return structured results (not just exit codes) to allow summary reporting

Error Handling Approach:
- Don't use panic() or fatal errors - catch command execution errors gracefully
- Store errors in a slice within the runner for later display
- Allow dependent operations (e.g., only run npm if it's installed) without failing the whole run

Command Execution:
- Use os/exec package to run shell commands
- In dry-run mode, print what would be executed instead of actually running it
- In verbose mode, stream command output to stderr in real-time
- Capture command output to show errors when commands fail

Output and Formatting:
- Use ANSI color codes or a simple library for colored output
  - ✓ (green) for successful operations
  - ✗ (red) for failed operations
  - → (blue) for current operations in progress
- Group all status messages together clearly
- Final summary should show overall results at a glance

What to Avoid:
- Don't use `os.Exit()` early - collect all errors and report them
- Don't suppress command output completely; let users see what's happening
- Don't create complex nested dependencies - keep the package structure flat
- Avoid global state; pass configuration through function parameters

Dependencies:
- Prefer stdlib where possible (os/exec, fmt, flag)
- Only add external dependencies if truly necessary (e.g., for color output, consider using conditional stdlib approach)
</implementation>

<output>
Create the following files:

1. `./cmd/system-update/main.go` - Entry point with CLI flag parsing and orchestration
2. `./internal/updaters/homebrew.go` - Homebrew update functions
3. `./internal/updaters/npm.go` - npm update functions (with existence check)
4. `./internal/runner/runner.go` - UpdateStep type and Runner orchestration
5. `./internal/output/output.go` - Formatted output with colors and status symbols
6. `./go.mod` - Go module file (module name: `system-update`)
7. `./Makefile` - Build target with `make build` creating `./system-update`
8. `./README.md` - Usage documentation with flag examples and sample output

All files should be production-ready with proper error handling, comments, and clean code structure.
</output>

<verification>
Before declaring the implementation complete, verify:

1. ✓ Binary builds successfully: `make build` creates `./system-update` executable
2. ✓ Help flag works: `./system-update --help` shows all available flags
3. ✓ Dry-run mode works: `./system-update --dry-run` shows what would execute without running
4. ✓ Selective updates work: `./system-update --brew-only` and `./system-update --npm-only` only run respective updates
5. ✓ Error collection works: If one step fails, others still execute and all errors appear in summary
6. ✓ npm existence check: If npm isn't installed, shows message and continues without failing
7. ✓ Summary displays properly: At end of run, shows what succeeded and what failed clearly

Test the implementation with at least one dry-run to verify output formatting.
</verification>

<success_criteria>
The Go application is complete and ready when:
- ✓ Modular package structure matches specification
- ✓ All CLI flags (--brew-only, --npm-only, --dry-run, --verbose) work correctly
- ✓ Error handling collects all errors and reports them in summary (not exiting early)
- ✓ Output is formatted clearly with status symbols and colors
- ✓ Makefile builds the binary to `./system-update`
- ✓ README documents all flags and usage examples
- ✓ Code follows Go best practices and is well-commented
</success_criteria>
