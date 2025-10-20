#!/bin/bash
# Master setup script for all ArgoCD environments
# Usage: ./setup-all-argocd.sh

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         ArgoCD Token Setup - All Environments                 ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Function to run individual setup scripts
run_setup() {
    local env_name=$1
    local script_name=$2
    local cluster_desc=$3

    echo ""
    echo "════════════════════════════════════════════════════════════════"
    echo "Setup: $env_name"
    echo "Cluster: $cluster_desc"
    echo "════════════════════════════════════════════════════════════════"
    echo ""

    if [ -f "$script_name" ]; then
        bash "$script_name"
    else
        echo "❌ Script not found: $script_name"
        echo "   Expected location: $(pwd)/$script_name"
        return 1
    fi
}

# Determine which setups to run
echo "Available setups:"
echo "  1. UAT (argocd-uat.ten.xyz) - for dev-testnet & uat-testnet"
echo "  2. Sepolia (argo-sepolia.ten.xyz) - for sepolia-testnet"
echo "  3. Mainnet (argocd-mainnet.ten.xyz) - for mainnet"
echo "  4. All (1, 2, 3)"
echo ""
read -p "Select setup (1/2/3/4) [default: 1]: " choice
choice=${choice:-1}

case $choice in
    1)
        run_setup "UAT" "setup-argocd-uat.sh" "argocd-uat.ten.xyz (shared by dev & uat)"
        ;;
    2)
        run_setup "Sepolia" "setup-argocd-sepolia.sh" "argo-sepolia.ten.xyz"
        ;;
    3)
        run_setup "Mainnet" "setup-argocd-mainnet.sh" "argocd-mainnet.ten.xyz"
        ;;
    4)
        run_setup "UAT" "setup-argocd-uat.sh" "argocd-uat.ten.xyz (shared by dev & uat)"
        run_setup "Sepolia" "setup-argocd-sepolia.sh" "argo-sepolia.ten.xyz"
        run_setup "Mainnet" "setup-argocd-mainnet.sh" "argocd-mainnet.ten.xyz"
        ;;
    *)
        echo "❌ Invalid choice: $choice"
        exit 1
        ;;
esac

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "✅ Setup Complete!"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "📝 Summary of what to do next:"
echo ""
echo "1. Copy the token(s) from above"
echo "2. Go to GitHub Repository Settings"
echo "3. For each environment:"
echo "   → Settings → Environments → [environment-name] → Secrets"
echo "   → New repository secret"
echo "   → Name: ARGOCD_TOKEN"
echo "   → Value: (paste token)"
echo ""
echo "4. For production environments (sepolia-testnet, mainnet):"
echo "   → Add required reviewers to the environment"
echo ""
