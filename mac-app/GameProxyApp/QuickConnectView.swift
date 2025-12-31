//
//  QuickConnectView.swift
//  GameProxyApp
//
//  Quick connect interface for proxy server
//

import SwiftUI

struct QuickConnectView: View {
    @ObservedObject var clientManager = ProxyClientManager()
    @State private var showCustomServer = false
    @State private var customIP = "194.233.73.36"
    @State private var customPort = "8080"
    @State private var testResult: String = ""
    
    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            // Header
            Text("Connect to Proxy Server")
                .font(.title2)
                .fontWeight(.bold)
                .padding(.bottom, 10)
            
            // Server Info
            VStack(alignment: .leading, spacing: 8) {
                Text("Server Information")
                    .font(.headline)
                
                HStack {
                    Text("IPv4:")
                    Text("194.233.73.36:8080")
                        .font(.system(.body, design: .monospaced))
                        .foregroundColor(.secondary)
                }
                
                HStack {
                    Text("IPv6:")
                    Text("[2407:3640:2299:7804::1]:8080")
                        .font(.system(.body, design: .monospaced))
                        .foregroundColor(.secondary)
                }
            }
            .padding()
            .background(Color(NSColor.controlBackgroundColor))
            .cornerRadius(8)
            
            // Connection Status
            HStack {
                Circle()
                    .fill(clientManager.isConnected ? Color.green : Color.red)
                    .frame(width: 12, height: 12)
                
                Text(clientManager.isConnected ? "Connected" : "Disconnected")
                    .font(.headline)
                
                Spacer()
                
                if clientManager.latency > 0 {
                    Text("\(clientManager.latency)ms")
                        .font(.system(.body, design: .monospaced))
                        .foregroundColor(.secondary)
                }
            }
            .padding()
            .background(Color(NSColor.controlBackgroundColor))
            .cornerRadius(8)
            
            // Test Result
            if !testResult.isEmpty {
                Text(testResult)
                    .foregroundColor(testResult.contains("✅") ? .green : .red)
                    .padding(.horizontal)
            }
            
            // Buttons
            HStack(spacing: 12) {
                Button(action: {
                    if clientManager.isConnected {
                        clientManager.disconnect()
                    } else {
                        clientManager.connect()
                    }
                }) {
                    Text(clientManager.isConnected ? "Disconnect" : "Connect")
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.borderedProminent)
                .controlSize(.large)
                
                Button(action: {
                    testServer()
                }) {
                    Text("Test")
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.bordered)
                .controlSize(.large)
            }
            
            // Custom Server
            Button(action: {
                showCustomServer.toggle()
            }) {
                HStack {
                    Image(systemName: showCustomServer ? "chevron.down" : "chevron.right")
                    Text("Custom Server")
                }
            }
            .buttonStyle(.plain)
            
            if showCustomServer {
                VStack(alignment: .leading, spacing: 12) {
                    HStack {
                        Text("IP Address:")
                        TextField("194.233.73.36", text: $customIP)
                            .textFieldStyle(.roundedBorder)
                    }
                    
                    HStack {
                        Text("Port:")
                        TextField("8080", text: $customPort)
                            .textFieldStyle(.roundedBorder)
                            .frame(width: 100)
                    }
                    
                    Button(action: {
                        if let port = Int(customPort) {
                            clientManager.connectToServer(ip: customIP, port: port)
                        }
                    }) {
                        Text("Connect to Custom Server")
                            .frame(maxWidth: .infinity)
                    }
                    .buttonStyle(.bordered)
                }
                .padding()
                .background(Color(NSColor.controlBackgroundColor))
                .cornerRadius(8)
            }
            
            Spacer()
        }
        .padding()
        .frame(width: 400, height: 500)
    }
    
    private func testServer() {
        testResult = "Testing..."
        clientManager.testServer(ip: customIP, port: Int(customPort) ?? 8080) { success, latency in
            if success {
                testResult = "✅ Connection successful! Latency: \(latency)ms"
            } else {
                testResult = "❌ Connection failed"
            }
        }
    }
}

#Preview {
    QuickConnectView()
}
