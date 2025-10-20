# Environment Setup Guide - Secrets & Variables

## 🎯 Overview

This guide shows you how to configure all environment-specific secrets and variables needed for the consolidated deployment workflow.

---

## 🏗️ Environment Structure

**GitHub Environments (4 total):**
- `dev-testnet` - Development environment
- `uat-testnet` - UAT environment (your focus)
- `sepolia-testnet` - Sepolia testnet (requires approval)
- `mainnet` - Production mainnet (requires approval)

---

## 🔑 ArgoCD Configuration

### ArgoCD Instances

| Environment | ArgoCD URL | Used By |
|---|---|---|
| **Dev & UAT** | https://argocd-uat.ten.xyz | dev-testnet, uat-testnet |
| **Sepolia** | https://argo-sepolia.ten.xyz | sepolia-testnet |
| **Mainnet** | https://argocd-mainnet.ten.xyz | mainnet |

### Getting ArgoCD Token

**Step 1: Create Service Account in ArgoCD Namespace**

For each ArgoCD instance, run:

```bash
# 1. SSH into the cluster
ssh <cluster-admin>

# 2. Create service account
kubectl create serviceaccount github-actions -n argocd

# 3. Create cluster role binding for permissions
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

# 4. Generate token (valid for 1 year)
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)
echo $TOKEN
```

