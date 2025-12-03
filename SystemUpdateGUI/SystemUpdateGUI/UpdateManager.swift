import Foundation
import SwiftUI

struct OutputLine: Identifiable {
    let id = UUID()
    let text: String
    let type: LineType
    let timestamp: Date

    enum LineType {
        case info
        case success
        case error
        case warning
        case progress
        case detail
    }

    var icon: String? {
        switch type {
        case .success: return "checkmark.circle.fill"
        case .error: return "xmark.circle.fill"
        case .warning: return "exclamationmark.triangle.fill"
        case .progress: return "arrow.right.circle"
        case .info: return "info.circle"
        case .detail: return nil
        }
    }

    var color: Color {
        switch type {
        case .success: return .green
        case .error: return .red
        case .warning: return .orange
        case .progress: return .blue
        case .info: return .cyan
        case .detail: return .secondary
        }
    }
}

@MainActor
class UpdateManager: ObservableObject {
    @Published var isRunning = false
    @Published var outputLines: [OutputLine] = []
    @Published var lastRunTime: Date?
    @Published var lastRunFailed = false

    @AppStorage("brewOnly") var brewOnly = false
    @AppStorage("npmOnly") var npmOnly = false
    @AppStorage("verboseOutput") var verboseOutput = false

    private var currentProcess: Process?
    private var outputPipe: Pipe?

    private var binaryPath: String {
        // Look for binary in app bundle Resources
        if let resourcePath = Bundle.main.resourcePath {
            let bundledPath = (resourcePath as NSString).appendingPathComponent("system-update")
            if FileManager.default.fileExists(atPath: bundledPath) {
                return bundledPath
            }
        }

        // Fallback to installed binary
        let installedPath = "/usr/local/bin/system-update"
        if FileManager.default.fileExists(atPath: installedPath) {
            return installedPath
        }

        // Development fallback - look in parent directory
        let developmentPath = Bundle.main.bundleURL
            .deletingLastPathComponent()
            .deletingLastPathComponent()
            .deletingLastPathComponent()
            .deletingLastPathComponent()
            .appendingPathComponent("system-update")
            .path

        return developmentPath
    }

    func runUpdate() {
        guard !isRunning else { return }

        isRunning = true
        lastRunFailed = false
        outputLines = []

        addLine("Starting system update...", type: .info)

        let process = Process()
        process.executableURL = URL(fileURLWithPath: binaryPath)

        var arguments: [String] = []
        if brewOnly {
            arguments.append("--brew-only")
        } else if npmOnly {
            arguments.append("--npm-only")
        }
        if verboseOutput {
            arguments.append("--verbose")
        }
        process.arguments = arguments

        // Set up environment to force color output (we'll parse ANSI codes)
        var environment = ProcessInfo.processInfo.environment
        environment["TERM"] = "xterm-256color"
        process.environment = environment

        let pipe = Pipe()
        process.standardOutput = pipe
        process.standardError = pipe
        outputPipe = pipe

        // Handle output asynchronously
        pipe.fileHandleForReading.readabilityHandler = { [weak self] handle in
            let data = handle.availableData
            if data.isEmpty { return }

            if let output = String(data: data, encoding: .utf8) {
                Task { @MainActor in
                    self?.processOutput(output)
                }
            }
        }

        process.terminationHandler = { [weak self] process in
            Task { @MainActor in
                self?.handleTermination(exitCode: process.terminationStatus)
            }
        }

        currentProcess = process

        do {
            try process.run()
        } catch {
            addLine("Failed to start: \(error.localizedDescription)", type: .error)
            isRunning = false
            lastRunFailed = true
        }
    }

    func cancel() {
        currentProcess?.terminate()
        addLine("Update cancelled by user", type: .warning)
    }

    private func processOutput(_ output: String) {
        let lines = output.components(separatedBy: .newlines)

        for line in lines {
            guard !line.isEmpty else { continue }

            let parsed = OutputParser.parse(line)
            addLine(parsed.text, type: parsed.type)
        }
    }

    private func handleTermination(exitCode: Int32) {
        outputPipe?.fileHandleForReading.readabilityHandler = nil

        if exitCode == 0 {
            addLine("Update completed successfully!", type: .success)
        } else if currentProcess?.terminationReason == .exit {
            addLine("Update finished with errors (exit code: \(exitCode))", type: .error)
            lastRunFailed = true
        }

        isRunning = false
        lastRunTime = Date()
        currentProcess = nil
        outputPipe = nil
    }

    private func addLine(_ text: String, type: OutputLine.LineType) {
        let line = OutputLine(text: text, type: type, timestamp: Date())
        outputLines.append(line)
    }
}
