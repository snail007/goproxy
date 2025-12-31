variable "server_ip" {
  description = "Server IP address"
  type        = string
  default     = "194.233.73.36"
}

variable "server_user" {
  description = "SSH username"
  type        = string
  default     = "root"
}

variable "server_password" {
  description = "SSH password"
  type        = string
  sensitive   = true
  default     = "JVxgREA3hbe6"
}

variable "server_ipv6" {
  description = "Server IPv6 address"
  type        = string
  default     = "2407:3640:2299:7804::1"
}

variable "proxy_port" {
  description = "Proxy server port"
  type        = number
  default     = 8080
}

variable "proxy_mode" {
  description = "Proxy mode (game-http, http, socks5)"
  type        = string
  default     = "game-http"
}

variable "enable_stealth" {
  description = "Enable stealth features"
  type        = bool
  default     = true
}

variable "enable_game_optimizations" {
  description = "Enable game-specific optimizations"
  type        = bool
  default     = true
}
