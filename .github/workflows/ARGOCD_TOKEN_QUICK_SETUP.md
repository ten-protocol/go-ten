# ArgoCD Token Generation - Quick Setup Guide

## 🎯 Quick Summary

| Environment | ArgoCD URL | GitHub Secret | Commands |
|---|---|---|---|
| **Dev & UAT** | https://argocd-uat.ten.xyz | `ARGOCD_TOKEN` | See "UAT Setup" |
| **Sepolia** | https://argo-sepolia.ten.xyz | `ARGOCD_TOKEN` | See "Sepolia Setup" |
| **Mainnet** | https://argocd-mainnet.ten.xyz | `ARGOCD_TOKEN` | See "Mainnet Setup" |

---

## 🔑 UAT Setup (Dev & UAT both use same ArgoCD)

### Step 1: Access ArgoCD Cluster

```bash
# SSH into cluster with argocd-uat
ssh <cluster-admin>@<uat-cluster>

# Verify kubectl access
kubectl get nodes
```

### Step 2: Create Service Account

```bash
# Create service account in argocd namespace
kubectl create serviceaccount github-actions -n argocd

# Verify it was created
kubectl get serviceaccount -n argocd | grep github-actions
```

### Step 3: Grant Permissions

```bash
# Create cluster role binding (full admin access for deployments)
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

# Verify
kubectl get clusterrolebinding | grep github-actions
```

### Step 4: Generate Token

```bash
# Generate token (valid for 1 year)
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)

# Print and copy
echo $TOKEN
```

**Output:** Long JWT token starting with `eyJhbGc...`

### Step 5: Verify Token Works

```bash
# Test the token
curl -s -H "Authorization: Bearer $TOKEN" \
  https://argocd-uat.ten.xyz/api/v1/applications | jq '.'

# Should return JSON list of apps (no 401 error)
```

### Step 6: Save to GitHub

**For both dev-testnet and uat-testnet environments:**

Go to:
```
GitHub Repo → Settings → Environments
→ dev-testnet → Secrets → New repository secret
```

**Add secret:**
- Name: `ARGOCD_TOKEN`
- Value: (paste the token from step 4)

**Repeat for uat-testnet environment**

---

## 🔑 Sepolia Setup

### Step 1: Access Sepolia ArgoCD Cluster

```bash
# SSH into cluster with argo-sepolia
ssh <cluster-admin>@<sepolia-cluster>

# Verify kubectl access
kubectl get nodes
```

### Step 2: Create Service Account

```bash
kubectl create serviceaccount github-actions -n argocd
kubectl get serviceaccount -n argocd | grep github-actions
```

### Step 3: Grant Permissions

```bash
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

kubectl get clusterrolebinding | grep github-actions
```

### Step 4: Generate Token

```bash
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)
echo $TOKEN
```

### Step 5: Verify Token Works

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  https://argo-sepolia.ten.xyz/api/v1/applications | jq '.'
```

### Step 6: Save to GitHub

Go to:
```
GitHub Repo → Settings → Environments
→ sepolia-testnet → Secrets → New repository secret
```

**Add secret:**
- Name: `ARGOCD_TOKEN`
- Value: (paste the token)

---

## 🔑 Mainnet Setup

### Step 1: Access Mainnet ArgoCD Cluster

```bash
ssh <cluster-admin>@<mainnet-cluster>
kubectl get nodes
```

### Step 2: Create Service Account

```bash
kubectl create serviceaccount github-actions -n argocd
kubectl get serviceaccount -n argocd | grep github-actions
```

### Step 3: Grant Permissions

```bash
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

kubectl get clusterrolebinding | grep github-actions
```

### Step 4: Generate Token

```bash
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)
echo $TOKEN
```

### Step 5: Verify Token Works

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  https://argocd-mainnet.ten.xyz/api/v1/applications | jq '.'
```

### Step 6: Save to GitHub

Go to:
```
GitHub Repo → Settings → Environments
→ mainnet → Secrets → New repository secret
```

**Add secret:**
- Name: `ARGOCD_TOKEN`
- Value: (paste the token)

---

## ✅ Verification Checklist

### After Setting Up Each Token

- [ ] Token generated successfully
- [ ] Test curl returns apps (not 401)
- [ ] Secret added to GitHub Environment
- [ ] Workflow can reference it

### Verify in Workflow

Run a test deployment:
```
GitHub → Actions → manual-deploy-k8s-consolidated
→ Run workflow
→ Select testnet_type: dev-testnet
→ Check workflow logs for: "✅ Syncing ArgoCD applications"
```

