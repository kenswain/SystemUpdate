# system-update

A professional CLI tool and native macOS menu bar app for updating Homebrew and npm packages with comprehensive error handling and user-friendly output.

## Features

- **Homebrew Updates**: Updates Homebrew itself, upgrades all packages, and cleans up old versions
- **npm Updates**: Updates npm CLI and all globally installed packages
- **Robust Error Handling**: Continues execution even if individual steps fail, with comprehensive error reporting
- **Flexible Configuration**: Command-line flags for selective updates and execution modes
- **User-Friendly Output**: Color-coded status messages and clear progress indicators
- **Dry-Run Mode**: Preview operations without executing them
- **Verbose Mode**: See detailed command output for debugging
- **macOS Menu Bar App**: Native SwiftUI app with real-time progress display

## Installation

### CLI - Build from Source

```bash
# Clone the repository
git clone https://github.com/kenswain/SystemUpdate.git
cd SystemUpdate

# Build the binary
make build

# Optional: Install to /usr/local/bin for system-wide access
make install
```

### macOS Menu Bar App

1. Open `SystemUpdateGUI/SystemUpdateGUI.xcodeproj` in Xcode
2. Build and run (⌘R) or archive for distribution (Product → Archive)
3. The app appears in your menu bar

**Requirements**: Xcode 15+, macOS 13.0+ (Ventura or later)

## Usage

### Basic Usage

```bash
# Update both Homebrew and npm
./system-update

# Get help and see all available options
./system-update --help

# Show version
./system-update --version
```

### Command-Line Flags

| Flag | Description |
|------|-------------|
| `--brew-only` | Only update Homebrew packages (skip npm) |
| `--npm-only` | Only update npm packages (skip Homebrew) |
| `--dry-run` | Show what would be executed without running commands |
| `--verbose` | Show detailed command output in real-time |
| `--version` | Show version information |

### Examples

#### Update Everything
```bash
./system-update
```

#### Update Only Homebrew
```bash
./system-update --brew-only
```

#### Update Only npm
```bash
./system-update --npm-only
```

#### Preview What Would Happen (Dry Run)
```bash
./system-update --dry-run
```
Output:
```
=== DRY RUN MODE - No commands will be executed ===

Starting system update...
→ Updating Homebrew...
[DRY RUN] Would execute: brew [update]
✓ Homebrew Update
→ Upgrading all Homebrew packages...
[DRY RUN] Would execute: brew [upgrade]
✓ Homebrew Upgrade
...
```

#### Show Detailed Command Output
```bash
./system-update --verbose
```

#### Combine Flags
```bash
# Dry run with verbose output for Homebrew only
./system-update --brew-only --dry-run --verbose
```

## Output Format

The tool provides clear, color-coded output:

- **✓** (green) - Successful operations
- **✗** (red) - Failed operations
- **→** (blue) - Operations in progress
- **⚠** (yellow) - Warnings

### Example Output

```
Starting system update...
→ Updating Homebrew...
✓ Homebrew Update
→ Upgrading all Homebrew packages...
✓ Homebrew Upgrade
→ Cleaning up old Homebrew versions and caches...
✓ Homebrew Cleanup
→ Updating npm CLI...
✓ npm CLI Update
→ Updating all globally installed npm packages...
✓ npm Global Update

=== Update Summary ===
✓ Homebrew Update
✓ Homebrew Upgrade
✓ Homebrew Cleanup
✓ npm CLI Update
✓ npm Global Update

✓ All 5 operation(s) completed successfully!
✓ System update complete!
```

### Handling Failures

If any operation fails, the tool continues with remaining operations and provides a detailed summary:

```
=== Update Summary ===
✓ Homebrew Update
✗ Homebrew Upgrade: exit status 1
Stdout: ...
Stderr: Error: Permission denied
✓ Homebrew Cleanup
✓ npm CLI Update
✓ npm Global Update

⚠ 4 succeeded, 1 failed
```

Exit codes:
- `0` - All operations succeeded
- `1` - One or more operations failed

## Project Structure

```
system-update/
├── cmd/
│   └── system-update/
│       └── main.go              # CLI entry point and flag parsing
├── internal/
│   ├── output/
│   │   ├── output.go            # Formatted output and colors
│   │   └── verbose.go           # Verbose mode output handling
│   ├── runner/
│   │   └── runner.go            # Orchestration and error collection
│   └── updaters/
│       ├── homebrew.go          # Homebrew update logic
│       └── npm.go               # npm update logic
├── SystemUpdateGUI/             # macOS Menu Bar App
│   ├── SystemUpdateGUI.xcodeproj
│   └── SystemUpdateGUI/
│       ├── SystemUpdateGUIApp.swift   # App entry point
│       ├── MenuBarView.swift          # Main menu bar UI
│       ├── UpdateManager.swift        # Process management
│       ├── OutputParser.swift         # CLI output parsing
│       └── SettingsView.swift         # Settings panel
├── go.mod                       # Go module definition
├── Makefile                     # Build automation
└── README.md                    # This file
```

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

### Cleaning Build Artifacts

```bash
make clean
```

## Requirements

### CLI
- Go 1.21 or later
- Homebrew (optional - will skip if not installed)
- npm (optional - will skip if not installed)

### macOS App
- macOS 13.0 (Ventura) or later
- Xcode 15+ (for building from source)

## Error Handling

The tool follows a "collect and report" error handling strategy:

1. Each operation is attempted regardless of previous failures
2. Errors are collected during execution
3. A comprehensive summary is displayed at the end
4. Exit code reflects overall success/failure

This approach ensures you get as many updates as possible even if some operations fail.

## Environment Variables

- `NO_COLOR` - Set to any value to disable color output

## macOS Menu Bar App

The SwiftUI menu bar app provides a native macOS experience:

- **Menu Bar Icon**: Lives in your menu bar with status indicators
- **Real-Time Output**: Watch updates as they happen with color-coded status
- **Settings Panel**: Configure Homebrew-only, npm-only, or verbose mode
- **Launch at Login**: Optionally start the app when you log in

The app bundles the Go CLI binary and executes it as a subprocess, streaming output in real-time.

## License

This tool is provided as-is for personal and professional use.
