#!/usr/bin/expect -f

# Direct deployment script using expect
# Usage: ./deploy-direct.sh

set timeout 300
set server_ip "194.233.73.36"
set server_user "root"
set server_pass "JVxgREA3hbe6"
set server_ipv6 "2407:3640:2299:7804::1"

puts "=========================================="
puts "Deploying Game Proxy to Server"
puts "=========================================="
puts "Server: $server_ip"
puts ""

# Connect and deploy
spawn ssh -o StrictHostKeyChecking=no ${server_user}@${server_ip}

expect {
    "password:" {
        send "${server_pass}\r"
    }
    timeout {
        puts "Connection timeout"
        exit 1
    }
}

expect "# "

# Update system
puts "Updating system..."
send "apt-get update -qq && apt-get install -y wget curl git ufw iptables net-tools\r"
expect "# "

# Install Go
puts "Installing Go..."
send "if ! command -v go &> /dev/null; then wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz -O /tmp/go.tar.gz && tar -C /usr/local -xzf /tmp/go.tar.gz && rm /tmp/go.tar.gz && echo 'export PATH=\\$PATH:/usr/local/go/bin' >> /root/.bashrc; fi\r"
expect "# "
send "export PATH=\\$PATH:/usr/local/go/bin\r"
expect "# "
send "/usr/local/go/bin/go version\r"
expect "# "

# Create directories
puts "Creating directories..."
send "mkdir -p /opt/game-proxy/{bin,config,logs}\r"
expect "# "

# Upload and extract files (we'll do this via separate SCP)
puts "Files need to be uploaded separately..."
puts "Please run: scp /tmp/proxy-deploy.tar.gz root@${server_ip}:/tmp/"

send "exit\r"
expect eof

puts ""
puts "=========================================="
puts "Please upload files manually:"
puts "  scp /tmp/proxy-deploy.tar.gz root@${server_ip}:/tmp/"
puts "Then continue with deployment..."
puts "=========================================="
