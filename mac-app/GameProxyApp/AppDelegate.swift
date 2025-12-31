//
//  AppDelegate.swift
//  GameProxyApp
//
//  Main application delegate
//

import Cocoa
import SwiftUI

@main
class AppDelegate: NSObject, NSApplicationDelegate {
    private var menuBarController: MenuBarController?
    
    func applicationDidFinishLaunching(_ aNotification: Notification) {
        // Hide dock icon - run as menu bar only app
        NSApp.setActivationPolicy(.accessory)
        
        // Initialize menu bar controller
        menuBarController = MenuBarController()
        
        // Auto-start proxy if configured
        if UserDefaults.standard.bool(forKey: "autoStart") {
            ProxyManager.shared.startProxy()
        }
    }
    
    func applicationWillTerminate(_ aNotification: Notification) {
        ProxyManager.shared.stopProxy()
    }
    
    func applicationShouldHandleReopen(_ sender: NSApplication, hasVisibleWindows flag: Bool) -> Bool {
        // Don't show windows on reopen - menu bar only
        return false
    }
}
