//
//  ClientConfig.swift
//  GameProxyApp
//
//  Client configuration for connecting to proxy server
//

import Foundation

struct ProxyServerConfig {
    static let shared = ProxyServerConfig()
    
    // Server details
    let serverIPv4 = "194.233.73.36"
    let serverIPv6 = "2407:3640:2299:7804::1"
    let httpPort = 8080
    let socks5Port = 1080
    
    // Connection URLs
    var httpProxyURL: String {
        return "http://\(serverIPv4):\(httpPort)"
    }
    
    var httpProxyURLIPv6: String {
        return "http://[\(serverIPv6)]:\(httpPort)"
    }
    
    var socks5ProxyURL: String {
        return "socks5://\(serverIPv4):\(socks5Port)"
    }
    
    // Use IPv6 if available, fallback to IPv4
    var preferredProxyURL: String {
        // Check IPv6 availability (simplified)
        if isIPv6Available() {
            return httpProxyURLIPv6
        }
        return httpProxyURL
    }
    
    private func isIPv6Available() -> Bool {
        // Simple check - in production, use proper network detection
        return false // Default to IPv4 for reliability
    }
}

// MARK: - Proxy Client Manager

class ProxyClientManager: ObservableObject {
    @Published var isConnected: Bool = false
    @Published var latency: Int = 0
    @Published var serverConfig = ProxyServerConfig.shared
    
    private var testTimer: Timer?
    
    // Connect to proxy server
    func connect() {
        // Configure system proxy to use our server
        configureSystemProxy()
        
        // Start connection monitoring
        startMonitoring()
        
        isConnected = true
    }
    
    // Disconnect from proxy server
    func disconnect() {
        // Remove system proxy configuration
        removeSystemProxy()
        
        // Stop monitoring
        stopMonitoring()
        
        isConnected = false
    }
    
    // Configure system proxy to route through server
    private func configureSystemProxy() {
        let config = serverConfig
        
        // Use networksetup to configure system proxy
        let script = """
        networksetup -setwebproxy "Wi-Fi" \(config.serverIPv4) \(config.httpPort)
        networksetup -setsecurewebproxy "Wi-Fi" \(config.serverIPv4) \(config.httpPort)
        networksetup -setwebproxystate "Wi-Fi" on
        networksetup -setsecurewebproxystate "Wi-Fi" on
        """
        
        // Execute with admin privileges
        executeWithAdmin(script: script)
    }
    
    // Remove system proxy configuration
    private func removeSystemProxy() {
        let script = """
        networksetup -setwebproxystate "Wi-Fi" off
        networksetup -setsecurewebproxystate "Wi-Fi" off
        """
        
        executeWithAdmin(script: script)
    }
    
    // Start monitoring connection health
    private func startMonitoring() {
        testTimer = Timer.scheduledTimer(withTimeInterval: 5.0, repeats: true) { [weak self] _ in
            self?.testConnection()
        }
    }
    
    // Stop monitoring
    private func stopMonitoring() {
        testTimer?.invalidate()
        testTimer = nil
    }
    
    // Test connection to proxy server
    private func testConnection() {
        let startTime = Date()
        let testURL = URL(string: "http://\(serverConfig.serverIPv4):\(serverConfig.httpPort)")!
        
        var request = URLRequest(url: testURL)
        request.timeoutInterval = 5.0
        
        URLSession.shared.dataTask(with: request) { [weak self] _, response, error in
            let elapsed = Int(Date().timeIntervalSince(startTime) * 1000)
            
            DispatchQueue.main.async {
                if error == nil {
                    self?.isConnected = true
                    self?.latency = elapsed
                } else {
                    self?.isConnected = false
                }
            }
        }.resume()
    }
    
    // Execute script with admin privileges
    private func executeWithAdmin(script: String) {
        // Note: In production, use Authorization Services API
        // This is a simplified example
        let task = Process()
        task.launchPath = "/bin/bash"
        task.arguments = ["-c", script]
        task.launch()
    }
}

// MARK: - Quick Connect Helper

extension ProxyClientManager {
    // Quick connect with default server
    func quickConnect() {
        connect()
    }
    
    // Connect with custom server
    func connectToServer(ip: String, port: Int) {
        // Update server config
        // Then connect
        connect()
    }
    
    // Test server connection without applying system proxy
    func testServer(ip: String, port: Int, completion: @escaping (Bool, Int) -> Void) {
        let startTime = Date()
        let testURL = URL(string: "http://\(ip):\(port)")!
        
        var request = URLRequest(url: testURL)
        request.timeoutInterval = 5.0
        
        URLSession.shared.dataTask(with: request) { _, response, error in
            let elapsed = Int(Date().timeIntervalSince(startTime) * 1000)
            let success = error == nil && (response as? HTTPURLResponse)?.statusCode == 200
            
            DispatchQueue.main.async {
                completion(success, elapsed)
            }
        }.resume()
    }
}
