#!/bin/bash

# Complete deployment script
# This will upload files and deploy to server

SERVER_IP="194.233.73.36"
SERVER_USER="root"
SERVER_PASS="JVxgREA3hbe6"

echo "=========================================="
echo "Deploying Game Proxy to Server"
echo "Server: $SERVER_IP"
echo "=========================================="

# Create deployment package
echo "Step 1: Creating deployment package..."
cd /workspace
tar -czf /tmp/proxy-deploy.tar.gz stealth/ game/ services/game_http.go main.go config.go go.mod
echo "✅ Package created: $(ls -lh /tmp/proxy-deploy.tar.gz | awk '{print $5}')"

# Upload using expect
echo "Step 2: Uploading files to server..."
expect <<EOF
set timeout 60
spawn scp -o StrictHostKeyChecking=no /tmp/proxy-deploy.tar.gz deploy/deploy-on-server.sh ${SERVER_USER}@${SERVER_IP}:/tmp/
expect {
    "password:" {
        send "${SERVER_PASS}\r"
        exp_continue
    }
    "100%" {
        puts "Upload successful"
    }
    timeout {
        puts "Upload timeout"
        exit 1
    }
}
expect eof
EOF

if [ $? -eq 0 ]; then
    echo "✅ Files uploaded"
else
    echo "❌ Upload failed, trying alternative method..."
    # Alternative: provide manual instructions
    echo "Please run manually:"
    echo "  scp /tmp/proxy-deploy.tar.gz root@${SERVER_IP}:/tmp/"
    exit 1
fi

# Execute deployment on server
echo "Step 3: Executing deployment on server..."
expect <<'DEPLOY_EOF'
set timeout 600
set server_ip "194.233.73.36"
set server_user "root"
set server_pass "JVxgREA3hbe6"

spawn ssh -o StrictHostKeyChecking=no ${server_user}@${server_ip}

expect {
    "password:" {
        send "${server_pass}\r"
        exp_continue
    }
    "# " {
        # Make script executable
        send "chmod +x /tmp/deploy-on-server.sh\r"
        expect "# "
        
        # Run deployment
        send "/tmp/deploy-on-server.sh\r"
        
        # Wait for completion
        expect {
            "Deployment Complete!" {
                puts "\n✅ Deployment completed successfully!"
            }
            "ERROR" {
                puts "\n❌ Deployment error occurred"
            }
            timeout {
                puts "\n⚠️  Deployment may still be running..."
            }
        }
        
        expect "# "
        send "exit\r"
    }
    timeout {
        puts "Connection timeout"
        exit 1
    }
}
expect eof
DEPLOY_EOF

echo ""
echo "=========================================="
echo "Deployment Process Complete!"
echo "=========================================="
echo ""
echo "Test your proxy:"
echo "  curl -x http://${SERVER_IP}:8080 http://httpbin.org/ip"
echo ""
echo "Check status:"
echo "  ssh ${SERVER_USER}@${SERVER_IP} 'systemctl status game-proxy'"
echo ""