**Step 2: Copy Token**
- Copy the output token (long alphanumeric string)
- Save it securely (you'll need it for GitHub)

**Step 3: Verify Token Works**
```bash
# Test the token
curl -H "Authorization: Bearer $TOKEN" \
  https://argocd-uat.ten.xyz/api/v1/applications

# Should return list of applications (not 401 error)
```

---

## 📝 GitHub Environment Secrets Setup

### For Each Environment in GitHub

Navigate to: **Settings → Environments → [environment-name] → Secrets**

#### 1. dev-testnet Environment

**Secrets to add:**

| Secret Name | Value | Source |
|---|---|---|
| `ARGOCD_SERVER_NONPROD` | `https://argocd-uat.ten.xyz` | Fixed URL |
| `ARGOCD_TOKEN` | `<token-from-step-above>` | Generated token |
| `DEPLOY_ACTIONS_PAT` | `<github-pat-with-repo-access>` | GitHub PAT |
| `REGISTRY_PASSWORD` | `<azure-registry-password>` | Azure ACR |
| `AZURE_CREDENTIALS` | `<azure-service-principal-json>` | Azure SP |
| `L1_HTTP_URL` | `https://eth-sepolia.drpc.org` | Sepolia RPC |
| `ACCOUNT_PK_WORKER` | `<deployer-private-key>` | Key vault |
| `ETHERSCAN_API_KEY` | `<etherscan-api-key>` | Etherscan |
| `L2_DEPLOYER_KEY` | `<l2-deployer-private-key>` | Key vault |
| `GH_TOKEN` | `<github-token-for-dispatch>` | GitHub token |

**Variables to add (if using):**

| Variable Name | Value |
|---|---|
| `ACCOUNT_ADDR_NODE_0` | Initial sequencer address |
| `CHAIN_ID` | 443 (or your chain ID) |
| `L1_CHALLENGE_PERIOD` | Challenge period value |
| `FAUCET_INITIAL_FUNDS` | Faucet funding amount |

#### 2. uat-testnet Environment

**Same as dev-testnet (identical secrets)**

```
ARGOCD_SERVER_NONPROD: https://argocd-uat.ten.xyz
ARGOCD_TOKEN: <token-from-uat-argocd>
DEPLOY_ACTIONS_PAT: <github-pat>
REGISTRY_PASSWORD: <azure-password>
AZURE_CREDENTIALS: <service-principal-json>
L1_HTTP_URL: <uat-l1-rpc>
ACCOUNT_PK_WORKER: <deployer-pk>
ETHERSCAN_API_KEY: <etherscan-key>
L2_DEPLOYER_KEY: <l2-deployer-pk>
GH_TOKEN: <github-token>
```

#### 3. sepolia-testnet Environment

**Secrets:**

| Secret Name | Value |
|---|---|
| `ARGOCD_SERVER_NONPROD` | `https://argo-sepolia.ten.xyz` |
| `ARGOCD_TOKEN` | `<token-from-sepolia-argocd>` |
| `DEPLOY_ACTIONS_PAT` | `<github-pat>` |
| `REGISTRY_PASSWORD` | `<azure-password>` |
| `AZURE_CREDENTIALS` | `<service-principal-json>` |
| `L1_HTTP_URL` | `<sepolia-l1-rpc>` |
| `ACCOUNT_PK_WORKER` | `<deployer-pk>` |
| `ETHERSCAN_API_KEY` | `<etherscan-api-key>` |
| `L2_DEPLOYER_KEY` | `<l2-deployer-pk>` |
| `GH_TOKEN` | `<github-token>` |

**Requires Approval Setup:**
```
Settings → Environments → sepolia-testnet → Required reviewers
Add team members who can approve production deployments
```

#### 4. mainnet Environment

**Secrets:**

| Secret Name | Value |
|---|---|
| `ARGOCD_SERVER_PROD` | `https://argocd-mainnet.ten.xyz` |
| `ARGOCD_TOKEN` | `<token-from-mainnet-argocd>` |
| `DEPLOY_ACTIONS_PAT` | `<github-pat>` |
| `REGISTRY_PASSWORD` | `<azure-password>` |
| `AZURE_CREDENTIALS` | `<service-principal-json>` |
| `L1_HTTP_URL` | `<mainnet-l1-rpc>` |
| `ACCOUNT_PK_WORKER` | `<deployer-pk>` |
| `ETHERSCAN_API_KEY` | `<etherscan-api-key>` |
| `L2_DEPLOYER_KEY` | `<l2-deployer-pk>` |
| `GH_TOKEN` | `<github-token>` |

**Requires Approval Setup:**
```
Settings → Environments → mainnet → Required reviewers
Add team members who can approve production deployments
```

---

## 📋 Complete Secret Reference

### ArgoCD Secrets

**ARGOCD_SERVER_NONPROD** (for dev & uat)
```
https://argocd-uat.ten.xyz
```

**ARGOCD_SERVER_PROD** (for sepolia & mainnet)
```
https://argocd-mainnet.ten.xyz  (Actually for mainnet)
https://argo-sepolia.ten.xyz    (For sepolia)
```

**ARGOCD_TOKEN**
```
How to get:
1. kubectl create serviceaccount github-actions -n argocd
2. kubectl create clusterrolebinding github-actions --clusterrole=cluster-admin --serviceaccount=argocd:github-actions
3. kubectl -n argocd create token github-actions --duration=8760h
4. Copy the output token

Format: Long alphanumeric string (JWT token)
Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### GitHub Secrets

**DEPLOY_ACTIONS_PAT**
```
How to get:
1. GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token"
3. Name: "Deploy Actions PAT"
4. Expiration: 90 days or custom
5. Select scopes: repo (full control of private repositories)
6. Generate and copy token

Format: ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Needs: Write access to ten-apps repo
```

**GH_TOKEN**
```
Same as DEPLOY_ACTIONS_PAT or separate token
Used for: Repository dispatches to ten-test
```

### Azure Secrets

**REGISTRY_PASSWORD**
```
How to get:
1. Azure Portal → Container Registries → testnetobscuronet → Access keys
2. Copy "password" or "password2"
3. Format: Long alphanumeric string
```

**AZURE_CREDENTIALS**
```
How to get:
1. Azure CLI: az ad sp create-for-rbac --role Contributor --scopes /subscriptions/<subscription-id>
2. Output is JSON with: clientId, clientSecret, subscriptionId, tenantId
3. Copy entire JSON

Format:
{
  "clientId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
  "subscriptionId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "tenantId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

### L1 & L2 Secrets

**L1_HTTP_URL** (RPC endpoint for L1)
```
Dev/UAT: https://eth-sepolia.drpc.org (or internal endpoint)
Sepolia: https://eth-sepolia.drpc.org
Mainnet: https://eth-mainnet.drpc.org (or internal endpoint)
```

**ACCOUNT_PK_WORKER** (Private key for L1 deployer)
```
Format: 0x followed by 64 hex characters
Example: 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
Source: Key vault or secure storage
```

**L2_DEPLOYER_KEY** (Private key for L2 deployer)
```
Format: 0x followed by 64 hex characters
Source: Key vault or secure storage
```

**ETHERSCAN_API_KEY** (For contract verification)
```
How to get:
1. Etherscan.io → Sign up or login
2. My Profile → API Keys
3. Create new API key
4. Copy the key

Format: Alphanumeric string
```

---

## ✅ Setup Checklist

### Step 1: Get ArgoCD Tokens

- [ ] Generate token for argocd-uat.ten.xyz
- [ ] Generate token for argo-sepolia.ten.xyz
- [ ] Generate token for argocd-mainnet.ten.xyz
- [ ] Test each token with curl command
- [ ] Save tokens securely

### Step 2: Get GitHub PATs

- [ ] Create DEPLOY_ACTIONS_PAT (with repo scope)
- [ ] Create GH_TOKEN (or reuse DEPLOY_ACTIONS_PAT)
- [ ] Verify PAT has write access to ten-apps

### Step 3: Get Azure Credentials

- [ ] Get REGISTRY_PASSWORD from Azure ACR
- [ ] Create AZURE_CREDENTIALS service principal
- [ ] Save JSON credentials securely

### Step 4: Get Deployer Keys

- [ ] Get ACCOUNT_PK_WORKER from key vault
- [ ] Get L2_DEPLOYER_KEY from key vault
- [ ] Verify keys are valid

### Step 5: Get API Keys

- [ ] Get ETHERSCAN_API_KEY from Etherscan
- [ ] Verify key is active

### Step 6: Configure GitHub Environments

For each environment (dev-testnet, uat-testnet, sepolia-testnet, mainnet):

- [ ] Go to Settings → Environments → [env-name]
- [ ] Add all secrets from Secret Reference above
- [ ] Add variables if using
- [ ] For prod environments (sepolia, mainnet):
  - [ ] Enable "Deployment branches and environments" rule
  - [ ] Add required reviewers
  - [ ] Save

### Step 7: Verify Setup

- [ ] Test dev-testnet environment (no approval)
- [ ] Run dummy deployment to verify secrets work
- [ ] Check ArgoCD token works: `argocd app list`
- [ ] Check GitHub PAT works: curl to GitHub API
- [ ] Check Azure credentials work: `az login --service-principal`

---

## 🔍 Verification Commands

### Test ArgoCD Token

```bash
# After setting ARGOCD_TOKEN and ARGOCD_SERVER
curl -H "Authorization: Bearer $ARGOCD_TOKEN" \
  https://argocd-uat.ten.xyz/api/v1/applications

# Should return JSON list of apps (not 401)
```

### Test GitHub PAT

```bash
curl -H "Authorization: token $GITHUB_PAT" \
  https://api.github.com/user

# Should return user info (not 401)
```

### Test Azure Credentials

```bash
# With AZURE_CREDENTIALS (service principal)
az login --service-principal \
  -u <clientId> \
  -p <clientSecret> \
  --tenant <tenantId>

# Should succeed
```

---

## 🌍 Environment-Specific Differences

### Dev-testnet
- ArgoCD: argocd-uat.ten.xyz
- No approval needed
- Safe for experimentation

### UAT-testnet
- ArgoCD: argocd-uat.ten.xyz (same as dev)
- No approval needed
- Your primary testing environment

### Sepolia-testnet
- ArgoCD: argo-sepolia.ten.xyz
- **Requires approval** (destructive deployments)
- Add required reviewers
- Production-like environment

### Mainnet
- ArgoCD: argocd-mainnet.ten.xyz
- **Requires approval** (all destructive deployments)
- Add required reviewers
- Production environment

---

## 📋 Minimal Setup for UAT Testing

If you only want to test UAT destructive deployments initially:

**Minimum secrets needed:**

1. `ARGOCD_SERVER_NONPROD`: https://argocd-uat.ten.xyz
2. `ARGOCD_TOKEN`: Token from argocd-uat
3. `DEPLOY_ACTIONS_PAT`: GitHub PAT
4. `REGISTRY_PASSWORD`: Azure ACR password
5. `AZURE_CREDENTIALS`: Azure service principal JSON
6. `L1_HTTP_URL`: L1 RPC endpoint
7. `ACCOUNT_PK_WORKER`: Deployer private key
8. `ETHERSCAN_API_KEY`: Etherscan API key
9. `L2_DEPLOYER_KEY`: L2 deployer private key
10. `GH_TOKEN`: GitHub token

**Can add later:**
- Variables (ACCOUNT_ADDR_NODE_0, CHAIN_ID, etc.)
- Sepolia & Mainnet environments

---

## 🔐 Security Notes

**Secrets Best Practices:**

✅ **Do:**
- Rotate tokens every 1 year
- Use service accounts (not personal accounts)
- Use least-privilege roles
- Store in secure vaults
- Audit secret access logs

❌ **Don't:**
- Commit secrets to git
- Share secrets in Slack/email
- Use personal passwords
- Reuse secrets across environments
- Keep secrets indefinitely

**Rotation Schedule:**
- ArgoCD tokens: Annually (from --duration=8760h)
- GitHub PATs: Every 90 days
- Service principals: As per corporate policy
- Private keys: Never share, rotate if compromised

---

## 🆘 Troubleshooting Setup

### Issue: "401 Unauthorized" from ArgoCD

**Causes:**
- Invalid or expired token
- Wrong ArgoCD URL
- Token not created properly
- Service account missing permissions

**Solution:**
```bash
# Regenerate token
kubectl -n argocd create token github-actions --duration=8760h

# Verify service account exists
kubectl get serviceaccount -n argocd | grep github

# Verify role binding exists
kubectl get clusterrolebinding | grep github

# Test with curl
curl -H "Authorization: Bearer $TOKEN" \
  https://argocd-uat.ten.xyz/api/v1/applications
```

### Issue: "403 Forbidden" from GitHub

**Causes:**
- PAT doesn't have repo scope
- PAT expired
- PAT doesn't have write access

**Solution:**
```bash
# Create new PAT
# GitHub → Settings → Developer settings → Personal access tokens
# Select: repo (full control)
# Generate and copy
```

### Issue: "Invalid credentials" from Azure

**Causes:**
- AZURE_CREDENTIALS JSON malformed
- Service principal deleted
- Wrong subscription ID

**Solution:**
```bash
# Recreate service principal
az ad sp create-for-rbac --role Contributor

# Copy output JSON exactly
```

---

## 📚 Next Steps

1. **Gather all secrets** using commands above
2. **Create service account** in each ArgoCD
3. **Generate tokens** for each ArgoCD instance
4. **Configure GitHub Environments** with secrets
5. **Test UAT setup** with dry-run deployment
6. **Add approval reviewers** for prod environments
7. **Document in your wiki** for team reference

---

## 🎯 Quick Summary

| Item | UAT Value | Sepolia Value | Mainnet Value |
|---|---|---|---|
| ArgoCD URL | https://argocd-uat.ten.xyz | https://argo-sepolia.ten.xyz | https://argocd-mainnet.ten.xyz |
| Approval? | ❌ No | ✅ Yes (destructive) | ✅ Yes (destructive) |
| Token source | argocd-uat cluster | argo-sepolia cluster | argocd-mainnet cluster |
| Use case | Your testing | Pre-production | Production |

**All other secrets (DEPLOY_ACTIONS_PAT, REGISTRY_PASSWORD, etc.) are shared across all environments.**

---

**Ready to set up? Start with ArgoCD tokens! 🚀**
