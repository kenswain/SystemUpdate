import SwiftUI

struct MenuBarIcon: View {
    let isRunning: Bool
    let lastRunFailed: Bool

    var body: some View {
        Group {
            if isRunning {
                if #available(macOS 15.0, *) {
                    Image(systemName: "arrow.triangle.2.circlepath")
                        .symbolEffect(.rotate, isActive: true)
                } else {
                    // Fallback for macOS 13-14: use a different indicator
                    Image(systemName: "arrow.triangle.2.circlepath.circle.fill")
                }
            } else if lastRunFailed {
                Image(systemName: "exclamationmark.triangle.fill")
            } else {
                Image(systemName: "arrow.triangle.2.circlepath")
            }
        }
    }
}

@main
struct SystemUpdateGUIApp: App {
    @StateObject private var updateManager = UpdateManager()

    var body: some Scene {
        MenuBarExtra {
            MenuBarView()
                .environmentObject(updateManager)
        } label: {
            MenuBarIcon(isRunning: updateManager.isRunning, lastRunFailed: updateManager.lastRunFailed)
        }
        .menuBarExtraStyle(.window)

        Settings {
            SettingsView()
                .environmentObject(updateManager)
        }
    }
}
