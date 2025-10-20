# ArgoCD Setup Scripts

Automated bash scripts to generate ArgoCD tokens for each environment.

## 📋 Available Scripts

| Script | ArgoCD Instance | GitHub Environments | Usage |
|--------|---|---|---|
| `setup-argocd-uat.sh` | argocd-uat.ten.xyz | dev-testnet, uat-testnet | Your primary setup |
| `setup-argocd-sepolia.sh` | argo-sepolia.ten.xyz | sepolia-testnet | When ready for production testing |
| `setup-argocd-mainnet.sh` | argocd-mainnet.ten.xyz | mainnet | Production deployment |
| `setup-all-argocd.sh` | All three | All | Run all setups with menu |

## 🚀 Quick Start

### Option 1: Setup UAT Only (Recommended First)

```bash
cd .github/workflows
./setup-argocd-uat.sh
```

**What it does:**
1. Creates service account in ArgoCD namespace
2. Grants cluster admin permissions
3. Generates 1-year token
4. Verifies token works
5. Displays token for copying

**Output:**
- Token to copy and paste
- Instructions for adding to GitHub

### Option 2: Setup Each Environment Individually

```bash
# Setup Sepolia
./setup-argocd-sepolia.sh

# Setup Mainnet
./setup-argocd-mainnet.sh
```

### Option 3: Interactive Setup for All

```bash
./setup-all-argocd.sh
```

**Menu options:**
- 1: UAT only
- 2: Sepolia only
- 3: Mainnet only
- 4: All three

---

## 📝 Step-by-Step Usage

### Step 1: Run the Script

```bash
# SSH into the appropriate cluster first (if running from outside)
ssh <user>@<cluster>

# Then run the script
./setup-argocd-uat.sh
```

### Step 2: Wait for Completion

Script will:
- ✅ Create service account (if not exists)
- ✅ Create role binding (if not exists)
- ✅ Generate token
- ✅ Verify token works with curl

### Step 3: Copy the Token

The script will display a long token:
```
================================================
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3M...
================================================
```

Copy this entire token.

### Step 4: Add to GitHub

1. Go to your GitHub repository
2. Settings → Environments
3. For dev-testnet:
   - Click "dev-testnet"
   - Secrets → New repository secret
   - Name: `ARGOCD_TOKEN`
   - Value: (paste the token)
   - Click "Add secret"

4. Repeat for uat-testnet using same token

---

## 🔍 What Each Script Does

### setup-argocd-uat.sh

**Cluster:** argocd-uat.ten.xyz
**GitHub Environments:** dev-testnet, uat-testnet (both use same token)
**When to run:** First setup
**Approval:** Not required

```bash
./setup-argocd-uat.sh
```

### setup-argocd-sepolia.sh

**Cluster:** argo-sepolia.ten.xyz
**GitHub Environments:** sepolia-testnet
**When to run:** When ready for pre-production testing
**Approval:** Required (destructive only)

```bash
./setup-argocd-sepolia.sh
```

### setup-argocd-mainnet.sh

**Cluster:** argocd-mainnet.ten.xyz
**GitHub Environments:** mainnet
**When to run:** Before first mainnet deployment
**Approval:** Required (all destructive)

```bash
./setup-argocd-mainnet.sh
```

### setup-all-argocd.sh

**Purpose:** Interactive menu to run any combination
**Usage:** When setting up multiple environments at once

```bash
./setup-all-argocd.sh
# Choose: 1 (UAT), 2 (Sepolia), 3 (Mainnet), or 4 (All)
```

---

## ✅ Script Output Example

```
================================================
ArgoCD UAT Setup (argocd-uat.ten.xyz)
================================================

This script will:
1. Create service account in ArgoCD namespace
2. Grant permissions
3. Generate token (valid 1 year)
4. Verify token works

[Step 1/5] Creating service account...
✓ Service account created

[Step 2/5] Creating cluster role binding...
✓ Cluster role binding created

[Step 3/5] Generating token (valid 1 year)...
✓ Token generated

[Step 4/5] Token ready for GitHub:
================================================
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJnaXRodWItYWN0aW9ucyJ9.ABC...
================================================

⚠️  COPY THIS TOKEN AND SAVE IT SECURELY
You will need to add this to GitHub Secrets:
  Environment: dev-testnet
  Secret name: ARGOCD_TOKEN
  Secret value: (paste token above)

REPEAT for uat-testnet environment too!

[Step 5/5] Verifying token works...
✓ Token verification PASSED

✅ Setup complete! Token is valid and working.

📝 Next steps:
1. Go to GitHub → Settings → Environments → dev-testnet → Secrets
2. Add secret: ARGOCD_TOKEN = (paste token above)
3. Repeat for uat-testnet environment
```

