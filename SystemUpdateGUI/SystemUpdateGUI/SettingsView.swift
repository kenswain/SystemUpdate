import SwiftUI
import ServiceManagement

struct SettingsView: View {
    @EnvironmentObject var updateManager: UpdateManager
    @AppStorage("brewOnly") private var brewOnly = false
    @AppStorage("npmOnly") private var npmOnly = false
    @AppStorage("verboseOutput") private var verboseOutput = false
    @State private var launchAtLogin = false
    @State private var loginItemError: String?

    var body: some View {
        Form {
            Section {
                VStack(alignment: .leading, spacing: 12) {
                    Text("Update Scope")
                        .font(.headline)

                    Picker("Update Mode", selection: updateModeBinding) {
                        Text("Update Both").tag(UpdateMode.both)
                        Text("Homebrew Only").tag(UpdateMode.brewOnly)
                        Text("npm Only").tag(UpdateMode.npmOnly)
                    }
                    .pickerStyle(.radioGroup)

                    Text("Choose which package managers to update when running.")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }

            Divider()
                .padding(.vertical, 8)

            Section {
                VStack(alignment: .leading, spacing: 12) {
                    Text("Output")
                        .font(.headline)

                    Toggle("Verbose output", isOn: $verboseOutput)

                    Text("Show detailed command output during updates.")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }

            Divider()
                .padding(.vertical, 8)

            Section {
                VStack(alignment: .leading, spacing: 12) {
                    Text("General")
                        .font(.headline)

                    Toggle("Launch at login", isOn: $launchAtLogin)
                        .onChange(of: launchAtLogin) { _, newValue in
                            setLaunchAtLogin(enabled: newValue)
                        }

                    Text("Start System Update automatically when you log in.")
                        .font(.caption)
                        .foregroundColor(.secondary)

                    if let error = loginItemError {
                        Text(error)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                }
            }
            .onAppear {
                launchAtLogin = SMAppService.mainApp.status == .enabled
            }

            Divider()
                .padding(.vertical, 8)

            Section {
                VStack(alignment: .leading, spacing: 8) {
                    Text("About")
                        .font(.headline)

                    HStack {
                        Text("Version")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.secondary)
                    }

                    HStack {
                        Text("CLI Binary")
                        Spacer()
                        Text(binaryStatus)
                            .foregroundColor(binaryFound ? .green : .red)
                    }

                    HStack {
                        Text("Login Item")
                        Spacer()
                        Text(loginItemStatus)
                            .foregroundColor(loginItemStatusColor)
                    }
                }
            }
        }
        .formStyle(.grouped)
        .frame(width: 400, height: 400)
        .navigationTitle("Settings")
    }

    private enum UpdateMode {
        case both, brewOnly, npmOnly
    }

    private var updateModeBinding: Binding<UpdateMode> {
        Binding(
            get: {
                if brewOnly { return .brewOnly }
                if npmOnly { return .npmOnly }
                return .both
            },
            set: { newValue in
                switch newValue {
                case .both:
                    brewOnly = false
                    npmOnly = false
                case .brewOnly:
                    brewOnly = true
                    npmOnly = false
                case .npmOnly:
                    brewOnly = false
                    npmOnly = true
                }
            }
        )
    }

    private var binaryFound: Bool {
        // Check bundled location first
        if let resourcePath = Bundle.main.resourcePath {
            let bundledPath = (resourcePath as NSString).appendingPathComponent("system-update")
            if FileManager.default.fileExists(atPath: bundledPath) {
                return true
            }
        }

        // Check installed location
        return FileManager.default.fileExists(atPath: "/usr/local/bin/system-update")
    }

    private var binaryStatus: String {
        if let resourcePath = Bundle.main.resourcePath {
            let bundledPath = (resourcePath as NSString).appendingPathComponent("system-update")
            if FileManager.default.fileExists(atPath: bundledPath) {
                return "Bundled"
            }
        }

        if FileManager.default.fileExists(atPath: "/usr/local/bin/system-update") {
            return "Installed"
        }

        return "Not Found"
    }

    private var loginItemStatus: String {
        switch SMAppService.mainApp.status {
        case .enabled:
            return "Enabled"
        case .notRegistered:
            return "Disabled"
        case .requiresApproval:
            return "Requires Approval"
        case .notFound:
            return "Not Found"
        @unknown default:
            return "Unknown"
        }
    }

    private var loginItemStatusColor: Color {
        switch SMAppService.mainApp.status {
        case .enabled:
            return .green
        case .notRegistered:
            return .secondary
        case .requiresApproval:
            return .orange
        default:
            return .red
        }
    }

    private func setLaunchAtLogin(enabled: Bool) {
        loginItemError = nil

        do {
            if enabled {
                try SMAppService.mainApp.register()
            } else {
                try SMAppService.mainApp.unregister()
            }
        } catch {
            loginItemError = "Failed to update login item: \(error.localizedDescription)"
            // Revert the toggle state
            launchAtLogin = !enabled
        }
    }
}

#Preview {
    SettingsView()
        .environmentObject(UpdateManager())
}
