#!/bin/bash

# Server-side configuration script
# Run this on the server after deployment

set -e

PROXY_PORT="8080"
SOCKS5_PORT="1080"
UDP_PORT="5353"

echo "Configuring Game Proxy Server..."

# Enable IP forwarding
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
echo "net.ipv6.conf.all.forwarding=1" >> /etc/sysctl.conf
sysctl -p

# Configure iptables for NAT (if needed)
iptables -t nat -A PREROUTING -p tcp --dport $PROXY_PORT -j DNAT --to-destination 127.0.0.1:$PROXY_PORT
iptables -t nat -A POSTROUTING -j MASQUERADE

# IPv6 configuration
ip6tables -t nat -A PREROUTING -p tcp --dport $PROXY_PORT -j DNAT --to-destination [::1]:$PROXY_PORT

# Create log rotation
cat > /etc/logrotate.d/game-proxy <<EOF
/opt/game-proxy/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
}
EOF

# Set up monitoring script
cat > /opt/game-proxy/monitor.sh <<'EOF'
#!/bin/bash
# Monitor proxy health

PROXY_PORT=8080
LOG_FILE=/opt/game-proxy/logs/health.log

check_proxy() {
    if curl -s --connect-timeout 5 http://127.0.0.1:$PROXY_PORT > /dev/null 2>&1; then
        echo "$(date): Proxy is healthy" >> $LOG_FILE
        return 0
    else
        echo "$(date): Proxy is down, restarting..." >> $LOG_FILE
        systemctl restart game-proxy
        return 1
    fi
}

check_proxy
EOF

chmod +x /opt/game-proxy/monitor.sh

# Add to crontab for health checks
(crontab -l 2>/dev/null; echo "*/5 * * * * /opt/game-proxy/monitor.sh") | crontab -

echo "Server configuration complete!"
