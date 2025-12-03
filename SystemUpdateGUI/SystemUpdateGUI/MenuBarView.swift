import SwiftUI

struct MenuBarView: View {
    @EnvironmentObject var updateManager: UpdateManager
    @State private var showingSettings = false

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // Header
            HStack {
                Image(systemName: "arrow.triangle.2.circlepath")
                    .font(.title2)
                    .foregroundColor(.accentColor)
                Text("System Update")
                    .font(.headline)
                Spacer()
                if updateManager.isRunning {
                    ProgressView()
                        .scaleEffect(0.7)
                }
            }
            .padding()

            Divider()

            // Status and output area
            ScrollViewReader { proxy in
                ScrollView {
                    VStack(alignment: .leading, spacing: 4) {
                        ForEach(updateManager.outputLines) { line in
                            OutputLineView(line: line)
                                .id(line.id)
                        }

                        if updateManager.outputLines.isEmpty && !updateManager.isRunning {
                            Text("Click 'Run Update' to check for updates")
                                .foregroundColor(.secondary)
                                .font(.caption)
                                .padding(.vertical, 20)
                                .frame(maxWidth: .infinity)
                        }
                    }
                    .padding(.horizontal)
                    .padding(.vertical, 8)
                }
                .frame(height: 200)
                .onChange(of: updateManager.outputLines.count) { _ in
                    if let lastLine = updateManager.outputLines.last {
                        withAnimation {
                            proxy.scrollTo(lastLine.id, anchor: .bottom)
                        }
                    }
                }
            }

            Divider()

            // Action buttons
            HStack {
                // Update mode indicator
                if updateManager.brewOnly {
                    Label("Homebrew only", systemImage: "mug")
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else if updateManager.npmOnly {
                    Label("npm only", systemImage: "shippingbox")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }

                Spacer()

                if #available(macOS 14.0, *) {
                    SettingsLink {
                        Image(systemName: "gear")
                    }
                    .buttonStyle(.borderless)
                    .help("Settings")
                } else {
                    Button {
                        NSApp.sendAction(Selector(("showSettingsWindow:")), to: nil, from: nil)
                    } label: {
                        Image(systemName: "gear")
                    }
                    .buttonStyle(.borderless)
                    .help("Settings")
                }

                Button(action: {
                    if updateManager.isRunning {
                        updateManager.cancel()
                    } else {
                        updateManager.runUpdate()
                    }
                }) {
                    if updateManager.isRunning {
                        Label("Cancel", systemImage: "xmark.circle")
                    } else {
                        Label("Run Update", systemImage: "play.fill")
                    }
                }
                .buttonStyle(.borderedProminent)
                .controlSize(.small)
                .disabled(false)
            }
            .padding()

            Divider()

            // Footer
            HStack {
                if let lastRun = updateManager.lastRunTime {
                    Text("Last run: \(lastRun, style: .relative) ago")
                        .font(.caption2)
                        .foregroundColor(.secondary)
                }
                Spacer()
                Button("Quit") {
                    NSApplication.shared.terminate(nil)
                }
                .buttonStyle(.borderless)
                .font(.caption)
            }
            .padding(.horizontal)
            .padding(.vertical, 8)
        }
        .frame(width: 380)
    }
}

struct OutputLineView: View {
    let line: OutputLine

    var body: some View {
        HStack(alignment: .top, spacing: 6) {
            if let icon = line.icon {
                Image(systemName: icon)
                    .font(.caption)
                    .foregroundColor(line.color)
                    .frame(width: 14)
            } else {
                Color.clear.frame(width: 14)
            }

            Text(line.text)
                .font(.system(.caption, design: .monospaced))
                .foregroundColor(line.color)
                .textSelection(.enabled)

            Spacer()
        }
    }
}

#Preview {
    MenuBarView()
        .environmentObject(UpdateManager())
        .frame(width: 380)
}
