#!/bin/bash
# Setup ArgoCD token for Mainnet (argocd-mainnet.ten.xyz)
# Usage: ./setup-argocd-mainnet.sh

set -e

echo "================================================"
echo "ArgoCD Mainnet Setup (argocd-mainnet.ten.xyz)"
echo "================================================"
echo ""
echo "This script will:"
echo "1. Create service account in ArgoCD namespace"
echo "2. Grant permissions"
echo "3. Generate token (valid 1 year)"
echo "4. Verify token works"
echo ""

# Step 1: Create service account
echo "[Step 1/5] Creating service account..."
if kubectl get serviceaccount -n argocd github-actions &>/dev/null; then
    echo "✓ Service account already exists"
else
    kubectl create serviceaccount github-actions -n argocd
    echo "✓ Service account created"
fi

# Step 2: Create cluster role binding
echo ""
echo "[Step 2/5] Creating cluster role binding..."
if kubectl get clusterrolebinding github-actions &>/dev/null; then
    echo "✓ Cluster role binding already exists"
else
    kubectl create clusterrolebinding github-actions \
        --clusterrole=cluster-admin \
        --serviceaccount=argocd:github-actions
    echo "✓ Cluster role binding created"
fi

# Step 3: Generate token
echo ""
echo "[Step 3/5] Generating token (valid 1 year)..."
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)
echo "✓ Token generated"

# Step 4: Display token
echo ""
echo "[Step 4/5] Token ready for GitHub:"
echo "================================================"
echo "$TOKEN"
echo "================================================"
echo ""
echo "⚠️  COPY THIS TOKEN AND SAVE IT SECURELY"
echo "You will need to add this to GitHub Secrets:"
echo "  Environment: mainnet"
echo "  Secret name: ARGOCD_TOKEN"
echo "  Secret value: (paste token above)"
echo ""

# Step 5: Verify token works
echo "[Step 5/5] Verifying token works..."
ARGOCD_URL="https://argocd-mainnet.ten.xyz"
if curl -s -H "Authorization: Bearer $TOKEN" \
    "$ARGOCD_URL/api/v1/applications" > /dev/null 2>&1; then
    echo "✓ Token verification PASSED"
    echo ""
    echo "✅ Setup complete! Token is valid and working."
else
    echo "❌ Token verification FAILED"
    echo "   This might mean:"
    echo "   - ArgoCD server is not accessible"
    echo "   - Token is invalid"
    echo "   - Permissions not set up correctly"
    exit 1
fi

echo ""
echo "📝 Next steps:"
echo "1. Go to GitHub → Settings → Environments → mainnet → Secrets"
echo "2. Add secret: ARGOCD_TOKEN = (paste token above)"
echo "3. Add required reviewers to mainnet environment"
echo ""
echo "⚠️  PRODUCTION ENVIRONMENT - DOUBLE CHECK BEFORE DEPLOYING!"
echo ""
echo "================================================"
