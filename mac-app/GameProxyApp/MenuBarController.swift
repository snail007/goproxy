//
//  MenuBarController.swift
//  GameProxyApp
//
//  Controls the menu bar interface
//

import AppKit
import SwiftUI

class MenuBarController {
    private var statusItem: NSStatusItem?
    private let proxyManager = ProxyManager.shared
    
    init() {
        setupMenuBar()
        setupObservers()
    }
    
    private func setupMenuBar() {
        statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
        
        if let button = statusItem?.button {
            button.image = NSImage(systemSymbolName: "network", accessibilityDescription: "Game Proxy")
            button.action = #selector(menuBarButtonClicked)
            button.target = self
        }
        
        updateStatusIcon()
    }
    
    private func setupObservers() {
        // Observe proxy state changes
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(proxyStateChanged),
            name: NSNotification.Name("ProxyStateChanged"),
            object: nil
        )
    }
    
    @objc private func menuBarButtonClicked() {
        guard let button = statusItem?.button else { return }
        
        let menu = NSMenu()
        
        // Status
        let statusItem = NSMenuItem()
        if proxyManager.isRunning {
            statusItem.title = "Status: Running"
        } else {
            statusItem.title = "Status: Stopped"
        }
        statusItem.isEnabled = false
        menu.addItem(statusItem)
        
        // Latency
        let latencyItem = NSMenuItem()
        latencyItem.title = "Latency: \(proxyManager.latency)ms"
        latencyItem.isEnabled = false
        menu.addItem(latencyItem)
        
        menu.addItem(NSMenuItem.separator())
        
        // Toggle
        let toggleItem = NSMenuItem()
        toggleItem.title = proxyManager.isRunning ? "Stop Proxy" : "Start Proxy"
        toggleItem.action = #selector(toggleProxy)
        toggleItem.target = self
        menu.addItem(toggleItem)
        
        menu.addItem(NSMenuItem.separator())
        
        // Settings
        let settingsItem = NSMenuItem()
        settingsItem.title = "Settings..."
        settingsItem.action = #selector(showSettings)
        settingsItem.target = self
        menu.addItem(settingsItem)
        
        // Quit
        let quitItem = NSMenuItem()
        quitItem.title = "Quit"
        quitItem.action = #selector(quitApp)
        quitItem.target = self
        menu.addItem(quitItem)
        
        statusItem?.menu = menu
        statusItem?.popUpMenu(menu)
    }
    
    @objc private func toggleProxy() {
        if proxyManager.isRunning {
            proxyManager.stopProxy()
        } else {
            proxyManager.startProxy()
        }
        updateStatusIcon()
    }
    
    @objc private func showSettings() {
        // Show settings window
        let settingsWindow = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 400, height: 300),
            styleMask: [.titled, .closable],
            backing: .buffered,
            defer: false
        )
        
        let settingsView = SettingsView()
        let hostingView = NSHostingView(rootView: settingsView)
        settingsWindow.contentView = hostingView
        settingsWindow.center()
        settingsWindow.makeKeyAndOrderFront(nil)
    }
    
    @objc private func quitApp() {
        proxyManager.stopProxy()
        NSApplication.shared.terminate(nil)
    }
    
    @objc private func proxyStateChanged() {
        updateStatusIcon()
    }
    
    private func updateStatusIcon() {
        guard let button = statusItem?.button else { return }
        
        if proxyManager.isRunning {
            if proxyManager.isConnected {
                button.image = NSImage(systemSymbolName: "network", accessibilityDescription: "Connected")
                button.contentTintColor = .systemGreen
            } else {
                button.image = NSImage(systemSymbolName: "network.slash", accessibilityDescription: "Disconnected")
                button.contentTintColor = .systemRed
            }
        } else {
            button.image = NSImage(systemSymbolName: "network", accessibilityDescription: "Stopped")
            button.contentTintColor = .systemGray
        }
    }
}

// Settings View
struct SettingsView: View {
    @ObservedObject var proxyManager = ProxyManager.shared
    
    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("Game Proxy Settings")
                .font(.title2)
                .padding()
            
            Toggle("Stealth Mode", isOn: $proxyManager.stealthMode)
                .padding(.horizontal)
            
            Toggle("Game Mode", isOn: $proxyManager.gameMode)
                .padding(.horizontal)
            
            HStack {
                Text("Latency: \(proxyManager.latency)ms")
                Spacer()
            }
            .padding(.horizontal)
            
            Spacer()
            
            HStack {
                Spacer()
                Button("Close") {
                    NSApplication.shared.keyWindow?.close()
                }
            }
            .padding()
        }
        .frame(width: 400, height: 300)
    }
}
