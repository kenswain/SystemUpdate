import Foundation

struct OutputParser {
    struct ParsedLine {
        let text: String
        let type: OutputLine.LineType
    }

    /// Parse a line of output from system-update CLI
    /// The CLI uses ANSI color codes and symbols to indicate status
    static func parse(_ rawLine: String) -> ParsedLine {
        // Strip ANSI escape codes
        let stripped = stripANSI(rawLine)
        let trimmed = stripped.trimmingCharacters(in: .whitespaces)

        guard !trimmed.isEmpty else {
            return ParsedLine(text: trimmed, type: .detail)
        }

        // Detect line type based on prefix symbols used by the Go CLI
        if trimmed.hasPrefix("✓") || trimmed.contains("completed successfully") {
            return ParsedLine(text: trimmed, type: .success)
        } else if trimmed.hasPrefix("✗") || trimmed.contains("failed") || trimmed.contains("error") {
            return ParsedLine(text: trimmed, type: .error)
        } else if trimmed.hasPrefix("→") || trimmed.hasPrefix("Updating") || trimmed.hasPrefix("Upgrading") || trimmed.hasPrefix("Cleaning") {
            return ParsedLine(text: trimmed, type: .progress)
        } else if trimmed.hasPrefix("⚠") || trimmed.contains("warning") || trimmed.contains("not found") || trimmed.contains("skipping") {
            return ParsedLine(text: trimmed, type: .warning)
        } else if trimmed.hasPrefix("•") || trimmed.hasPrefix("  ") {
            return ParsedLine(text: trimmed, type: .detail)
        } else if trimmed.hasPrefix("===") || trimmed.contains("Summary") {
            return ParsedLine(text: trimmed, type: .info)
        } else if trimmed.hasPrefix("Starting") {
            return ParsedLine(text: trimmed, type: .info)
        } else if trimmed.hasPrefix("Total time:") {
            return ParsedLine(text: trimmed, type: .info)
        } else {
            return ParsedLine(text: trimmed, type: .detail)
        }
    }

    /// Strip ANSI escape sequences from a string
    private static func stripANSI(_ string: String) -> String {
        // Match ANSI escape sequences: ESC [ ... m
        let pattern = "\\x1B\\[[0-9;]*[a-zA-Z]"
        guard let regex = try? NSRegularExpression(pattern: pattern) else {
            return string
        }
        let range = NSRange(string.startIndex..., in: string)
        return regex.stringByReplacingMatches(in: string, range: range, withTemplate: "")
    }
}
