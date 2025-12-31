output "server_ip" {
  description = "Proxy server IP address"
  value       = var.server_ip
}

output "server_ipv6" {
  description = "Proxy server IPv6 address"
  value       = var.server_ipv6
}

output "proxy_url" {
  description = "Proxy URL"
  value       = "http://${var.server_ip}:${var.proxy_port}"
}

output "proxy_url_ipv6" {
  description = "Proxy URL (IPv6)"
  value       = "http://[${var.server_ipv6}]:${var.proxy_port}"
}

output "connection_test" {
  description = "Test proxy connection"
  value       = "curl -x ${var.server_ip}:${var.proxy_port} http://httpbin.org/ip"
}

output "status_check" {
  description = "Check proxy status"
  value       = "ssh ${var.server_user}@${var.server_ip} 'systemctl status game-proxy'"
}

output "logs" {
  description = "View proxy logs"
  value       = "ssh ${var.server_user}@${var.server_ip} 'tail -f /opt/game-proxy/logs/proxy.log'"
}