If logs show `401 Unauthorized`, token is invalid.

---

## 🔐 Token Security

### Token Properties

- **Validity:** 1 year (from --duration=8760h)
- **Scope:** Cluster admin (can access all resources)
- **Rotation:** Regenerate annually

### Rotation Process

When token expires (or every 1 year):

```bash
# Step 1: Generate new token
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)

# Step 2: Update GitHub secret
# GitHub → Environments → [env] → Secrets → ARGOCD_TOKEN → Update value

# Step 3: Test
curl -H "Authorization: Bearer $TOKEN" https://argocd-uat.ten.xyz/api/v1/applications
```

### What NOT to Do

❌ Don't share tokens in Slack/email
❌ Don't commit tokens to git
❌ Don't use personal account tokens
❌ Don't use expired tokens

---

## 🆘 Troubleshooting

### Issue: "serviceaccount github-actions created" but can't find it

```bash
# Check if it exists
kubectl get serviceaccount -n argocd

# If not there, create it again
kubectl create serviceaccount github-actions -n argocd
```

### Issue: "Token create" fails

```bash
# Make sure you're in the right namespace
kubectl get namespace | grep argocd

# Try creating again with explicit namespace
kubectl -n argocd create token github-actions --duration=8760h
```

### Issue: curl returns 401 Unauthorized

**Possible causes:**
- Token is invalid or expired
- Wrong ArgoCD URL
- Service account doesn't have permissions

**Solutions:**
```bash
# 1. Verify role binding exists
kubectl get clusterrolebinding | grep github-actions

# 2. Re-create if missing
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

# 3. Generate new token
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)

# 4. Test again
curl -H "Authorization: Bearer $TOKEN" \
  https://argocd-uat.ten.xyz/api/v1/applications
```

### Issue: Can't SSH into cluster

```bash
# Check SSH key
ssh-add ~/.ssh/id_rsa

# Try with verbose output
ssh -vvv <user>@<host>

# Ask team for correct access credentials
```

---

## 📋 Commands Cheat Sheet

### UAT (argocd-uat.ten.xyz)

```bash
# SSH & setup
ssh <user>@<uat-cluster>
kubectl create serviceaccount github-actions -n argocd
kubectl create clusterrolebinding github-actions --clusterrole=cluster-admin --serviceaccount=argocd:github-actions

# Generate token
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h) && echo $TOKEN

# Verify
curl -H "Authorization: Bearer $TOKEN" https://argocd-uat.ten.xyz/api/v1/applications
```

### Sepolia (argo-sepolia.ten.xyz)

```bash
# SSH & setup
ssh <user>@<sepolia-cluster>
kubectl create serviceaccount github-actions -n argocd
kubectl create clusterrolebinding github-actions --clusterrole=cluster-admin --serviceaccount=argocd:github-actions

# Generate token
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h) && echo $TOKEN

# Verify
curl -H "Authorization: Bearer $TOKEN" https://argo-sepolia.ten.xyz/api/v1/applications
```

### Mainnet (argocd-mainnet.ten.xyz)

```bash
# SSH & setup
ssh <user>@<mainnet-cluster>
kubectl create serviceaccount github-actions -n argocd
kubectl create clusterrolebinding github-actions --clusterrole=cluster-admin --serviceaccount=argocd:github-actions

# Generate token
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h) && echo $TOKEN

# Verify
curl -H "Authorization: Bearer $TOKEN" https://argocd-mainnet.ten.xyz/api/v1/applications
```

---

## 🎯 Next Steps

1. ✅ SSH into UAT cluster
2. ✅ Create service account & generate token
3. ✅ Test token with curl
4. ✅ Add `ARGOCD_TOKEN` to GitHub dev-testnet environment
5. ✅ Add `ARGOCD_TOKEN` to GitHub uat-testnet environment
6. ✅ Repeat for Sepolia & Mainnet (when ready)
7. ✅ Test deployment workflow

---

## 📞 Common Questions

**Q: Do dev and uat use the same token?**
A: Yes! Both use argocd-uat.ten.xyz, so same token for both GitHub environments.

**Q: How long is the token valid?**
A: 1 year (from --duration=8760h).

**Q: What if token expires?**
A: Regenerate new token and update GitHub secret (see Rotation Process).

**Q: Can I use the same token for all environments?**
A: No! Each cluster has different token. Generate separate tokens for uat, sepolia, mainnet.

**Q: Do I need to update anything when token expires?**
A: Just regenerate and update GitHub secret. Everything else stays same.

---

**All set! Your tokens are ready to use! 🎉**