---

## 🆘 Troubleshooting

### Script Fails: "command not found: kubectl"

**Problem:** kubectl is not in PATH

**Solution:**
```bash
# Add kubectl to PATH
export PATH=$PATH:$(which kubectl | xargs dirname)

# Or use full path
/usr/bin/kubectl -n argocd create token ...

# Or run from cluster with kubectl installed
ssh <user>@<cluster>
./setup-argocd-uat.sh
```

### Script Fails: "service account already exists"

**This is OK!** Script checks and reuses existing service account.

### Token Verification Fails: "401 Unauthorized"

**Problem:** Token or permissions not set up correctly

**Solution:**
```bash
# Verify service account permissions
kubectl get clusterrolebinding github-actions

# If missing, manually create:
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

# Generate token again
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)
echo $TOKEN
```

### Can't SSH into Cluster

**Problem:** SSH access not configured

**Solution:**
```bash
# Check SSH key
ssh-add ~/.ssh/id_rsa

# Try with verbose output
ssh -vvv <user>@<cluster>

# Ask team for SSH access instructions
```

---

## 🔐 Security Notes

### Token Security

- ✅ Valid for 1 year
- ✅ Can be rotated annually
- ✅ Service account, not personal credentials
- ✅ Limited to automation use (GitHub Actions)

### After Running Script

- ☐ Save token securely (1Password, LastPass, etc.)
- ☐ Don't share token in Slack/email
- ☐ Don't commit token to git
- ☐ Rotate annually

---

## 📋 Verification Checklist

After running each script:

- [ ] Script ran without errors
- [ ] Token displayed at end
- [ ] Token verification PASSED
- [ ] Token copied to secure location
- [ ] Added to GitHub Secrets for correct environment
- [ ] Verified in GitHub: Settings → Environments → [env] → Secrets

---

## 🎯 Recommended Setup Order

**Week 1:**
```bash
./setup-argocd-uat.sh
# Add token to dev-testnet & uat-testnet in GitHub
# Test deployment workflow
```

**Week 2 (After UAT Testing):**
```bash
./setup-argocd-sepolia.sh
# Add token to sepolia-testnet in GitHub
# Add required reviewers
# Test sepolia deployment
```

**Week 3 (After Sepolia Testing):**
```bash
./setup-argocd-mainnet.sh
# Add token to mainnet in GitHub
# Add required reviewers
# Prepare for production deployment
```

---

## 💡 Tips & Tricks

### Check if Service Account Already Exists

```bash
kubectl get serviceaccount -n argocd | grep github-actions
```

### Check if Role Binding Exists

```bash
kubectl get clusterrolebinding | grep github-actions
```

### Manually Generate Token (if script fails)

```bash
kubectl -n argocd create token github-actions --duration=8760h
```

### Test Token Manually

```bash
TOKEN="<paste-token-here>"
curl -H "Authorization: Bearer $TOKEN" \
  https://argocd-uat.ten.xyz/api/v1/applications
# Should return JSON list (not 401)
```

---

## 📞 Common Questions

**Q: Do I need to run script for both dev-testnet and uat-testnet?**
A: No! Run once for uat-testnet and use same token for both dev-testnet and uat-testnet.

**Q: How often do I need to run the script?**
A: Once per ArgoCD instance. Token is valid 1 year.

**Q: Can I rotate the token?**
A: Yes! Just run the script again (it generates new token) and update GitHub secret.

**Q: What if I lost the token?**
A: Run script again to generate new token.

**Q: Can I use for multiple repos?**
A: Yes! But each repo needs its own GitHub secret added.

---

**Ready to set up? Run the script! 🚀**

```bash
cd .github/workflows
./setup-argocd-uat.sh
```
