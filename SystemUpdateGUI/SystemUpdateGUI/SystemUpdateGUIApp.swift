import SwiftUI

@main
struct SystemUpdateGUIApp: App {
    @StateObject private var updateManager = UpdateManager()

    var body: some Scene {
        MenuBarExtra {
            MenuBarView()
                .environmentObject(updateManager)
        } label: {
            Label {
                Text("System Update")
            } icon: {
                if updateManager.isRunning {
                    Image(systemName: "arrow.triangle.2.circlepath")
                        .symbolEffect(.rotate)
                } else if updateManager.lastRunFailed {
                    Image(systemName: "exclamationmark.triangle.fill")
                } else {
                    Image(systemName: "arrow.triangle.2.circlepath")
                }
            }
        }
        .menuBarExtraStyle(.window)

        Settings {
            SettingsView()
                .environmentObject(updateManager)
        }
    }
}
