package output

import (
	"fmt"
	"os"
)

// VerboseWriter writes command output to stderr in verbose mode
type VerboseWriter struct{}

// Write implements io.Writer interface
func (v VerboseWriter) Write(p []byte) (n int, err error) {
	// Write to stderr so it doesn't interfere with structured output
	return os.Stderr.Write(p)
}

// VerbosePrint prints a message only in verbose mode
func VerbosePrint(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
