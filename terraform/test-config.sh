#!/bin/bash
# Test Terraform configuration without actually deploying

set -e

echo "=========================================="
echo "Testing Terraform Configuration"
echo "=========================================="

cd "$(dirname "$0")"

# Check files exist
echo "1. Checking required files..."
files=("main.tf" "variables.tf" "outputs.tf" "providers.tf" "scripts/deploy.sh")
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✅ $file exists"
    else
        echo "   ❌ $file missing"
        exit 1
    fi
done

# Check deploy script syntax
echo ""
echo "2. Checking deploy script syntax..."
if bash -n scripts/deploy.sh; then
    echo "   ✅ Deploy script syntax is valid"
else
    echo "   ❌ Deploy script has syntax errors"
    exit 1
fi

# Check Terraform files structure
echo ""
echo "3. Checking Terraform file structure..."
if grep -q "resource \"null_resource\"" main.tf; then
    echo "   ✅ Resources defined"
else
    echo "   ⚠️  No resources found in main.tf"
fi

if grep -q "variable" variables.tf; then
    echo "   ✅ Variables defined"
else
    echo "   ⚠️  No variables found"
fi

if grep -q "output" outputs.tf; then
    echo "   ✅ Outputs defined"
else
    echo "   ⚠️  No outputs found"
fi

# Check if terraform.tfvars exists
echo ""
echo "4. Checking configuration..."
if [ -f "terraform.tfvars" ]; then
    echo "   ✅ terraform.tfvars exists"
    if grep -q "server_ip" terraform.tfvars; then
        echo "   ✅ Server IP configured"
    fi
else
    echo "   ⚠️  terraform.tfvars not found (use terraform.tfvars.example)"
fi

# Check if Terraform is installed
echo ""
echo "5. Checking Terraform installation..."
if command -v terraform > /dev/null 2>&1; then
    terraform_version=$(terraform version | head -1)
    echo "   ✅ Terraform installed: $terraform_version"
    
    echo ""
    echo "6. Running terraform init (dry-run check)..."
    if terraform init -backend=false > /dev/null 2>&1; then
        echo "   ✅ Terraform init successful"
    else
        echo "   ⚠️  Terraform init had issues (may need providers)"
    fi
    
    echo ""
    echo "7. Running terraform validate..."
    if terraform validate > /dev/null 2>&1; then
        echo "   ✅ Terraform configuration is valid"
    else
        echo "   ⚠️  Terraform validation issues (check above)"
        terraform validate
    fi
else
    echo "   ⚠️  Terraform not installed"
    echo "   Install: https://developer.hashicorp.com/terraform/downloads"
fi

echo ""
echo "=========================================="
echo "Test Complete!"
echo "=========================================="
echo ""
echo "To deploy:"
echo "  1. Install Terraform (if not installed)"
echo "  2. cd terraform"
echo "  3. terraform init"
echo "  4. terraform plan"
echo "  5. terraform apply"
echo ""
