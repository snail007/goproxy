//
//  ProxyManager.swift
//  GameProxyApp
//
//  Manages the proxy service lifecycle and configuration
//

import Foundation
import Network

class ProxyManager: ObservableObject {
    @Published var isRunning: Bool = false
    @Published var isConnected: Bool = false
    @Published var latency: Int = 0
    @Published var stealthMode: Bool = true
    @Published var gameMode: Bool = true
    
    private var proxyProcess: Process?
    private var monitor: NWPathMonitor?
    private let queue = DispatchQueue(label: "ProxyManager")
    
    static let shared = ProxyManager()
    
    private init() {
        setupNetworkMonitor()
    }
    
    // MARK: - Network Monitoring
    
    private func setupNetworkMonitor() {
        monitor = NWPathMonitor()
        monitor?.pathUpdateHandler = { [weak self] path in
            DispatchQueue.main.async {
                self?.isConnected = path.status == .satisfied
            }
        }
        monitor?.start(queue: queue)
    }
    
    // MARK: - Proxy Control
    
    func startProxy(port: Int = 8080, parentProxy: String? = nil) {
        guard !isRunning else { return }
        
        let process = Process()
        process.executableURL = Bundle.main.url(forResource: "proxy", withExtension: nil)
        
        var arguments = ["http"]
        arguments.append("-t")
        arguments.append("tcp")
        arguments.append("-p")
        arguments.append(":\(port)")
        
        if let parent = parentProxy {
            arguments.append("-T")
            arguments.append("tcp")
            arguments.append("-P")
            arguments.append(parent)
        }
        
        process.arguments = arguments
        process.standardOutput = Pipe()
        process.standardError = Pipe()
        
        do {
            try process.run()
            proxyProcess = process
            isRunning = true
            
            // Configure system proxy
            configureSystemProxy(port: port)
            
            // Start latency monitoring
            startLatencyMonitoring()
            
        } catch {
            print("Failed to start proxy: \(error)")
        }
    }
    
    func stopProxy() {
        guard isRunning else { return }
        
        proxyProcess?.terminate()
        proxyProcess = nil
        isRunning = false
        
        // Remove system proxy configuration
        removeSystemProxy()
    }
    
    // MARK: - System Proxy Configuration
    
    private func configureSystemProxy(port: Int) {
        // Create PAC file
        let pacContent = """
        function FindProxyForURL(url, host) {
            return "PROXY 127.0.0.1:\(port); DIRECT";
        }
        """
        
        let pacURL = FileManager.default.temporaryDirectory.appendingPathComponent("gameproxy.pac")
        try? pacContent.write(to: pacURL, atomically: true, encoding: .utf8)
        
        // Configure system proxy using networksetup (requires admin)
        let script = """
        networksetup -setwebproxystate "Wi-Fi" on
        networksetup -setwebproxy "Wi-Fi" 127.0.0.1 \(port)
        networksetup -setsecurewebproxystate "Wi-Fi" on
        networksetup -setsecurewebproxy "Wi-Fi" 127.0.0.1 \(port)
        networksetup -setautoproxyurl "Wi-Fi" file://\(pacURL.path)
        """
        
        // Note: This requires admin privileges
        // In production, use Authorization Services API
    }
    
    private func removeSystemProxy() {
        let script = """
        networksetup -setwebproxystate "Wi-Fi" off
        networksetup -setsecurewebproxystate "Wi-Fi" off
        networksetup -setautoproxyurl "Wi-Fi" ""
        """
        
        // Note: This requires admin privileges
    }
    
    // MARK: - Latency Monitoring
    
    private func startLatencyMonitoring() {
        Timer.scheduledTimer(withTimeInterval: 5.0, repeats: true) { [weak self] _ in
            self?.measureLatency()
        }
    }
    
    private func measureLatency() {
        let startTime = Date()
        
        // Simple latency measurement
        URLSession.shared.dataTask(with: URL(string: "http://127.0.0.1:8080")!) { _, _, _ in
            let elapsed = Int(Date().timeIntervalSince(startTime) * 1000)
            DispatchQueue.main.async {
                self.latency = elapsed
            }
        }.resume()
    }
    
    // MARK: - Settings
    
    func setStealthMode(_ enabled: Bool) {
        stealthMode = enabled
        // Send command to proxy process to enable/disable stealth
    }
    
    func setGameMode(_ enabled: Bool) {
        gameMode = enabled
        // Send command to proxy process to enable/disable game mode
    }
}
