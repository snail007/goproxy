# Terraform Deployment for Game Proxy

This directory contains Terraform configuration to deploy the Game Proxy server infrastructure.

## Prerequisites

1. **Terraform installed** (>= 1.0)
   ```bash
   # Install Terraform
   curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
   sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
   sudo apt-get update && sudo apt-get install terraform
   ```

2. **SSH access** to the target server

## Quick Start

1. **Configure variables**:
   ```bash
   cd terraform
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your server details
   ```

2. **Initialize Terraform**:
   ```bash
   terraform init
   ```

3. **Plan deployment**:
   ```bash
   terraform plan
   ```

4. **Apply deployment**:
   ```bash
   terraform apply
   ```

5. **Verify deployment**:
   ```bash
   terraform output
   ```

## Configuration

### Variables

Edit `terraform.tfvars` to configure:

- `server_ip`: Server IPv4 address
- `server_user`: SSH username (default: root)
- `server_password`: SSH password
- `server_ipv6`: Server IPv6 address
- `proxy_port`: Proxy port (default: 8080)
- `proxy_mode`: Proxy mode (default: game-http)
- `enable_stealth`: Enable stealth features (default: true)
- `enable_game_optimizations`: Enable game optimizations (default: true)

### Outputs

After deployment, Terraform outputs:

- `proxy_url`: HTTP proxy URL
- `proxy_url_ipv6`: IPv6 proxy URL
- `connection_test`: Command to test proxy
- `status_check`: Command to check service status
- `logs`: Command to view logs

## Usage Examples

### Test Proxy Connection

```bash
# Using Terraform output
terraform output -raw connection_test | bash

# Or manually
curl -x http://194.233.73.36:8080 http://httpbin.org/ip
```

### Check Service Status

```bash
terraform output -raw status_check | bash
```

### View Logs

```bash
terraform output -raw logs | bash
```

## Management Commands

### Update Deployment

After making code changes:

```bash
terraform apply
```

### Destroy Deployment

To remove the proxy from the server:

```bash
terraform destroy
```

**Note**: This will stop and disable the service but won't remove installed files.

### Refresh State

If you manually changed something on the server:

```bash
terraform refresh
```

## Troubleshooting

### Connection Issues

If SSH connection fails:

1. Verify server IP and credentials in `terraform.tfvars`
2. Check firewall allows SSH (port 22)
3. Test SSH manually: `ssh root@194.233.73.36`

### Build Failures

If Go build fails:

1. Check logs: `ssh root@194.233.73.36 'tail -50 /opt/game-proxy/logs/proxy-error.log'`
2. Verify Go installation: `ssh root@194.233.73.36 'go version'`
3. Check source files: `ssh root@194.233.73.36 'ls -la /opt/game-proxy/'`

### Service Not Starting

```bash
# Check service status
ssh root@194.233.73.36 'systemctl status game-proxy'

# Check logs
ssh root@194.233.73.36 'journalctl -u game-proxy -n 50'

# Restart service
ssh root@194.233.73.36 'systemctl restart game-proxy'
```

## File Structure

```
terraform/
├── main.tf              # Main Terraform configuration
├── variables.tf          # Variable definitions
├── outputs.tf           # Output values
├── terraform.tfvars     # Variable values (not in git)
├── terraform.tfvars.example  # Example variables
├── scripts/
│   └── deploy.sh        # Deployment script (executed on server)
└── README.md            # This file
```

## Security Notes

- **Never commit `terraform.tfvars`** with real passwords to git
- Use environment variables or Terraform Cloud for sensitive values
- Consider using SSH keys instead of passwords
- Use `terraform.tfvars.example` as a template

## Advanced Usage

### Using Environment Variables

Instead of `terraform.tfvars`, use environment variables:

```bash
export TF_VAR_server_password="your-password"
terraform apply
```

### Using Terraform Cloud

1. Create a workspace in Terraform Cloud
2. Set variables in the workspace
3. Mark sensitive variables appropriately
4. Run `terraform init` with backend configuration

### Custom Provider Configuration

Edit `providers.tf` to use custom provider settings:

```hcl
provider "null" {
  # Custom configuration
}
```
