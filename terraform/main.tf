terraform {
  required_version = ">= 1.0"
  
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

# Archive source files for deployment
data "archive_file" "proxy_source" {
  type        = "tar"
  source_dir  = "${path.module}/../"
  output_path = "${path.module}/proxy-deploy.tar.gz"
  
  excludes = [
    ".git",
    ".gitignore",
    "terraform",
    "*.md",
    "docs",
    "docker",
    "gui",
    "*.sh",
    "*.log",
    "*.tmp",
    "proxy-deploy.tar.gz",
    ".DS_Store",
    "*.swp",
    "*~"
  ]
}

# Deploy proxy to server
resource "null_resource" "deploy_proxy" {
  triggers = {
    source_hash = data.archive_file.proxy_source.output_md5
    server_ip   = var.server_ip
  }

  connection {
    type     = "ssh"
    host     = var.server_ip
    user     = var.server_user
    password = var.server_password
    timeout  = "5m"
  }

  # Upload source files
  provisioner "file" {
    source      = data.archive_file.proxy_source.output_path
    destination = "/tmp/proxy-deploy.tar.gz"
  }

  # Upload deployment script
  provisioner "file" {
    content     = file("${path.module}/scripts/deploy.sh")
    destination = "/tmp/deploy-proxy.sh"
  }

  # Execute deployment
  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/deploy-proxy.sh",
      "bash /tmp/deploy-proxy.sh"
    ]
  }
}

# Configure firewall rules
resource "null_resource" "configure_firewall" {
  depends_on = [null_resource.deploy_proxy]

  connection {
    type     = "ssh"
    host     = var.server_ip
    user     = var.server_user
    password = var.server_password
    timeout  = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "ufw allow ${var.proxy_port}/tcp || true",
      "ufw allow 1080/tcp || true",
      "ufw allow 5353/udp || true",
      "ufw allow 22/tcp || true",
      "ufw --force enable || true"
    ]
  }
}

# Health check
resource "null_resource" "health_check" {
  depends_on = [null_resource.deploy_proxy, null_resource.configure_firewall]

  connection {
    type     = "ssh"
    host     = var.server_ip
    user     = var.server_user
    password = var.server_password
    timeout  = "1m"
  }

  provisioner "remote-exec" {
    inline = [
      "sleep 5",
      "systemctl is-active game-proxy || exit 1",
      "ss -tulpn | grep ${var.proxy_port} || exit 1",
      "echo 'Proxy is running on port ${var.proxy_port}'"
    ]
  }
}
